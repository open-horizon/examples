package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/open-horizon/edge-sync-service/common"
)

const (
	HZN_ESS_AUTH_VAR         = "HZN_ESS_AUTH"
	HZN_ESS_CERT_VAR         = "HZN_ESS_CERT"
	HZN_ESS_API_ADDRESS_VAR  = "HZN_ESS_API_ADDRESS"
	HZN_ESS_API_PORT_VAR     = "HZN_ESS_API_PORT"
	HZN_ESS_API_PROTOCOL_VAR = "HZN_ESS_API_PROTOCOL"
	MMS_OBJECT_TYPES_VAR     = "MMS_OBJECT_TYPES" // MMS_OBJECT_TYPES should be passed as an array in the userinput. it will be put into hzn-env-<agid> configmap
	// operator should set values in hzn-env-<agid> configmap as env var in MMS helper and consumer deployment
)

const (
	MMS_HELPER_STORAGE = "/ess-store"
	ESS_AUTH_FILE      = "/ess-auth/auth.json"
	ESS_CERT_FILE      = "/ess-cert/cert.pem"
)

var essApiAddress string
var essApiPort string

type AuthenticationCredential struct {
	Id      string `json:"id"`
	Token   string `json:"token"`
	Version string `json:"version"`
}

func ReadCredFromAuthFile(authFilePath string) (*AuthenticationCredential, error) {
	if authFile, err := os.Open(authFilePath); err != nil {
		return nil, fmt.Errorf("unable to open auth file %v, error: %v", authFilePath, err)
	} else if bytes, err := ioutil.ReadAll(authFile); err != nil {
		return nil, fmt.Errorf("unable to read auth file %v, error: %v", authFilePath, err)
	} else {
		authObj := new(AuthenticationCredential)
		if err := json.Unmarshal(bytes, authObj); err != nil {
			return nil, fmt.Errorf("unable to demarshal auth file %v, error: %v", authFilePath, err)
		} else {
			return authObj, nil
		}
	}
}

func getFromEnv(envName string, defaultVal string) string {
	val := os.Getenv(envName)
	if val == "" {
		val = defaultVal
	}
	return val
}

func getHttpClient(essCertFile string) (*http.Client, error) {
	caCert, err := ioutil.ReadFile(essCertFile)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
		},
	}

	client := http.Client{Transport: t, Timeout: 30 * time.Second}
	return &client, nil
}

func constructHttpGetRequest(url string, auth string) (*http.Request, error) {
	if req, err := http.NewRequest("GET", url, nil); err != nil {
		return nil, err
	} else {
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(auth))))
		return req, nil
	}
}

func invokeESSGet(httpClient *http.Client, url string, authValue string, respStruct interface{}) error {
	req, err := constructHttpGetRequest(url, authValue)
	if err != nil {
		return err
	}

	httpResp, _ := httpClient.Do(req)
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}
	if httpResp != nil {
		glog.V(2).Infof("Get response returns %v for %v", httpResp.Status, url)

		if outBytes, readErr := ioutil.ReadAll(httpResp.Body); readErr != nil {
			glog.Errorf("Error reading response body for %v, status code: %v , error was: %v", url, httpResp.Status, readErr)
			return readErr
		} else if httpResp.StatusCode != 200 && httpResp.StatusCode != 404 {
			msg := fmt.Sprintf("Get response returns %v, for %v", httpResp.Status, url)
			glog.Errorf(msg)
			return errors.New(msg)
		} else if httpResp.StatusCode == 404 {
			msg := fmt.Sprintf("Get response returns %v, for %v", httpResp.Status, url)
			glog.V(3).Infof(msg)
		} else {
			glog.V(3).Infof("Get response returns %v, body: %v for %v", httpResp.Status, string(outBytes), url)
			switch s := respStruct.(type) {
			case *[]byte:
				*s = outBytes
			case *string:
				*s = string(outBytes)
			default:
				if err = json.Unmarshal(outBytes, respStruct); err != nil {
					msg := fmt.Sprintf("Error unmarshal the response body from %v, error was: %v", url, err)
					glog.Errorf(msg)
					return errors.New(msg)
				}
			}
		}
	} else {
		glog.Errorf("received nil response from Get objects call")
	}
	return nil
}

func constructHttpPutRequest(url string, auth string) (*http.Request, error) {
	if req, err := http.NewRequest("PUT", url, nil); err != nil {
		return nil, err
	} else {
		req.Header.Add("Authorization", fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(auth))))
		return req, nil
	}
}

func invokeESSPut(httpClient *http.Client, url string, authValue string) error {
	req, err := constructHttpPutRequest(url, authValue)
	if err != nil {
		return err
	}

	httpResp, _ := httpClient.Do(req)
	if httpResp != nil && httpResp.Body != nil {
		defer httpResp.Body.Close()
	}

	if httpResp != nil && httpResp.StatusCode != 204 {
		return fmt.Errorf("receive error status code: %v from %v", httpResp.StatusCode, url)
	} else if httpResp == nil {
		glog.Errorf("received nil response from PUT objects call")
	}
	return nil
}

func writeDateStreamToFile(dataReader io.Reader, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if _, err := io.Copy(file, dataReader); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func convertToObjectList(objListString string) []string {
	objList := strings.Split(objListString, " ")
	return objList
}

func checkingMMSObject(authValue string, objectType string) {
	glog.V(2).Infof("Start checking MMS update for objectType: %v", objectType)

	// curl --cacert /path/to/cert.pem https://agent-service.<namespace>.svc.cluster.local:8443/api/v1/objects/<objectType> -u <org>/<service>:<token>
	getUpdatedObjectUrl := fmt.Sprintf("https://%v:%v/api/v1/objects/%v", essApiAddress, essApiPort, objectType)

	var getUpdatedObjectDataUrl string
	var markObjectReceivedUrl string
	var markObjectDeletedUrl string
	for {
		httpClient, err := getHttpClient(ESS_CERT_FILE)
		if err != nil {
			glog.Error("Error opening cert file %s, Error: %s", ESS_CERT_FILE, err)
			return
		}

		var metas []common.MetaData
		if err := invokeESSGet(httpClient, getUpdatedObjectUrl, authValue, &metas); err != nil {
			glog.Errorf("Failed to call ESS API: %v, error was: %v", getUpdatedObjectUrl, err)
		} else {
			glog.V(3).Infof("Receive %v objects updates", len(metas))
			for _, meta := range metas {
				filePath := fmt.Sprintf("%v/%v-%v", MMS_HELPER_STORAGE, meta.ObjectType, meta.ObjectID)
				if meta.Deleted {
					// this update is about to delete the object
					glog.V(3).Infof("MMS file %v/%v was deleted", meta.ObjectType, meta.ObjectID)

					// delete the file /ess-store/<objectType>-<objectId>
					if err = os.Remove(filePath); err != nil {
						glog.Errorf("failed to remove MMS file %v/%v, error was: %v", meta.ObjectType, meta.ObjectID, err)
					}

					// ack deleted
					// curl -sSLw "%{http_code}" -X PUT ${AUTH} ${CERT} $SOCKET $BASEURL/$OBJECT_TYPE/$OBJECT_ID/deleted
					// expect 204
					markObjectDeletedUrl = fmt.Sprintf("%v/%v/deleted", getUpdatedObjectUrl, meta.ObjectID)
					if err = invokeESSPut(httpClient, markObjectDeletedUrl, authValue); err != nil {
						glog.Errorf("failed to mark object %v/%v as deleted, error was: %v", meta.ObjectType, meta.ObjectID, err)
					}
					glog.V(3).Infof("mark object %v/%v is marked as 'deleted'", meta.ObjectType, meta.ObjectID)
				} else {
					// this update is about to get an object
					// get object and save to a file in the shared volume: /ess-store/<objectType>-<objectId>
					//curl -sSLw "%{http_code}" -o $OBJECT_ID ${AUTH} ${CERT} $SOCKET $BASEURL/$OBJECT_TYPE/$OBJECT_ID/data
					glog.V(3).Infof("Received new/updated %v/%v from MMS", meta.ObjectType, meta.ObjectID)
					getUpdatedObjectDataUrl = fmt.Sprintf("%v/%v/data", getUpdatedObjectUrl, meta.ObjectID)
					var dataBytes []byte

					if err = invokeESSGet(httpClient, getUpdatedObjectDataUrl, authValue, &dataBytes); err != nil {
						glog.Errorf("Failed to get object data from %v, error was: %v", getUpdatedObjectDataUrl, err)
						continue
					}
					glog.V(3).Infof("get data for %v/%v", meta.ObjectType, meta.ObjectID)

					glog.V(3).Infof("saving data at %v", filePath)
					r := bytes.NewReader(dataBytes)
					if err := writeDateStreamToFile(r, filePath); err != nil {
						glog.Errorf("failed to save data at %v", filePath)
						continue
					}

					// ack received
					// curl -sSLw "%{http_code}" -X PUT ${AUTH} ${CERT} $SOCKET $BASEURL/$OBJECT_TYPE/$OBJECT_ID/received
					// expect 200 or 204
					glog.V(3).Infof("mark object %v/%v is received", meta.ObjectType, meta.ObjectID)
					markObjectReceivedUrl = fmt.Sprintf("%v/%v/received", getUpdatedObjectUrl, meta.ObjectID)
					if err = invokeESSPut(httpClient, markObjectReceivedUrl, authValue); err != nil {
						glog.Errorf("failed to mark object %v/%v as received, error was: %v", meta.ObjectType, meta.ObjectID, err)
					}
					glog.V(3).Infof("mark object %v/%v is marked as 'received'", meta.ObjectType, meta.ObjectID)
				}
			}
		}
		time.Sleep(5 * time.Second)
	}

}

func main() {
	glog.V(3).Info("Starting checking updates for MMS objects")
	flag.Parse()
	/*
		env example:
			HZN_ESS_API_ADDRESS: agent-service.<namespace>.svc.cluster.local
			HZN_ESS_API_PORT: "8443"
			HZN_ESS_API_PROTOCOL: secure

		operator should bind these following 2 secrets into the /ess-auth/auth.json and /ess-cert/cert.pem
			HZN_ESS_AUTH: ess-auth-46e44a7d46530ecad0719e0ca24797054863717172407d9eb6fd755674c737fb
			HZN_ESS_CERT: ess-cert-46e44a7d46530ecad0719e0ca24797054863717172407d9eb6fd755674c737fb
	*/

	// get ess env
	essApiAddress = getFromEnv(HZN_ESS_API_ADDRESS_VAR, "")
	essApiPort = getFromEnv(HZN_ESS_API_PORT_VAR, "8443")
	objectTypesAsString := getFromEnv(MMS_OBJECT_TYPES_VAR, "") // need to add MMS_HELPER_OBJECT_TYPES as array in the userinput. echo $MMS_OBJECT_TYPE will return: model model1 model2 model3
	objectTypes := convertToObjectList(objectTypesAsString)

	// Get auth from auth.json
	var authValue string
	if auth, err := ReadCredFromAuthFile(ESS_AUTH_FILE); err != nil {
		glog.Error("error reading /ess-auth/auth.json file")
		return
	} else {
		// auth.json content: {"id":"<org>/abc","token":".........","version":"1.0.0"}
		authValue = fmt.Sprintf("%v:%v", auth.Id, auth.Token)
	}

	glog.V(3).Infof("Object types are: %v", objectTypesAsString)
	glog.V(3).Infof("authValue: %v", authValue)

	var wg sync.WaitGroup
	for _, objType := range objectTypes {
		wg.Add(1)
		go checkingMMSObject(authValue, objType)
	}
	wg.Wait()
	glog.V(3).Info("Done")
}

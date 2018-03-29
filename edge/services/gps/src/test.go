package main

/*

  A Blue Horizon Firmware Device API test shell

  Written by Glen Darling, November 2016.

*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"crypto/tls"
	"crypto/x509"
	"strconv"
	"strings"
	"envutil"
	"hgps"

	"gopkg.in/go-playground/validator.v8"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// JSON in this format is what is expected from the container under test
// Note that if spurious additional keys are found, it will still pass!
type LocationData struct {
	Latitude   float64 `json:"latitude" validate:"required"`
	Longitude  float64 `json:"longitude" validate:"required"`
	ElevationM float64 `json:"elevation" validate:"required"`
	AccuracyKM float64 `json:"accuracy_km"` // May be 0, considered empty
	LocSource  string  `json:"loc_source" validate:"required"`
	LastUpdate int64   `json:"loc_last_update" validate:"required"`
}

// JSON struct for satellite data
type SatellitesData struct {
        Satellites []SatelliteData `json:"satellites" description:"Array of data for each of the attached satellites."`
}

// JSON struct for data about a single satellite
type SatelliteData struct {
        PRN        float64  `json:"PRN"  description:"PRN ID of the satellite. 1-63 are GNSS satellites, 64-96 are GLONASS satellites, 100-164 are SBAS satellites"`
        Az         float64  `json:"az"   description:"Azimuth, degrees from true north."`
        El         float64  `json:"el"   description:"Elevation in degrees."`
        Ss         float64  `json:"ss"   description:"Signal strength in dB."`
        Used       bool     `json:"used" description:"Used in current solution? (SBAS/WAAS/EGNOS satellites may be flagged used if the solution has corrections from them, but not all drivers make this information available.)
"`
}

// Complete hgps data
type GpsData struct {
	Location  LocationData `json:"location" validate:"required"`
        Satellites []SatelliteData `json:"satellites" description:"Data about discovered satellites."`

}

// For sending to mqtt
type HeartbeatLocData struct {
	Lat float64  `json:"lat"`
	Lon float64  `json:"lon"`
	Alt float64  `json:"alt"`
}

type HeartbeatData struct {
	T int64  `json:"t"`
	R HeartbeatLocData  `json:"r"`
}

// Configuration variables and default values
const (
	REST_SERVICE_HOST = "REST_SERVICE_HOST"
	DEFAULT_REST_SERVICE_HOST = "gps"
	GPS_PORT = "HZN_GPS_PORT"
	DEFAULT_GPS_PORT = 31779
	// The reporting interval during the long running portion. Will report status and show full output and verbose output regularly after this many tests
	REPORTING_INTERVAL = "REPORTING_INTERVAL"
	DEFAULT_REPORTING_INTERVAL = 100000
	// The sleep time (in seconds) between each interval during the long running portion. 0 means it hammers the gps rest service as fast as it can.
	INTERVAL_SLEEP = "INTERVAL_SLEEP"
	DEFAULT_INTERVAL_SLEEP = 0

	HEARTBEAT_TO_MQTT = "HEARTBEAT_TO_MQTT"
	DEFAULT_HEARTBEAT_TO_MQTT = false
	HZN_DEVICE_ID = "HZN_DEVICE_ID"
	DEFAULT_HZN_DEVICE_ID = ""
	HZN_AGREEMENTID = "HZN_AGREEMENTID"
	DEFAULT_HZN_AGREEMENTID = ""
	HZN_HASH = "HZN_HASH"
	DEFAULT_HZN_HASH = ""
	POC_NUMBER = 20 	// the number assigned this poc in the verne topics
	MQTT_BROKER = "staging.bluehorizon.hovitos.engineering" 		// since this is a temporary test, we are hardcoding this
	CA_FILE = "ca-staging.pem"
)

// JSON format validator
var validate *validator.Validate

// Standard JSON headers
var json_hdrs = map[string]string{
	"Content-Type": "application/json",
}

// Format JSON for the human eye
func json_prettyprint(in []byte, prefix string, indent string) string {
	var out bytes.Buffer
	err := json.Indent(&out, in, prefix, indent)
	if err != nil {
		return "<< ERROR prettyprinting JSON response >>"
	}
	return string(out.Bytes())
}

// Assert the condition, and display if it fails, and return 1 if it succeeds
func assert(cond bool, cond_str string, msg string) int {
	if !cond {
		fmt.Printf("ASSERTION FAILED (%s):\n%s\n", cond_str, msg)
		return 0
	}
	return 1
}

// Increment the test number, announce a test by number, and return the number
func announce(test_num int) int {
	test_num += 1
	fmt.Printf("Running test #%d...\n", test_num)
	return test_num
}

// Make a REST API call and return err, status, headers, and body
func rest_call(method string,
	uri string,
	headers map[string]string,
	body []byte) (error, int, map[string][]string, []byte) {
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(body))
	for k, _ := range headers {
		req.Header.Set(k, headers[k])
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	// if 0 == assert(nil == err,
	// 	"nil == err",
	// 	fmt.Sprintf("REST call failed: %s %q\n", method, uri)) {
	// 	return err, 0, nil, nil
	// }
	if err != nil { return err, 0, nil, nil }
	defer resp.Body.Close()
	response_body, err := ioutil.ReadAll(resp.Body)
	// if 0 == assert(nil == err,
	// 	"nil == err",
	// 	fmt.Sprintf("Unable to read response from REST call: %s %q", method, uri)) {
	// 	return err, 0, nil, nil
	// }
	if err != nil { return err, 0, nil, nil }
	return nil, resp.StatusCode, resp.Header, response_body
}

// Expect a status code from a REST API (ignore headers and response body)
func test_status(expected_status int,
	method string,
	uri string,
	headers map[string]string,
	body []byte) int {
	err, status, _, r := rest_call(method, uri, headers, body)
	if nil != err {
		fmt.Printf("REST call failed: %s %q\n", method, uri)
		return 0
	}
	response_body := json_prettyprint(r, "    ", "  ")
	// fmt.Printf("test_status: uri=%q\nresp=%s\n", uri, response_body)
	if 0 == assert(expected_status == status,
			"expected_status == status",
			fmt.Sprintf("    Expected %d, but got %d. %s, %s :\n    %s", expected_status, status, method, uri, response_body)) {
		return 0
	}

	// All checks passed
	return 1
}

// Check status and validate JSON for the /v1/gps/location API
func test_content_location(expected_status int,
	method string,
	uri string,
	headers map[string]string,
	body []byte) int {
	err, status, _, r := rest_call(method, uri, headers, body)
	if nil != err {
		fmt.Printf("REST call failed: %s %q\n", method, uri)
		return 0
	}
	response_body := json_prettyprint(r, "    ", "  ")
	if 0 == assert(expected_status == status,
		"expected_status == status",
		fmt.Sprintf("    Expected %d, but got %d. %s, %s :\n    %s",
			expected_status, status, method, uri, response_body)) {
		return 0
	}
	// fmt.Printf("    Received:%s\n", response_body)

	// Decode the JSON body of the response
	var response LocationData
	err = json.Unmarshal(r, &response)
	if 0 == assert(nil == err,
		"nil == err",
		fmt.Sprintf("ERROR: failed to unmarshall JSON: '%s'", response_body)) {
		return 0
	}
	// fmt.Printf("    Unmarshalled: %v\n", response)

	// If not still searching, try to validate the unmarshalled structure
	if string(hgps.SEARCHING) != response.LocSource {
		validate_err := validate.Struct(response)
		if 0 == assert(nil == validate_err,
			"nil == validate_err",
			fmt.Sprintf("ERROR: failed to validate JSON location content (%v): '%s'", validate_err, response_body)) {
			return 0
		}
	}

	// All checks passed
	return 1
}

// Check status and validate JSON for the /v1/gps/location API
func test_content_satellites(expected_status int,
	method string,
	uri string,
	headers map[string]string,
	body []byte) int {
	err, status, _, r := rest_call(method, uri, headers, body)
	if nil != err {
		fmt.Printf("REST call failed: %s %q\n", method, uri)
		return 0
	}
	response_body := json_prettyprint(r, "    ", "  ")
	if 0 == assert(expected_status == status,
		"expected_status == status",
		fmt.Sprintf("    Expected %d, but got %d. %s, %s :\n    %s",
			expected_status, status, method, uri, response_body)) {
		return 0
	}
	// fmt.Printf("    Received:%s\n", response_body)

	// Decode the JSON body of the response
	var response SatellitesData
	err = json.Unmarshal(r, &response)
	if 0 == assert(nil == err,
		"nil == err",
		fmt.Sprintf("ERROR: failed to unmarshall JSON: '%s'", response_body)) {
		return 0
	}
	// fmt.Printf("    Unmarshalled: %v\n", response)

	// Validate each of the satellite entres in the unmarshalled structure
	for i, sat := range response.Satellites {
		validate_err := validate.Struct(sat)
		if 0 == assert(nil == validate_err,
			"nil == validate_err",
			fmt.Sprintf("ERROR: failed to validate JSON satellite #%d content (%v): '%s'", i, validate_err, response_body)) {
			return 0
		}
	}

	// All checks passed
	return 1
}

// Check status and validate JSON for the /v1/gps API
func test_content_gps(expected_status int,
	method string,
	uri string,
	headers map[string]string,
	body []byte) int {
	err, status, _, r := rest_call(method, uri, headers, body)
	if nil != err {
		fmt.Printf("REST call failed: %s %q\n", method, uri)
		return 0
	}
	response_body := json_prettyprint(r, "    ", "  ")
	if 0 == assert(expected_status == status,
		"expected_status == status",
		fmt.Sprintf("    Expected %d, but got %d. %s, %s :\n    %s",
			expected_status, status, method, uri, response_body)) {
		return 0
	}
	// fmt.Printf("    Received:%s\n", response_body)

	// Decode the JSON body of the response
	var response GpsData
	err = json.Unmarshal(r, &response)
	if 0 == assert(nil == err,
		"nil == err",
		fmt.Sprintf("ERROR: failed to unmarshall JSON: '%s'", response_body)) {
		return 0
	}
	// fmt.Printf("    Unmarshalled: %v\n", response)

	// If not still searching, try to validate the unmarshalled structure
	if string(hgps.SEARCHING) != response.Location.LocSource {
		validate_err := validate.Struct(response)
		if 0 == assert(nil == validate_err,
			"nil == validate_err",
			fmt.Sprintf("ERROR: failed to validate JSON gps content (%v): '%s'", validate_err, response_body)) {
			return 0
		}
	}

	// All checks passed
	return 1
}

// This "verbose mode" routine probes a REST API once, showing the results
func probe(expected_status int,
	method string,
	uri string,
	headers map[string]string,
	body []byte) int {
	fmt.Printf("--> REST request:\n")
	fmt.Printf("      Method=%q\n", method)
	fmt.Printf("      URI=%q\n", uri)
	err, status, _, r := rest_call(method, uri, headers, body)
	if nil != err {
		fmt.Printf("      *** ERROR: %v: REST call failed!\n", err)
		return 0
	}
	response_body := json_prettyprint(r, "    ", "  ")
	if 0 == assert(expected_status == status,
		"expected_status == status",
		fmt.Sprintf("      *** ERROR: Expected status %d, but got %d. %s, %s :\n    %s",
			expected_status, status, method, uri, response_body)) {
		return 0
	}
	fmt.Printf("<-- REST Response:\n")
	fmt.Printf("      Status: %d\n", status)
	fmt.Printf("      Body:%s\n", response_body)

	// Decode the JSON body of the response
	var response GpsData
	err = json.Unmarshal(r, &response)
	if 0 == assert(nil == err,
		"nil == err",
		fmt.Sprintf("      *** ERROR: %v: failed to unmarshall JSON: '%s'", err, response_body)) {
		return 0
	}
	// fmt.Printf("      Unmarshalled: %v\n", response)
	return 1
}

func NewTlsConfig() (*tls.Config, error) {
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(CA_FILE)
	if err != nil {
		return nil, fmt.Errorf("NewTlsConfig: unable to open certificate file: %v", err)
	}

	certpool.AppendCertsFromPEM(pemCerts)
	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
	}, nil
}

func min(a, b int) int {
	if a < b { return a }
	return b
}

func GetMqttClient(username string, password string) (MQTT.Client, error) {
	clientId := username[0:min(22,len(username))] 
	tlsconfig, err := NewTlsConfig()
	if err != nil { return nil, err }
	broker := "ssl://" + MQTT_BROKER + ":8883"
	opts := MQTT.NewClientOptions().AddBroker(broker)
	opts.SetClientID(clientId).SetTLSConfig(tlsconfig)
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetAutoReconnect(true) 	// i think this is the default, but just to make sure
	opts.SetWriteTimeout(30 * time.Second) 		// the default for this is no timeout

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.WaitTimeout(30 * time.Second) && token.Error() != nil {
		err := fmt.Errorf("GetMqttClient: unable to connect: %v", token.Error())
		return c, err
	}

	return c, nil
}

func RegisterAgreement(username string, password string) bool {
	c, err := GetMqttClient("public", "public")
    if err != nil {
    	fmt.Println("ERROR: could not get mqtt client for registration: ", err)
    	return false
    }
    defer c.Disconnect(250)
	msg := username + " " + password
    fmt.Println("RegisterAgreement: sending registration request for " + username + "...")
    if token := c.Publish("/registration", byte(2), false, msg); token.WaitTimeout(30 * time.Second) && token.Error() != nil {
        fmt.Printf("ERROR: RegisterAgreement: failed to send registration request: %v. Retrying registration...\n", token.Error())
	    if token := c.Publish("/registration", byte(2), false, msg); token.WaitTimeout(30 * time.Second) && token.Error() != nil {
	        fmt.Printf("ERROR: RegisterAgreement: failed to send registration request: %v. Giving up on registration.\n", token.Error())
	        return false
	    }
    } else {
    	fmt.Println("RegisterAgreement: registration request successfully sent.")
    	time.Sleep(20 * time.Second)
    }
    return true
}

func SendHeartbeat(notused MQTT.Client, username string, password string, serial string, gps_loc_data []byte) {
	// gps_loc_data comes from the gps rest api like: {"latitude": 45.403229, "longitude": -72.734066, "elevation": 101.07, "accuracy_km": 0, "loc_source": "Manual", "loc_last_update": 1491481763}
	// Decode the JSON body of the response
	var locData LocationData
	if err := json.Unmarshal(gps_loc_data, &locData); err != nil {
		fmt.Printf("ERROR: SendHeartbeat: unable to decode gps location data json: %v", err)
		return
	}

	// need to get the gps loc data to json like {"t":$ts,"r":{"lat":$HZN_LAT,"lon":$HZN_LON,"alt":0}}
	hbData := HeartbeatData{T: locData.LastUpdate, R: HeartbeatLocData{Lat: locData.Latitude, Lon: locData.Longitude, Alt: locData.ElevationM} }
	hbByteStr, err := json.Marshal(hbData)
	if err != nil {
		fmt.Printf("ERROR: SendHeartbeat: unable to encode gps location data to json: %v", err)
		return
	}
	hbStr := string(hbByteStr)

	c, err := GetMqttClient(username, password)
    if err != nil {
    	if strings.Contains(strings.ToLower(err.Error()), "not authori") {
    		// We are not registered properly, register and then try sending 1 more time
    		fmt.Printf("ERROR: not registered with mqtt properly (%v). Will register and trying sending again.\n", err)
			if !RegisterAgreement(username, password) { return } 	// RegisterAgreement() will display the error msg
			if c, err = GetMqttClient(username, password); err != nil {
		    	fmt.Println("ERROR: still could not get mqtt client for SendHeartbeat: ", err)
		    	return
			}
    	} else {
    		// Do not know what the connect problem was, bail on sending this HB
	    	fmt.Println("ERROR: could not get mqtt client for SendHeartbeat: ", err)
	    	return
	    }
    }
    defer c.Disconnect(2000) 		// the arg is the number of milliseconds to wait for in-flight msgs to complete
    topic := "/applications/in/" + username + "/public/h/"+serial+"/"+strconv.Itoa(POC_NUMBER)+"/"
    fmt.Printf("SendHeartbeat: sending hb: %v %v\n", topic, hbStr)
    if token :=  c.Publish(topic, 0, true, hbStr); token.WaitTimeout(30 * time.Second) && token.Error() != nil {
        fmt.Printf("ERROR: SendHeartbeat: unable to publish heartbeat to mqtt: %v", token.Error())
    } else {
    	fmt.Printf("SendHeartbeat: back from sending hb.\n")
    }
}

func main() {

	// Announce launch
	fmt.Printf("\nTest program started at: %s:\n", time.Now().UTC().Format(time.UnixDate))

	//
	// Get configuration from process environment
	//
	port := envutil.GetInt(GPS_PORT, DEFAULT_GPS_PORT, false)
	rest_host := envutil.GetString(REST_SERVICE_HOST, DEFAULT_REST_SERVICE_HOST, false)
	base_uri := fmt.Sprintf("http://%s:%d/", rest_host, port)
	reporting_interval := envutil.GetInt(REPORTING_INTERVAL, DEFAULT_REPORTING_INTERVAL, false)
	interval_sleep := envutil.GetInt(INTERVAL_SLEEP, DEFAULT_INTERVAL_SLEEP, false)
	hb_to_mqtt := envutil.GetBool(HEARTBEAT_TO_MQTT, DEFAULT_HEARTBEAT_TO_MQTT, false)
	device_id := envutil.GetString(HZN_DEVICE_ID, DEFAULT_HZN_DEVICE_ID, false)
	agreement_id := envutil.GetString(HZN_AGREEMENTID, DEFAULT_HZN_AGREEMENTID, false)
	hzn_hash := envutil.GetString(HZN_HASH, DEFAULT_HZN_HASH, false)
	
	// Check env vars
	if interval_sleep == 0 && hb_to_mqtt {
		fmt.Printf("ERROR: can not set %v to 0 and %v to true, you will kill mqtt.\n", INTERVAL_SLEEP, HEARTBEAT_TO_MQTT)
		os.Exit(2)
	}
	if hb_to_mqtt && (device_id=="" || agreement_id=="" || hzn_hash=="") {
		fmt.Printf("ERROR: if %v==true you must also set %v, %v, and %v.\n", HZN_DEVICE_ID, HEARTBEAT_TO_MQTT, HZN_AGREEMENTID, HZN_HASH)
		os.Exit(2)
	}

	fmt.Printf("\nBlue Horizon gps workload configuration:\n")
	fmt.Printf("  %v=%v\n", GPS_PORT, port)
	fmt.Printf("  %v=%v\n", "base_uri", base_uri)
	fmt.Printf("  %v=%v\n", REPORTING_INTERVAL, reporting_interval)
	fmt.Printf("  %v=%v\n", INTERVAL_SLEEP, interval_sleep)
	fmt.Printf("  %v=%v\n", HEARTBEAT_TO_MQTT, hb_to_mqtt)

	//
	// Initialize the JSON validator
	//
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)

	// Wait for the gps container to get started, in case we were started before it
	for {
		if err, status, _, _ := rest_call("GET", base_uri+"v1/gps", json_hdrs, []byte{}); err == nil && status == 200 { break }
		fmt.Println("The gps rest service is not available yet, still waiting...")
		time.Sleep(time.Duration(5) * time.Second)
	}
	fmt.Println("The gps rest service is now available.")

	// t will be incremented for each test, and p will be incremented for each successful test
	t := 0
	p := 0

	// Basic tests (just check status code)
	fmt.Printf("\nBasic REST API tests (checking HTTP status code only)...\n")

	// Bad methods, bad uris, empty bodies
	t = announce(t)
	p += test_status(400, "GET", base_uri, json_hdrs, []byte{})
	t = announce(t)
	p += test_status(400, "POST", base_uri, json_hdrs, []byte{})
	t = announce(t)
	p += test_status(400, "GET", base_uri+"foo", json_hdrs, []byte{})
	t = announce(t)
	p += test_status(400, "POST", base_uri+"bar", json_hdrs, []byte{})

	// Bad method for v1/gps API
	t = announce(t)
	p += test_status(405, "POST", base_uri+"v1/gps", json_hdrs, []byte{})
	t = announce(t)
	p += test_status(405, "PUT", base_uri+"v1/gps", json_hdrs, []byte{})
	t = announce(t)
	p += test_status(405, "DELETE", base_uri+"v1/gps", json_hdrs, []byte{})

	// Bad method for v1/gps/location API
	t = announce(t)
	p += test_status(405, "POST", base_uri+"v1/gps/location", json_hdrs, []byte{})
	t = announce(t)
	p += test_status(405, "PUT", base_uri+"v1/gps/location", json_hdrs, []byte{})
	t = announce(t)
	p += test_status(405, "DELETE", base_uri+"v1/gps/location", json_hdrs, []byte{})

	// Bad method for v1/gps/satellites API
	t = announce(t)
	p += test_status(405, "POST", base_uri+"v1/gps/satellites", json_hdrs, []byte{})
	t = announce(t)
	p += test_status(405, "PUT", base_uri+"v1/gps/satellites", json_hdrs, []byte{})
	t = announce(t)
	p += test_status(405, "DELETE", base_uri+"v1/gps/satellites", json_hdrs, []byte{})

	// Valid request for v1/gps (status check only)
	t = announce(t)
	p += test_status(200, "GET", base_uri+"v1/gps", json_hdrs, []byte{})

	// Valid request for v1/gps/location (status check only)
	t = announce(t)
	p += test_status(200, "GET", base_uri+"v1/gps/location", json_hdrs, []byte{})

	// Valid request for v1/gps/satellites (status check only)
	t = announce(t)
	p += test_status(200, "GET", base_uri+"v1/gps/satellites", json_hdrs, []byte{})

	// Content validation test (check status code and JSON fomat of response)
	fmt.Printf("\nResponse content tests (checking status code and response JSON)...\n")

	// Check content of valid request to the v1/gps/location API
	t = announce(t)
	p += test_content_location(200, "GET", base_uri+"v1/gps/location", json_hdrs, []byte{})

	// Check content of valid request to the v1/gps/satellites API
	t = announce(t)
	p += test_content_satellites(200, "GET", base_uri+"v1/gps/satellites", json_hdrs, []byte{})

	// Check content of valid request to the v1/gps API
	t = announce(t)
	p += test_content_gps(200, "GET", base_uri+"v1/gps", json_hdrs, []byte{})

	//
	// If all functional tests passed so far, start long-running test loop (until failure), logging results periodically
	//
	if p == t {
		fmt.Printf("\nAll functional tests passed.\nStress testing... (will report every %d tests; will halt on failure)\n\n", reporting_interval)
		// var mqttClient MQTT.Client
		if hb_to_mqtt {
			RegisterAgreement(agreement_id, hzn_hash)
			// var err error 		// if we do not create a new client each time we send, we eventually get hangs
			// if mqttClient, err = GetMqttClient(agreement_id, hzn_hash); err != nil { fmt.Printf("ERROR: could not get mqtt client: %v\n", err) }
		}
		for p == t {
			if 0 == (t % reporting_interval) {
				// time to report
				t = announce(t)
				p += probe(200, "GET", base_uri+"v1/gps", json_hdrs, []byte{})
				if hb_to_mqtt {
					err, status, _, r := rest_call("GET", base_uri+"v1/gps/location", json_hdrs, []byte{})
					if err != nil || status != 200 {
						fmt.Printf("ERROR: getting location to heartbeat to mqtt failed: %v, %v\n", status, err)
					} else {
						// if mqttClient != nil { SendHeartbeat(mqttClient, agreement_id, hzn_hash, device_id, r) }
						SendHeartbeat(nil, agreement_id, hzn_hash, device_id, r)
					}
				}
			} else {
				t += 1
				p += test_content_gps(200, "GET", base_uri+"v1/gps", json_hdrs, []byte{})
			}
			if interval_sleep > 0 { time.Sleep(time.Duration(interval_sleep) * time.Second) }
		}
		// mqttClient.Disconnect(250)
	}

	// Announce termination
	fmt.Printf("\n\nTest program ended at: %s:\n", time.Now().UTC().Format(time.UnixDate))

	// Summarize results
	if p == t {
		fmt.Printf("\n\nSUCCESS!  All %d tests passed.\n", t)
		os.Exit(0)
	} else {
		fmt.Printf("\n\nFAILURE!  Only %d of %d tests passed (%d errors)\n", p, t, t-p)
		fmt.Printf("Information about the %d failed test case(s) is above.\n", t-p)
		os.Exit(1)
	}
}

package wiotp

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func getEnv(keys ...string) (val string) {
	if len(keys) == 0 {
		panic("must give at least one key")
	}
	for _, key := range keys {
		val = os.Getenv(key)
		if val != "" {
			return
		}
	}
	if val == "" {
		fmt.Println("none of", keys, "are set")
		panic("can't any find set value")
	}
	return
}

func newTLSconfig() (tlsConfig *tls.Config, err error) {
	certpool := x509.NewCertPool()
	certPath := "messaging.pem"
	pemCerts, err := ioutil.ReadFile(certPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	certpool.AppendCertsFromPEM(pemCerts)
	tlsConfig = &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
	}
	return
}

// Connect attempts to connect to the correct MQTT broker, and returns a func to publish data.
func Connect() func([]byte) {
	testEnv := os.Getenv("WIOTP_TEST_ENV") // may be set to null string, so don't panic if it looks like it is unset.
	gwToken := getEnv("WIOTP_GW_TOKEN")
	orgID := getEnv("HZN_ORGANIZATION", "HZN_ORG_ID")
	hznDevID := getEnv("HZN_DEVICE_ID", "WIOTP_GW_ID")
	edgeMQTTip := getEnv("WIOTP_EDGE_MQTT_IP")
	agreementID := getEnv("HZN_AGREEMENTID")
	substrings := strings.Split(hznDevID, "@")
	if len(substrings) != 3 {
		panic("can't parse HZN_DEVICE_ID: " + hznDevID)
	}
	classID := substrings[0]
	gwType := substrings[1]
	devID := substrings[2]
	topic := "iot-2/type/" + gwType + "/id/" + devID + "/evt/status/fmt/json"

	// By default, connect to the local MQTT broker.
	connectFunc := func() (client mqtt.Client, err error) {
		return connectToEdgeConnector(classID, orgID, testEnv, gwType, devID, edgeMQTTip, agreementID)
	}
	// but if we don't an edgeMQTTip set, then fall back to talking directly to the cloud.
	if edgeMQTTip == "-" {
		fmt.Println("connecting to cloud")
		connectFunc = func() (client mqtt.Client, err error) {
			return connectToCloudMQTT(classID, orgID, testEnv, gwType, devID, gwToken, agreementID)
		}
	}
	// now that we have the correct connector func, we can try to use it to connect.
	mqttClient := tryConnect(connectFunc, 3, time.Second) // if something goes wrong and it has no tries left, it will panic.
	fmt.Println("connected to WIoTP")
	return func(msg []byte) {
		token := mqttClient.Publish(topic, 2, true, msg)
		token.Wait()
		if token.Error() != nil {
			fmt.Println(token.Error())
		}
	}
}

func tryConnect(connectToMQTT func() (mqtt.Client, error), ntryes int, sleepTime time.Duration) mqtt.Client {
	mqttClient, err := connectToMQTT()
	if ntryes == 0 {
		panic(err)
	}
	if err != nil {
		fmt.Println("can't connect to MQTT broker, will try", ntryes, "more times")
		time.Sleep(sleepTime)
		mqttClient = tryConnect(connectToMQTT, ntryes-1, sleepTime*2)
	}
	return mqttClient
}

func connectToCloudMQTT(classID, orgID, testEnv, gwType, gwID, gwToken, agreementID string) (mqttClient mqtt.Client, err error) {
	tlsConfig, err := newTLSconfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	broker := "ssl://" + orgID + ".messaging" + testEnv + ".internetofthings.ibmcloud.com:8883"
	clientID := classID + ":" + orgID + ":" + gwType + ":" + gwID
	fmt.Println("clientID:", clientID)
	fmt.Println("broker:", broker)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetTLSConfig(tlsConfig)
	opts.SetUsername("use-token-auth")
	opts.SetPassword(gwToken)

	mqttClient = mqtt.NewClient(opts)
	token := mqttClient.Connect()
	token.Wait()
	if err = token.Error(); err != nil {
		fmt.Println(err)
		return
	}
	return
}

func connectToEdgeConnector(classID, orgID, testEnv, gwType, gwID, edgeMQTTip, agreementID string) (mqttClient mqtt.Client, err error) {
	tlsConfig, err := newTLSconfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	broker := "ssl://" + edgeMQTTip + ":8883"
	clientID := "a:" + agreementID[0:36]
	fmt.Println("clientID:", clientID)
	fmt.Println("broker:", broker)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetTLSConfig(tlsConfig)

	mqttClient = mqtt.NewClient(opts)
	token := mqttClient.Connect()
	token.Wait()
	if err = token.Error(); err != nil {
		fmt.Println(err)
		return
	}
	return
}

package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var logger = log.New(os.Stdout, "audio_demo: ", log.Lshortfile)

type readSamplesMessage struct {
	NumSamples int `json:"numsamples"`
}

func getEnv(key string) (val string) {
	val = os.Getenv(key)
	if val == "" {
		panic(key + " is empty")
	}
	return
}

func main() {
	testEnv := os.Getenv("WIOTP_TEST_ENV") // may be set to null string, so don't panic if it looks like it is unset.
	orgID := getEnv("HZN_ORGANIZATION")
	gwType := getEnv("WIOTP_DEVICE_TYPE")
	hznDevID := getEnv("HZN_DEVICE_ID")
	gwToken := getEnv("WIOTP_DEVICE_AUTH_TOKEN")
	substrings := strings.Split(hznDevID, "@")
	if len(substrings) != 3 {
		panic("can't parse HZN_DEVICE_ID: " + hznDevID)
	}
	classID := substrings[0]
	devID := substrings[2]
	topic := "iot-2/type/" + gwType + "/id/" + devID + "/evt/status/fmt/json"
	mqttClient, err := getMqttClient(classID, orgID, testEnv, gwType, devID, gwToken)
	if err != nil {
		panic(err)
	}
	fmt.Println("connected to mqtt")
	// connect to the audio stream
	conn, err := net.Dial("tcp", "microphone:48926")
	if err != nil {
		panic(err)
	}
	buff := make([]byte, 1000)
	var numBytes int
	var numBytesMutex = sync.Mutex{}
	go func() {
		for {
			n, err := conn.Read(buff)
			if err != nil {
				panic(err)
			}
			numBytesMutex.Lock()
			numBytes += n
			//logger.Println("read", numBytes, "bytes")
			numBytesMutex.Unlock()
		}
	}()
	for {
		time.Sleep(10 * time.Second)
		numBytesMutex.Lock()
		byteCount := numBytes
		numBytes = 0
		numBytesMutex.Unlock()
		numSamples := readSamplesMessage{NumSamples: byteCount / 2}
		payload, err := json.Marshal(numSamples)
		if err != nil {
			logger.Println(err)
			continue
		}
		token := mqttClient.Publish(topic, 2, true, []byte(payload))
		token.Wait()
		if token.Error() != nil {
			logger.Println(token.Error())
			continue
		}
		fmt.Println("Published message:", string(payload))
	}
}

func newTLSconfig() (tlsConfig *tls.Config, err error) {
	certpool := x509.NewCertPool()
	certPath := "messaging.pem"
	pemCerts, err := ioutil.ReadFile(certPath)
	if err != nil {
		logger.Println(err)
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

func getMqttClient(classID, orgID, testEnv, gwType, gwID, gwToken string) (mqttClient mqtt.Client, err error) {
	tlsConfig, err := newTLSconfig()
	if err != nil {
		logger.Println(err)
		return
	}

	//mqtt.DEBUG = log.New(os.Stdout, "DEBUG", 0)
	//mqtt.ERROR = log.New(os.Stdout, "ERR", 0)

	broker := "ssl://" + orgID + ".messaging" + testEnv + ".internetofthings.ibmcloud.com:8883"
	clientID := classID + ":" + orgID + ":" + gwType + ":" + gwID
	logger.Println("broker:", broker)
	logger.Println("client id:", clientID)
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
		logger.Println(err)
		return
	}
	return
}

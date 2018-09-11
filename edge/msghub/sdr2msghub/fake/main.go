package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/open-horizon/examples/edge/msghub/sdr2msghub/audiolib"
)

type msghubConn struct {
	Producer sarama.SyncProducer
	Topic    string
}

// taken from cloud/sdr/data-ingest/example-go-clients/util/util.go
func populateConfig(config *sarama.Config, user, pw, apiKey string) error {
	config.ClientID = apiKey
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Net.TLS.Enable = true
	config.Net.SASL.User = user
	config.Net.SASL.Password = pw
	config.Net.SASL.Enable = true
	return nil
}

func connect(topic string) (conn msghubConn, err error) {
	conn.Topic = topic
	apiKey := getEnv("MSGHUB_API_KEY")
	username := apiKey[:16]
	password := apiKey[16:]
	brokerStr := getEnv("MSGHUB_BROKER_URL")
	brokers := strings.Split(brokerStr, ",")
	config := sarama.NewConfig()
	err = populateConfig(config, username, password, apiKey)
	if err != nil {
		return
	}
	conn.Producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return
	}
	return
}

// func (conn *msghubConn) publishAudio(audioMsg *audiolib.AudioMsg) (err error) {
func (conn *msghubConn) publishAudio(audioMsg *audiolib.AudioMsg) (err error) {
	// as AudioMsg implements the sarama.Encoder interface, we can pass it directly to ProducerMessage.
	msg := &sarama.ProducerMessage{Topic: conn.Topic, Key: nil, Value: audioMsg}
	partition, offset, err := conn.Producer.SendMessage(msg)
	if err != nil {
		log.Printf("FAILED to send message: %s\n", err)
	} else {
		log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
	}
	return
}

// read env vars from system with fall back.
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
		panic("can't find any set env var value")
	}
	return
}

func main() {
	devID := getEnv("HZN_ORG_ID") + "/" + getEnv("HZN_DEVICE_ID")
	mockAudio := "../../services/sdr/mock_audio.mp3" // if running it from the Makefile in edge/msghub/sdr2msghub
	audioBytes, err := ioutil.ReadFile(mockAudio)
	if err != nil {
		audioBytes, err = ioutil.ReadFile("../" + mockAudio) // if running it via go run in edge/msghub/sdr2msghub/fake
		if err != nil {
			panic(err)
		}
	}
	topic := getEnv("MSGHUB_TOPIC")
	fmt.Printf("using topic %s\n", topic)
	conn, err := connect(topic)
	if err != nil {
		panic(err)
	}
	fmt.Println("connected to msghub")
	msg := &audiolib.AudioMsg{
		Audio:         base64.StdEncoding.EncodeToString(audioBytes),
		Ts:            time.Now().Unix(),
		Freq:          123.45,
		ExpectedValue: 0.9,
		DevID:         devID,
		Lat:           42.214607,
		Lon:           -73.959494,
	}
	fmt.Println("sending sample")
	err = conn.publishAudio(msg)
	if err != nil {
		fmt.Println(err)
	}
}

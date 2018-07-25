package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/open-horizon/examples/edge/services/sdr/data_broker/audiolib"
)

type msghubConn struct {
	Producer sarama.SyncProducer
	Topic    string
}

// taken from cloud/sdr/sdr-data-ingest/example-go-clients/util/util.go
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
		panic("can't any find set value")
	}
	return
}

func main() {
	devID := getEnv("HZN_DEVICE_ID")
	audio, err := ioutil.ReadFile("../../librtlsdr/mock_audio.raw")
	if err != nil {
		panic(err)
	}
	conn, err := connect("sdr-audio")
	if err != nil {
		panic(err)
	}
	fmt.Println("connected to msghub")
	msg := &audiolib.AudioMsg{
		Audio:         audio,
		Ts:            time.Now(),
		Freq:          123.45,
		ExpectedValue: 0.9,
		DevID:         devID,
	}
	fmt.Println("sending sample")
	err = conn.publishAudio(msg)
	if err != nil {
		fmt.Println(err)
	}
}

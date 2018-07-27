// Example for consuming  messages from IBM Cloud Message Hub (kafka) using go.
// See README.md for setup requirements.

package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/Shopify/sarama"             // doc: https://godoc.org/github.com/Shopify/sarama
	cluster "github.com/bsm/sarama-cluster" // doc: http://godoc.org/github.com/bsm/sarama-cluster
	"github.com/open-horizon/examples/cloud/sdr/data-ingest/example-go-clients/util"
	"github.com/open-horizon/examples/cloud/sdr/data-processing/watson/nlu"
	"github.com/open-horizon/examples/cloud/sdr/data-processing/watson/stt"
	"github.com/open-horizon/examples/edge/services/sdr/data_broker/audiolib"
)

func Usage(exitCode int) {
	fmt.Printf("Usage: %s [-t <topic>] [-h] [-v]\n\nEnvironment Variables: MSGHUB_API_KEY, MSGHUB_BROKER_URL, MSGHUB_TOPIC\n", os.Args[0])
	os.Exit(exitCode)
}

const minConfidence = 0.5

func main() {
	sttUsername := util.RequiredEnvVar("STT_USERNAME", "")
	sttPassword := util.RequiredEnvVar("STT_PASSWORD", "")

	nluUsername := util.RequiredEnvVar("NLU_USERNAME", "")
	nluPassword := util.RequiredEnvVar("NLU_PASSWORD", "")

	// Get all of the input options
	var topic string
	flag.StringVar(&topic, "t", "", "topic")
	var help bool
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&util.VerboseBool, "v", false, "verbose")
	flag.Parse()
	if help {
		Usage(1)
	}

	apiKey := util.RequiredEnvVar("MSGHUB_API_KEY", "")
	username := apiKey[:16]
	password := apiKey[16:]
	util.Verbose("username: %s, password: %s\n", username, password)
	brokerStr := util.RequiredEnvVar("MSGHUB_BROKER_URL", "kafka01-prod02.messagehub.services.us-south.bluemix.net:9093,kafka02-prod02.messagehub.services.us-south.bluemix.net:9093,kafka03-prod02.messagehub.services.us-south.bluemix.net:9093,kafka04-prod02.messagehub.services.us-south.bluemix.net:9093,kafka05-prod02.messagehub.services.us-south.bluemix.net:9093")
	brokers := strings.Split(brokerStr, ",")
	if topic == "" {
		topic = util.RequiredEnvVar("MSGHUB_TOPIC", "sdr-audio")
	}

	util.Verbose("starting message hub consuming example...")

	if util.VerboseBool {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	err := util.PopulateConfig(&config.Config, username, password, apiKey) // add creds and tls info
	util.ExitOnErr(err)

	// init consumer
	consumer, err := cluster.NewConsumer(brokers, "my-consumer-group", []string{topic}, config)
	util.ExitOnErr(err)
	defer consumer.Close()

	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			if util.VerboseBool {
				log.Printf("Rebalanced: %+v\n", ntf)
			}
		}
	}()

	// consume messages, watch signals
	fmt.Printf("Consuming messages produced to %s...\n", topic)
	for {
		select {
		case msg, ok := <-consumer.Messages():
			if ok {
				// Got an audio clip, convert it to text
				var audioMsg audiolib.AudioMsg
				dec := gob.NewDecoder(bytes.NewReader(msg.Value))
				err := dec.Decode(&audioMsg)
				if err != nil {
					log.Println(err)
				}
				fmt.Println("got audio from device:", audioMsg.DevID, "on station:", audioMsg.Freq)
				transcript, err := stt.Transcribe(audioMsg.Audio, sttUsername, sttPassword)
				if err != nil {
					panic(err)
				}
				fmt.Println(transcript.Results)
				if util.VerboseBool {
					json, err := json.MarshalIndent(transcript.Results, "", "    ")
					if err != nil {
						panic(err)
					}
					fmt.Printf("STT: %s\n", string(json))
				}

				// Send each string of text with good confidence to NLU
				for _, r := range transcript.Results {
					altNum := 0 //todo: we only seem to get 1 alternative, not sure if it will always be that way
					if r.Final && r.Alternatives[altNum].Confidence > minConfidence {
						sentiments, err := nlu.Sentiment(r.Alternatives[altNum].Transcript, nluUsername, nluPassword)
						if err != nil {
							panic(err)
						}
						fmt.Printf("NLU: %s\n", sentiments)
					} else {
						fmt.Printf("Skipping: Final: %v, Confidence: %f\n", r.Final, r.Alternatives[altNum].Confidence)
					}
				}

				consumer.MarkOffset(msg, "") // mark message as processed
			}
		case <-signals:
			return
		}
	}
}

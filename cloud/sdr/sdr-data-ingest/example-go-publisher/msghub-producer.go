// Example for producing messages to IBM Cloud Message Hub (kafka) using go.
// See README.md for setup requirements.

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
	"flag"
	"github.com/Shopify/sarama"
	"github.com/open-horizon/examples/cloud/sdr/sdr-data-ingest/example-go-publisher/util"
)

func Usage(exitCode int) {
	fmt.Printf("Usage: %s [-t <topic>] [-s] [-h] [-v] [<message>]\n\nEnvironment Variables: MSGHUB_API_KEY, MSGHUB_BROKER_URL, MSGHUB_TOPIC\n", os.Args[0])
	os.Exit(exitCode)
}

func main() {
	// Get all of the input options
	var topic string
	flag.StringVar(&topic, "t", "", "topic")
	var sync, help bool
	flag.BoolVar(&sync, "s", false, "synchronous")
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&util.VerboseBool, "v", false, "verbose")
	flag.Parse()
	if help { Usage(1) }

	message := ""
	if flag.NArg() >= 1 {
		message = flag.Arg(0)
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

	if !sync {

	} else {
		util.Verbose("starting message hub producing example...")

		if util.VerboseBool {
			sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
		}

		client, err := util.NewClient(username, password, apiKey, brokers)
		util.ExitOnErr(err)

		producer, err := sarama.NewSyncProducerFromClient(client)
		util.ExitOnErr(err)

		defer util.Close(client, producer, nil, nil)

		if message != "" {
			util.Verbose("producing the specified msg to %s...\n", topic)
			err = util.SendSyncMessage(producer, topic, message)
			util.ExitOnErr(err)
		} else {
			numMsgs := 10
			util.Verbose("producing %d generated msgs to %s...\n", numMsgs, topic)
			msgs := make([]string, numMsgs)
			for i := 0; i < numMsgs; i++ {
				msgs[i] = "message "+strconv.Itoa(i)
			}
			err = util.SendSyncMessages(producer, topic, msgs)
			util.ExitOnErr(err)
			/* can do this in a single call instead...
			for i := 0; i < numMsgs; i++ {
				err = util.SendSyncMessage(producer, topic, "message "+strconv.Itoa(i))
				util.ExitOnErr(err)
			}
			*/
		}
	}

	util.Verbose("message hub producing example complete")
}

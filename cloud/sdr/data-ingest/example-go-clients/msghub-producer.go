// Example for producing messages to IBM Cloud Message Hub (kafka) using go.
// See README.md for setup requirements.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/Shopify/sarama" // doc: https://godoc.org/github.com/Shopify/sarama
	"github.com/open-horizon/examples/cloud/sdr/data-ingest/example-go-clients/util"
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
	if help {
		Usage(1)
	}

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

	util.Verbose("starting message hub producing example...")

	if util.VerboseBool {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	config, err := util.NewConfig(username, password, apiKey)
	util.ExitOnErr(err)

	if !sync {
		// Produce msgs asynchronously
		producer, err := sarama.NewAsyncProducer(brokers, config)
		util.ExitOnErr(err)

		defer func() {
			if err := producer.Close(); err != nil {
				log.Fatalln(err)
			}
		}()

		// Trap SIGINT to trigger a shutdown.
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		// First queue up the msgs in our own buffered channel. In a real application you would probably just send
		// your msg to the producer.Input() channel right now. We are instead queuing them up here so in a single
		// select below we can send the msgs and list for results.
		numMsgs := 10
		ch := make(chan *sarama.ProducerMessage, numMsgs)
		for i := 0; i < numMsgs; i++ {
			ch <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(message + " " + strconv.Itoa(i))}
			// fmt.Printf("DEBUG: adding msg %d to ch\n", i)
		}

		// Now enqueue the msgs in the async producer while also listening for errors and successes
		var enqueued, errors, successes int
	ProducerLoop:
		for {
			select {
			case producerMsg := <-ch:
				producer.Input() <- producerMsg
				enqueued++
				// fmt.Printf("DEBUG: enqueue %d\n", enqueued)
			case err := <-producer.Errors():
				log.Println("Failed to produce message", err)
				errors++
				// fmt.Printf("DEBUG: error %d\n", errors)
				if (errors + successes) >= numMsgs {
					break ProducerLoop
				}
			case <-producer.Successes():
				successes++
				// fmt.Printf("DEBUG: success %d\n", successes)
				if (errors + successes) >= numMsgs {
					break ProducerLoop
				}
			case <-signals:
				break ProducerLoop
			}
		}

		fmt.Printf("%d messages produced to topic: %s; successes: %d errors: %d\n", enqueued, topic, successes, errors)
	} else {
		// Produce msgs synchronously
		producer, err := sarama.NewSyncProducer(brokers, config)
		util.ExitOnErr(err)

		defer func() {
			if err := producer.Close(); err != nil {
				log.Fatalln(err)
			}
		}()

		if message != "" {
			util.Verbose("producing msg '%s' to %s...\n", message, topic)
			err = SendSyncMessage(producer, topic, message)
			util.ExitOnErr(err)
		} else {
			numMsgs := 10
			util.Verbose("producing %d generated msgs to %s...\n", numMsgs, topic)
			msgs := make([]string, numMsgs)
			for i := 0; i < numMsgs; i++ {
				msgs[i] = "message " + strconv.Itoa(i)
			}
			err = SendSyncMessages(producer, topic, msgs)
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

func SendSyncMessage(producer sarama.SyncProducer, topic, msg string) error {
	pMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	partition, offset, err := producer.SendMessage(pMsg)
	if err != nil {
		return err
	}
	fmt.Printf("Message produced to topic: %s, partition: %d, offset: %d\n", topic, partition, offset)
	return nil
}

func SendSyncMessages(producer sarama.SyncProducer, topic string, msgs []string) error {
	pMsgs := make([]*sarama.ProducerMessage, len(msgs))
	for i, m := range msgs {
		pMsgs[i] = &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(m),
		}
	}

	err := producer.SendMessages(pMsgs)
	if err != nil {
		return err
	}
	fmt.Printf("%d messages produced to topic: %s\n", len(msgs), topic)
	return nil
}

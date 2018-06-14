// Example for consuming  messages from IBM Cloud Message Hub (kafka) using go.
// See README.md for setup requirements.

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"flag"
	"github.com/Shopify/sarama"		// doc: https://godoc.org/github.com/Shopify/sarama
	cluster "github.com/bsm/sarama-cluster"		// doc: http://godoc.org/github.com/bsm/sarama-cluster
	"github.com/open-horizon/examples/cloud/sdr/sdr-data-ingest/example-go-clients/util"
)

func Usage(exitCode int) {
	fmt.Printf("Usage: %s [-t <topic>] [-h] [-v]\n\nEnvironment Variables: MSGHUB_API_KEY, MSGHUB_BROKER_URL, MSGHUB_TOPIC\n", os.Args[0])
	os.Exit(exitCode)
}

func main() {
	// Get all of the input options
	var topic string
	flag.StringVar(&topic, "t", "", "topic")
	var help bool
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&util.VerboseBool, "v", false, "verbose")
	flag.Parse()
	if help { Usage(1) }

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
	err := util.PopulateConfig(&config.Config, username, password, apiKey)		// add creds and tls info
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
				//fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				if util.VerboseBool {
					fmt.Printf("%s: %s (partition: %d, offset: %d)\n", msg.Topic, msg.Value, msg.Partition, msg.Offset)
				} else {
					fmt.Printf("%s: %s\n", msg.Topic, msg.Value)
				}
				consumer.MarkOffset(msg, "")	// mark message as processed
			}
		case <-signals:
			return
		}
	}

	/* This can only listen to 1 partition, or a hardcoded number of partitions...
	client, err := util.NewClient(username, password, apiKey, brokers)
	util.ExitOnErr(err)
	consumer, err := sarama.NewConsumerFromClient(client)
	util.ExitOnErr(err)
	defer util.Close(client, nil, nil, consumer)
	callback := func(msg *sarama.ConsumerMessage) {
			if util.VerboseBool {
				fmt.Printf("%s: %s (partition: %d, offset: %d)\n", msg.Topic, string(msg.Value), msg.Partition, msg.Offset)
			} else {
				fmt.Printf("%s: %s\n", msg.Topic, string(msg.Value))
			}
		}
	fmt.Printf("Consuming messages produced to %s...\n", topic)
	err = util.ConsumePartition(consumer, topic, 0, callback)
	util.ExitOnErr(err)
	*/

	util.Verbose("message hub consuming example complete")   // we should never get here
}


// Not currently used, because can only listen to 1 partition...
func ConsumePartition(consumer sarama.Consumer, topic string, partition int32, callback func(*sarama.ConsumerMessage)) error {
	partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		return err
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			callback(msg)
		}
	}
}

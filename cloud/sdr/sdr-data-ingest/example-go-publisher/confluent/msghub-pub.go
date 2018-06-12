// !!This example currently does not work with IBM Message Hub!!
// Example for publishing messages to IBM Cloud Message Hub (kafka) using go

/* Current build/run requirements:
- install librdkafka (on MacOS X: brew install librdkafka pkg-config
- go get -u github.com/confluentinc/confluent-kafka-go/kafka
*/ 

package main

import (
	"fmt"
	"os"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func requiredEnvVar(name, defaultVal string) string {
	v := os.Getenv(name)
	if defaultVal != "" {
		v = defaultVal
	}
	if v == "" {
		fmt.Printf("Error: environment variable '%s' must be defined.\n", name)
		os.Exit(2)
	}
	return v
}

func main() {
	fmt.Println("Starting message hub publishing example...")

	apiKey := requiredEnvVar("MSGHUB_API_KEY", "")
	username := apiKey[:16]
	password := apiKey[16:]
	fmt.Printf("username: %s, password: %s\n", username, password)
	brokerUrls := requiredEnvVar("MSGHUB_BROKER_URL", "kafka01-prod02.messagehub.services.us-south.bluemix.net:9093,kafka02-prod02.messagehub.services.us-south.bluemix.net:9093,kafka03-prod02.messagehub.services.us-south.bluemix.net:9093,kafka04-prod02.messagehub.services.us-south.bluemix.net:9093,kafka05-prod02.messagehub.services.us-south.bluemix.net:9093")
	topic := requiredEnvVar("MSGHUB_TOPIC", "sdr-audio")

	// For valid kafka config values, see https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
	kafkaConfig := kafka.ConfigMap{
		//"bootstrap.servers": brokerUrls,
		"metadata.broker.list": brokerUrls,
		"sasl.mechanisms": "PLAIN",
		"sasl.username": username,
		"sasl.password": password,
	}

	p, err := kafka.NewProducer(&kafkaConfig)
	if err != nil {
		panic(err)
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	fmt.Printf("publishing a few msgs to %s...\n", topic)
	for _, word := range []string{"someaudiodata" /*, "moreaudiodata"*/} {
		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
		if err != nil {
			fmt.Printf("Error from Produce(): %v\n", err)
		}
	}

	// Wait for message deliveries
	p.Flush(15 * 1000)

	fmt.Println("Message hub publishing example complete.")
}

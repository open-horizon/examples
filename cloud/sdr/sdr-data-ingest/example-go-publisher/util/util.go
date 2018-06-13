// Utility functions for the examples for publishing and consuming messages to/fromm IBM Cloud Message Hub (kafka) using go

/* Todos:
- implement async producer
- decide what to do with tls key and pem
*/

package util

import (
	"fmt"
	"log"
	"os"
	"strings"
	"crypto/tls"
	"github.com/Shopify/sarama"
)

var VerboseBool bool

func Verbose(msg string, args ...interface{}) {
	if !VerboseBool {
		return
	}
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(os.Stderr, "[verbose] "+msg, args...) // send to stderr so it doesn't mess up stdout if they are piping that to jq or something like that
}

// RequiredEnvVar gets an env var value. If a default value is not supplied and the env var is not defined, a fatal error is displayed.
func RequiredEnvVar(name, defaultVal string) string {
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

func ExitOnErr(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(2)
	}
}

func TlsConfig(certFile, keyFile string) (*tls.Config, error) {
	cer, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{cer}}, nil
}

func NewClient(user, pw, apiKey string, brokers []string) (sarama.Client, error) {
	config := sarama.NewConfig()
	err := PopulateConfig(config, user, pw, apiKey)
	if err != nil {
		return nil, err
	}

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func PopulateConfig(config *sarama.Config, user, pw, apiKey string) error {
	tlsConfig, err := TlsConfig("server.pem", "server.key")
	if err != nil {
		return err
	}

	config.ClientID = apiKey
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = tlsConfig
	config.Net.SASL.User = user
	config.Net.SASL.Password = pw
	config.Net.SASL.Enable = true
	return nil
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
	fmt.Printf("Message published to topic: %s, partition: %d, offset: %d\n", topic, partition, offset)
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
	fmt.Printf("%d messages published to topic: %s\n", len(msgs), topic)
	return nil
}

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

func Close(client sarama.Client, syncProducer sarama.SyncProducer, asyncProducer sarama.AsyncProducer, consumer sarama.Consumer) {
	if syncProducer != nil {
		if err := syncProducer.Close(); err != nil {
			log.Fatalln(err)
		}
	}
	if asyncProducer != nil {
		if err := asyncProducer.Close(); err != nil {
			log.Fatalln(err)
		}
	}
	if consumer != nil {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}
	if err := client.Close(); err != nil {
		log.Fatalln(err)
	}
}

// Example for publishing messages to IBM Cloud Message Hub (kafka) using go

/* Current build/run requirements:
- go get github.com/Shopify/sarama
- openssl genrsa -out server.key 2048
- openssl req -new -x509 -key server.key -out server.pem -days 3650
- export MSGHUB_API_KEY='abcdefg'
*/ 

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"strconv"
	"flag"
	"crypto/tls"
	"github.com/Shopify/sarama"
)

var VerboseBool bool

func Usage(exitCode int) {
	fmt.Printf("Usage: %s [-h] [-v] [<message>]\n\nEnvironment Variables: MSGHUB_API_KEY, MSGHUB_BROKER_URL, MSGHUB_TOPIC\n", os.Args[0])
	os.Exit(exitCode)
}

func main() {
	var help bool
	flag.BoolVar(&help, "h", false, "help")
	flag.BoolVar(&VerboseBool, "v", false, "verbose")
	flag.Parse()
	if help { Usage(1) }

	message := ""
	if flag.NArg() >= 1 {
		message = flag.Arg(0)
	}

	apiKey := requiredEnvVar("MSGHUB_API_KEY", "")
	username := apiKey[:16]
	password := apiKey[16:]
	//fmt.Printf("username: %s, password: %s\n", username, password)
	brokerStr := requiredEnvVar("MSGHUB_BROKER_URL", "kafka01-prod02.messagehub.services.us-south.bluemix.net:9093,kafka02-prod02.messagehub.services.us-south.bluemix.net:9093,kafka03-prod02.messagehub.services.us-south.bluemix.net:9093,kafka04-prod02.messagehub.services.us-south.bluemix.net:9093,kafka05-prod02.messagehub.services.us-south.bluemix.net:9093")
	brokers := strings.Split(brokerStr, ",")
	topic := requiredEnvVar("MSGHUB_TOPIC", "sdr-audio")

	Verbose("starting message hub publishing example...")

	if VerboseBool {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	client, err := NewClient(username, password, apiKey, brokers)
	exitOnErr(err)

	producer, err := sarama.NewSyncProducerFromClient(client)
	exitOnErr(err)

	defer Close(client, producer, nil)

	if message != "" {
		Verbose("publishing the specified msg to %s...\n", topic)
		err = SendMessage(producer, topic, message)
		exitOnErr(err)
	} else {
		Verbose("publishing a few generated msgs to %s...\n", topic)
		for i := 0; i < 10; i++ {
			err = SendMessage(producer, topic, "message "+strconv.Itoa(i))
			exitOnErr(err)
		}
	}

	Verbose("message hub publishing example complete")
}


func Verbose(msg string, args ...interface{}) {
	if !VerboseBool {
		return
	}
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Fprintf(os.Stderr, "[verbose] "+msg, args...) // send to stderr so it doesn't mess up stdout if they are piping that to jq or something like that
}

func requiredEnvVar(name, defaultVal string) string {
	v := os.Getenv(name)
	if defaultVal != "" {
		v = defaultVal
	}
	if v == "" {
		fmt.Printf("Error: environment variable '%s' must be defined.\n", name)
		Usage(2)
	}
	return v
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(2)
	}
}

func tlsConfig(certFile, keyFile string) (*tls.Config, error) {
	cer, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{cer}}, nil
}

func NewClient(user, pw, apiKey string, brokers []string) (sarama.Client, error) {
	tlsConfig, err := tlsConfig("server.pem", "server.key")
	if err != nil {
		return nil, err
	}

	config := sarama.NewConfig()
	config.ClientID = apiKey
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = tlsConfig
	config.Net.SASL.User = user
	config.Net.SASL.Password = pw
	config.Net.SASL.Enable = true

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func SendMessage(producer sarama.SyncProducer, topic, msg string) error {
	pMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	partition, offset, err := producer.SendMessage(pMsg)
	if err != nil {
		return err
	}
	fmt.Printf("Message published to topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

func Close(client sarama.Client, producer sarama.SyncProducer, consumer sarama.Consumer) {
	if producer != nil {
		if err := producer.Close(); err != nil {
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

// Utility functions for the examples for producing and consuming messages to/fromm IBM Cloud Message Hub (kafka) using go

package util

import (
	"fmt"
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

/* If you want to create your own client object, it can be done like this. We create the producer and
	consumer objects directly, and let them own the client (so they also close them at the end).
func NewClient(user, pw, apiKey string, brokers []string) (sarama.Client, error) {
	config, err := NewConfig(user, pw, apiKey)
	if err != nil {
		return nil, err
	}

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
*/

func NewConfig(user, pw, apiKey string) (*sarama.Config, error) {
	config := sarama.NewConfig()
	err := PopulateConfig(config, user, pw, apiKey)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func PopulateConfig(config *sarama.Config, user, pw, apiKey string) error {
	/* If you want to create your own certificate and use it, you can...
	tlsConfig, err := TlsConfig("server.pem", "server.key")
	if err != nil {
		return err
	}
	*/

	config.ClientID = apiKey
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Net.TLS.Enable = true
	// config.Net.TLS.Config = tlsConfig
	config.Net.SASL.User = user
	config.Net.SASL.Password = pw
	config.Net.SASL.Enable = true
	return nil
}

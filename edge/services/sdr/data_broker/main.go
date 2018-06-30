package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/open-horizon/examples/edge/services/sdr/data_broker/audiolib"
	rtlsdr "github.com/open-horizon/examples/edge/services/sdr/librtlsdr/rtlsdrclientlib"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

func opIsSafe(a string) bool {
	safeOPtypes := []string{
		"Const",
		"Placeholder",
		"Conv2D",
		"Cast",
		"Div",
		"StatelessRandomNormal",
		"ExpandDims",
		"AudioSpectrogram",
		"DecodeRaw",
		"Reshape",
		"MatMul",
		"Sum",
		"Softmax",
		"Squeeze",
		"RandomUniform",
		"Identity",
	}
	for _, b := range safeOPtypes {
		if b == a {
			return true
		}
	}
	return false
}

type model struct {
	Sess    *tf.Session
	InputPH tf.Output
	Output  tf.Output
}

func (m *model) goodness(audio []byte) (value float32, err error) {
	inputTensor, err := tf.NewTensor(string(audio))
	if err != nil {
		return
	}
	result, err := m.Sess.Run(map[tf.Output]*tf.Tensor{m.InputPH: inputTensor}, []tf.Output{m.Output}, nil)
	if err != nil {
		return
	}
	value = result[0].Value().([]float32)[0]
	return
}

func newModel(path string) (m model, err error) {
	def, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	graph := tf.NewGraph()
	err = graph.Import(def, "")
	if err != nil {
		panic(err)
	}
	ops := graph.Operations()
	unsafeOPs := map[string]bool{}
	graphIsUnsafe := false
	for _, op := range ops {
		if !opIsSafe(op.Type()) {
			unsafeOPs[op.Type()] = true
			graphIsUnsafe = true
		}
	}
	if graphIsUnsafe {
		fmt.Println("The following OP types are not in whitelist:")
		for op := range unsafeOPs {
			fmt.Println(op)
		}
		err = errors.New("unsafe OPs")
		return
	}
	outputOP := graph.Operation("output")
	if outputOP == nil {
		err = errors.New("output OP not found")
		return
	}
	m.Output = outputOP.Output(0)

	inputPHOP := graph.Operation("input/Placeholder")
	if inputPHOP == nil {
		err = errors.New("input OP not found")
		return
	}
	m.InputPH = inputPHOP.Output(0)
	m.Sess, err = tf.NewSession(graph, nil)
	return
}

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
	msg := &sarama.ProducerMessage{Topic: conn.Topic, Key: nil, Value: audioMsg}
	partition, offset, err := conn.Producer.SendMessage(msg)
	if err != nil {
		log.Printf("FAILED to send message: %s\n", err)
	} else {
		log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
	}
	return
}

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

var hostname string = "rtlsdr"

func main() {
	alt_addr := os.Getenv("RTLSDR_ADDR")
	if alt_addr != "" {
		fmt.Println("connecting to remote rtlsdr:", alt_addr)
		hostname = alt_addr
	}
	devID := getEnv("HZN_DEVICE_ID")
	m, err := newModel("model.pb")
	if err != nil {
		panic(err)
	}
	fmt.Println("model loaded")
	conn, err := connect("sdr-audio")
	if err != nil {
		panic(err)
	}
	fmt.Println("connected to msghub")
	stations, err := rtlsdr.GetCeilingSignals(hostname, -8)
	if err != nil {
		panic(err)
	}
	fmt.Println("found", len(stations), "stations")
	for {
		for _, station := range stations {
			audio, err := rtlsdr.GetAudio(hostname, int(station))
			if err != nil {
				panic(err)
			}
			val, err := m.goodness(audio)
			if err != nil {
				panic(err)
			}
			fmt.Println(station, val)
			if val > 0.5 {
				msg := &audiolib.AudioMsg{
					Audio:         audio,
					Ts:            time.Now(),
					Freq:          station,
					ExpectedValue: val,
					DevID:         devID,
				}
				fmt.Println("sending sample")
				err = conn.publishAudio(msg)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("not sending")
			}
		}
	}
}

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/open-horizon/examples/cloud/sdr/sdr-data-ingest/example-go-clients/util"
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
	outputOP := graph.Operation("Squeeze")
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

func connect(topic string) (conn msghubConn, err error) {
	conn.Topic = topic
	apiKey := util.RequiredEnvVar("MSGHUB_API_KEY", "")
	username := apiKey[:16]
	password := apiKey[16:]
	util.Verbose("username: %s, password: %s\n", username, password)
	brokerStr := util.RequiredEnvVar("MSGHUB_BROKER_URL", "kafka01-prod02.messagehub.services.us-south.bluemix.net:9093,kafka02-prod02.messagehub.services.us-south.bluemix.net:9093,kafka03-prod02.messagehub.services.us-south.bluemix.net:9093,kafka04-prod02.messagehub.services.us-south.bluemix.net:9093,kafka05-prod02.messagehub.services.us-south.bluemix.net:9093")
	brokers := strings.Split(brokerStr, ",")
	config, err := util.NewConfig(username, password, apiKey)
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

func main() {
	//m, err := newModel("train/conv_model.pb")
	m, err := newModel("train/random_model/random_model.pb")
	if err != nil {
		panic(err)
	}
	fmt.Println("model loaded")
	conn, err := connect("sdr-audio")
	if err != nil {
		panic(err)
	}
	fmt.Println("connected to msghub")
	stations, err := rtlsdr.GetCeilingSignals("localhost", -13)
	if err != nil {
		panic(err)
	}
	fmt.Println("found", len(stations), "stations")
	for {
		for _, station := range stations {
			fmt.Println("starting freq", station)
			audio, err := rtlsdr.GetAudio("localhost", int(station))
			if err != nil {
				panic(err)
			}
			val, err := m.goodness(audio)
			if err != nil {
				panic(err)
			}
			if val > 0.5 {
				msg := &audiolib.AudioMsg{
					Audio:         audio,
					Ts:            time.Now(),
					Freq:          station,
					ExpectedValue: val,
					DevID:         "isaac_test_desktop",
				}
				err = conn.publishAudio(msg)
				if err != nil {
					fmt.Println(err)
				}
			}
			fmt.Println(station, val)
		}
	}
}

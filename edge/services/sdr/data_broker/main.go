package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

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

func getAudio(freq int) (audio []byte, err error) {
	resp, err := http.Get("http://localhost:8080/audio/" + strconv.Itoa(freq))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("bad resp")
		return
	}
	defer resp.Body.Close()
	audio, err = ioutil.ReadAll(resp.Body)
	return
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

func main() {
	//m, err := newModel("train/conv_model.pb")
	m, err := newModel("train/random_model/random_model.pb")
	if err != nil {
		panic(err)
	}
	for {
		//audio, err := getAudio(91100000)
		audio, err := getAudio(89700000)
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile("/tmp/dat1", audio, 0644)
		val, err := m.goodness(audio)
		if err != nil {
			panic(err)
		}
		fmt.Println(val)
	}
}

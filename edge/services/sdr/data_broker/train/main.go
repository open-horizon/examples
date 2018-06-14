package main

import (
	"fmt"
	"io/ioutil"
	"os"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"

	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

func model(s *op.Scope, input tf.Output) (class tf.Output) {
	zero := op.Const(s.SubScope("zero"), int64(0))
	one := op.Const(s.SubScope("one"), int64(1))
	//two := op.Const(s.SubScope("two"), int64(2))
	three := op.Const(s.SubScope("three"), int64(3))
	//fzero := op.Const(s.SubScope("f0"), float32(0))
	seed := op.Const(s.SubScope("seed"), []int64{9, 4})

	filter1 := op.StatelessRandomNormal(s.SubScope("filter1"), op.Const(s.SubScope("filter1_dims"), []int64{7, 5, 1, 3}), seed)
	filter2 := op.StatelessRandomNormal(s.SubScope("filter2"), op.Const(s.SubScope("filter2_dims"), []int64{7, 5, 3, 5}), seed)
	filter3 := op.StatelessRandomNormal(s.SubScope("filter3"), op.Const(s.SubScope("filter3_dims"), []int64{7, 5, 5, 5}), seed)
	fc := op.StatelessRandomNormal(s.SubScope("fc"), op.Const(s.SubScope("fc_shape"), []int64{5 * 5, 7}), seed)
	readout := op.StatelessRandomNormal(s.SubScope("readout"), op.Const(s.SubScope("readout_shape"), []int64{7, 2}), seed)

	//rawData := op.ReadFile(s.SubScope("read_audio"), op.Const(s.SubScope("filename"), "audio.raw"))

	pcm := op.Div(s.SubScope("div_by_2_16"),
		op.Cast(s, op.DecodeRaw(s.SubScope("decode_u16"), input, tf.Uint16), tf.Float),
		op.Const(s.SubScope("65536"), float32(65536)),
	)
	spectrogram := op.AudioSpectrogram(s, op.ExpandDims(s, pcm, one), int64(100), int64(100))
	conv1 := op.Conv2D(s.SubScope("conv1"),
		op.ExpandDims(s.SubScope("add_chan"), spectrogram, three),
		filter1,
		[]int64{1, 5, 2, 1},
		"VALID",
	)
	conv2 := op.Conv2D(s.SubScope("conv2"),
		conv1,
		filter2,
		[]int64{1, 5, 2, 1},
		"VALID",
	)
	conv3 := op.Conv2D(s.SubScope("conv3"),
		conv2,
		filter3,
		[]int64{1, 5, 2, 1},
		"VALID",
	)
	flat := op.Reshape(s, conv3, op.Const(s.SubScope("flat"), []int64{74, 25}))
	timeOutput := op.MatMul(s, flat, fc)
	sum := op.Sum(s, timeOutput, zero)
	class = op.Squeeze(s, op.Softmax(s.SubScope("output"), op.MatMul(s.SubScope("readout"), op.ExpandDims(s.SubScope("readout"), sum, zero), readout)))
	return
}

func trainingDataQueue(fileName, target tf.Output) (readBatch tf.Output, init tf.Operation) {
  dataset := op.TensorSliceDataset(s, components, [])
}

func main() {
	fmt.Println(tf.Version())
	s := op.NewScope()
	rawDataPH := op.Placeholder(s.SubScope("input"), tf.String)
	output := model(s.SubScope("model"), rawDataPH)
	fmt.Println(output.Op.Name())
	fmt.Println(rawDataPH.Op.Name())
	graph, err := s.Finalize()
	if err != nil {
		panic(err)
	}
	file, err := os.Create("conv_model.pb")
	if err != nil {
		panic(err)
	}
	_, err = graph.WriteTo(file)
	if err != nil {
		panic(err)
	}
	sess, err := tf.NewSession(graph, nil)
	if err != nil {
		panic(err)
	}
	rawBytes, err := ioutil.ReadFile("audio.raw")
	if err != nil {
		panic(err)
	}
	inputTensor, err := tf.NewTensor(string(rawBytes))
	if err != nil {
		panic(err)
	}
	result, err := sess.Run(map[tf.Output]*tf.Tensor{rawDataPH: inputTensor}, []tf.Output{output}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(result[0].Value())
}

package main

import (
	"fmt"
	"os"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"

	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

func main() {
	fmt.Println(tf.Version())
	s := op.NewScope()
	rawDataPH := op.Placeholder(s.SubScope("input"), tf.String)
	_ = rawDataPH

	output := op.Squeeze(s, op.Softmax(s.SubScope("output"), op.RandomUniform(s, op.Const(s, []int64{1, 2}), tf.Float)))
	_ = output
	graph, err := s.Finalize()
	if err != nil {
		panic(err)
	}
	file, err := os.Create("random_model.pb")
	if err != nil {
		panic(err)
	}
	_, err = graph.WriteTo(file)
	if err != nil {
		panic(err)
	}
}

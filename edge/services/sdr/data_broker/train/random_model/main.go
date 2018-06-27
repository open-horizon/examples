package main

import (
	"fmt"
	"os"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"

	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

func namedIdentity(scope *op.Scope, input tf.Output, name string) (output tf.Output) {
	if scope.Err() != nil {
		return
	}
	opspec := tf.OpSpec{
		Type: "Identity",
		Input: []tf.Input{
			input,
		},
		Name: name,
	}
	op := scope.AddOperation(opspec)
	return op.Output(0)
}

func main() {
	fmt.Println(tf.Version())
	s := op.NewScope()
	rawDataPH := op.Placeholder(s.SubScope("input"), tf.String)
	_ = rawDataPH

	output := namedIdentity(s, op.Squeeze(s, op.Softmax(s.SubScope("output"), op.RandomUniform(s, op.Const(s, []int64{1, 2}), tf.Float))), "output")
	_ = output
	graph, err := s.Finalize()
	if err != nil {
		panic(err)
	}
	file, err := os.Create("rand_model.pb")
	if err != nil {
		panic(err)
	}
	_, err = graph.WriteTo(file)
	if err != nil {
		panic(err)
	}

	s = op.NewScope()
	rawDataPH = op.Placeholder(s.SubScope("input"), tf.String)
	_ = rawDataPH

	output = namedIdentity(s, op.Const(s, []float32{1.0, 0.0}), "output")
	_ = output
	graph, err = s.Finalize()
	if err != nil {
		panic(err)
	}
	file, err = os.Create("yes_model.pb")
	if err != nil {
		panic(err)
	}
	_, err = graph.WriteTo(file)
	if err != nil {
		panic(err)
	}

	s = op.NewScope()
	rawDataPH = op.Placeholder(s.SubScope("input"), tf.String)
	_ = rawDataPH

	output = namedIdentity(s, op.Const(s, []float32{0.0, 1.0}), "output")
	_ = output
	graph, err = s.Finalize()
	if err != nil {
		panic(err)
	}
	file, err = os.Create("no_model.pb")
	if err != nil {
		panic(err)
	}
	_, err = graph.WriteTo(file)
	if err != nil {
		panic(err)
	}
}

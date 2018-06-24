package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/is8ac/tfutils/descend"
	"github.com/is8ac/tfutils/descend/models"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"

	"github.com/tensorflow/tensorflow/tensorflow/go/op"
)

func model(s *op.Scope, input tf.Output, filter1, filter2, filter3, fcWeights, readout tf.Output) (class tf.Output) {
	batchSize := input.Shape().Size(0)

	conv1 := op.Conv2D(s.SubScope("conv1"),
		input,
		filter1,
		[]int64{1, 2, 2, 1},
		"VALID",
	)
	conv2 := op.Conv2D(s.SubScope("conv2"),
		conv1,
		filter2,
		[]int64{1, 2, 2, 1},
		"VALID",
	)
	conv3 := op.Conv2D(s.SubScope("conv3"),
		conv2,
		filter3,
		[]int64{1, 2, 2, 1},
		"VALID",
	)
	flat := op.Reshape(s, conv3, op.Const(s.SubScope("flat"), []int64{batchSize, -1}))
	timeOutput := op.MatMul(s, flat, fcWeights)
	class = op.Softmax(s.SubScope("output"), op.MatMul(s.SubScope("readout"), timeOutput, readout))
	return
}

func preprocessAudio(s *op.Scope, audio tf.Output) (ffts tf.Output) {
	pcm := op.Div(s.SubScope("div_by_2_16"),
		op.Cast(s, op.DecodeRaw(s.SubScope("decode_u16"), audio, tf.Uint16), tf.Float),
		op.Const(s.SubScope("65536"), float32(65536)),
	)
	ffts = op.ExpandDims(s.SubScope("add_chan"),
		op.Squeeze(s,
			op.AudioSpectrogram(s,
				op.ExpandDims(s, pcm, op.Const(s.SubScope("one"), int64(1))),
				100,
				100,
			),
		),
		op.Const(s.SubScope("two"), int64(2)),
	)
	return
}

//filter1 := op.StatelessRandomNormal(randomS.SubScope("filter1"), op.Const(s.SubScope("filter1_dims"), []int64{5, 5, 1, 3}), seed)
//filter2 := op.StatelessRandomNormal(randomS.SubScope("filter2"), op.Const(s.SubScope("filter2_dims"), []int64{5, 5, 3, 5}), seed)
//filter3 := op.StatelessRandomNormal(randomS.SubScope("filter3"), op.Const(s.SubScope("filter3_dims"), []int64{5, 5, 5, 5}), seed)
//fc := op.StatelessRandomNormal(randomS.SubScope("fc"), op.Const(s.SubScope("fc_shape"), []int64{5 * 5 * 583, 7}), seed)
//readout := op.StatelessRandomNormal(randomS.SubScope("readout"), op.Const(s.SubScope("readout_shape"), []int64{7, 2}), seed)

func makeConvModel(input, target tf.Output) (
	lossFunc descend.LossFunc,
	size int64,
	makeFinalizeAccuracy func(*op.Scope, tf.Output, tf.Output, tf.Output) tf.Output,
) {
	paramDefs := []models.ParamDef{
		models.ParamDef{Name: "filter1", Shape: tf.MakeShape(5, 5, 1, 3)},
		models.ParamDef{Name: "filter2", Shape: tf.MakeShape(5, 5, 3, 5)},
		models.ParamDef{Name: "filter3", Shape: tf.MakeShape(5, 5, 5, 5)},
		models.ParamDef{Name: "fc", Shape: tf.MakeShape(5*5*583, 7)},
		models.ParamDef{Name: "readout", Shape: tf.MakeShape(7, 2)},
	}
	unflatten, size := models.MakeUnflatten(paramDefs)

	lossFunc = func(s *op.Scope, params tf.Output) (loss tf.Output) {
		layerParams := unflatten(s.SubScope("unflatten"), params)
		output := model(s.SubScope("model"), input, layerParams[0], layerParams[1], layerParams[2], layerParams[3], layerParams[4])
		loss = op.Mean(s,
			op.Sum(s,
				op.Square(s, op.Sub(s, output, target)),
				op.Const(s.SubScope("one"), int64(1)),
			),
			op.Const(s.SubScope("zero"), int64(0)),
		)
		return
	}
	return
}

func initDataQueue(s *op.Scope,
	preprocess func(*op.Scope, tf.Output) tf.Output,
) (
	initLoadDatum func(*tf.Session) (func(string, bool) error, error),
	closeQueue *tf.Operation,
	dequeueFFTs, dequeueLabels tf.Output,
) {
	queue := op.FIFOQueueV2(s, []tf.DataType{tf.Float, tf.Float}, op.FIFOQueueV2Shapes([]tf.Shape{tf.MakeShape(4692, 65, 1), tf.MakeShape(2)}))

	sizeVar := op.VarHandleOp(s.SubScope("n"), tf.Int32, tf.ScalarShape(), op.VarHandleOpContainer("n"))
	increment := op.AssignAddVariableOp(s, sizeVar, op.Const(s.SubScope("1i32"), int32(1)))
	reset := op.AssignVariableOp(s, sizeVar, op.Const(s.SubScope("0i32"), int32(0)))
	readSize := op.ReadVariableOp(s, sizeVar, tf.Int32)

	fileNamePH := op.Placeholder(s.SubScope("filename"), tf.String, op.PlaceholderShape(tf.ScalarShape()))
	readFile := op.ReadFile(s, fileNamePH)
	fft := preprocess(s.SubScope("preprocess"), readFile)
	labelIndexPH := op.Placeholder(s.SubScope("label_index"), tf.Bool)
	label := op.OneHot(s,
		op.Cast(s, labelIndexPH, tf.Int64),
		op.Const(s.SubScope("two"), int32(2)),
		op.Const(s.SubScope("1f32"), float32(1)),
		op.Const(s.SubScope("0f32"), float32(0)),
	)
	enqueue := op.QueueEnqueueV2(s.WithControlDependencies(increment), queue, []tf.Output{fft, label})
	dequeueComponents := op.QueueDequeueUpToV2(s, queue, readSize, []tf.DataType{tf.Float, tf.Float})
	dequeueFFTs = dequeueComponents[0]
	dequeueLabels = dequeueComponents[1]
	closeQueue = op.QueueCloseV2(s, queue)
	initLoadDatum = func(sess *tf.Session) (loadDatum func(string, bool) error, err error) {
		_, err = sess.Run(nil, nil, []*tf.Operation{reset})
		if err != nil {
			return
		}
		loadDatum = func(path string, label bool) (err error) {
			fileNameTensor, err := tf.NewTensor(path)
			if err != nil {
				return
			}
			labelTensor, err := tf.NewTensor(label)
			if err != nil {
				return
			}
			_, err = sess.Run(map[tf.Output]*tf.Tensor{fileNamePH: fileNameTensor, labelIndexPH: labelTensor}, nil, []*tf.Operation{enqueue})
			return
		}
		return
	}
	return
}

func nextBatch(s *op.Scope, ffts, labels, seed tf.Output, n int64) (batchFFTs, batchLabels tf.Output, init *tf.Operation) {
	outputTypes := []tf.DataType{ffts.DataType(), labels.DataType()}
	outputShapes := []tf.Shape{tf.MakeShape(n, 4692, 65, 1), tf.MakeShape(n, 2)}
	preBatchOutputShapes := []tf.Shape{tf.MakeShape(4692, 65, 1), tf.MakeShape(2)}
	dataset := op.TensorSliceDataset(s, []tf.Output{ffts, labels}, preBatchOutputShapes)
	repeatDataset := op.RepeatDataset(s, dataset, op.Const(s.SubScope("count"), int64(-1)), outputTypes, preBatchOutputShapes)
	shuffleDataset := op.ShuffleDataset(s,
		repeatDataset,
		op.Const(s.SubScope("buffer_size"), int64(1000)),
		seed,
		seed,
		outputTypes,
		preBatchOutputShapes,
	)
	batchDataset := op.BatchDataset(s, shuffleDataset, op.Const(s.SubScope("batch_size"), n), outputTypes, outputShapes)
	iterator := op.Iterator(s, "", "", outputTypes, outputShapes)
	next := op.IteratorGetNext(s, iterator, outputTypes, outputShapes)
	init = op.MakeIterator(s, batchDataset, iterator)
	batchFFTs = next[0]
	batchLabels = next[1]
	return
}

func loadClass(path string, load func(string) error) (err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, file := range files {
		err = load(filepath.Join(path, file.Name()))
		if err != nil {
			log.Println("skipping", file.Name(), ":", err.Error())
			continue
		}
	}
	return
}

const gobalSeed int64 = 0

func main() {
	fmt.Println(tf.Version())

	s := op.NewScope()
	initLoadDatum, closeQueue, dequeueFFTs, dequeueLabels := initDataQueue(s.SubScope("queue"), preprocessAudio)
	seed := op.Const(s.SubScope("seed"), []int64{gobalSeed, gobalSeed})
	scalarSeed := op.Const(s.SubScope("scalar_seed"), int64(gobalSeed))
	fftsBatch, labelsBatch, initOP := nextBatch(s.SubScope("dataset"), dequeueFFTs, dequeueLabels, scalarSeed, 30)

	randomS := s.SubScope("random_params")
	filter1 := op.StatelessRandomNormal(randomS.SubScope("filter1"), op.Const(s.SubScope("filter1_dims"), []int64{5, 5, 1, 3}), seed)
	filter2 := op.StatelessRandomNormal(randomS.SubScope("filter2"), op.Const(s.SubScope("filter2_dims"), []int64{5, 5, 3, 5}), seed)
	filter3 := op.StatelessRandomNormal(randomS.SubScope("filter3"), op.Const(s.SubScope("filter3_dims"), []int64{5, 5, 5, 5}), seed)
	fc := op.StatelessRandomNormal(randomS.SubScope("fc"), op.Const(s.SubScope("fc_shape"), []int64{5 * 5 * 583, 7}), seed)
	readout := op.StatelessRandomNormal(randomS.SubScope("readout"), op.Const(s.SubScope("readout_shape"), []int64{7, 2}), seed)

	output := model(s.SubScope("model"), fftsBatch, filter1, filter2, filter3, fc, readout)
	loss := op.Mean(s,
		op.Sum(s,
			op.Square(s,
				op.Sub(s, output, labelsBatch),
			),
			op.Const(s.SubScope("one"), int64(1)),
		),
		op.Const(s.SubScope("zero"), int64(0)),
	)

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
	loadDatum, err := initLoadDatum(sess)
	if err != nil {
		panic(err)
	}
	loadClass("label/good", func(fileName string) error {
		return loadDatum(fileName, true)
	})
	loadClass("label/nongood", func(fileName string) error {
		return loadDatum(fileName, false)
	})

	_, err = sess.Run(nil, nil, []*tf.Operation{initOP})
	if err != nil {
		panic(err)
	}
	// now we can close the queue
	_, err = sess.Run(nil, nil, []*tf.Operation{closeQueue})
	if err != nil {
		panic(err)
	}

	result, err := sess.Run(map[tf.Output]*tf.Tensor{}, []tf.Output{loss, labelsBatch, output}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(result[0].Value())
	fmt.Println(result[1].Value())
	fmt.Println(result[2].Value())
}

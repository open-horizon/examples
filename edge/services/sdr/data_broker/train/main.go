package main

import (
	"errors"
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

func convModel(s *op.Scope, input tf.Output, filter1, filter2, filter3, fcWeights, readout tf.Output) (class tf.Output) {
	batchSize := input.Shape().Size(0)

	conv1 := op.Conv2D(s.SubScope("conv1"),
		input,
		filter1,
		[]int64{1, 4, 2, 1},
		"VALID",
	)
	conv2 := op.Conv2D(s.SubScope("conv2"),
		op.Relu(s.SubScope("l2_relu"), conv1),
		filter2,
		[]int64{1, 4, 3, 1},
		"VALID",
	)
	conv3 := op.Conv2D(s.SubScope("conv3"),
		op.Relu(s.SubScope("l3_relu"), conv2),
		filter3,
		[]int64{1, 4, 3, 1},
		"VALID",
	)
	flat := op.Reshape(s, conv3, op.Const(s.SubScope("flat"), []int64{batchSize, -1}))
	timeOutput := op.MatMul(s, flat, fcWeights)
	class = op.Softmax(s.SubScope("output"), op.MatMul(s.SubScope("readout"), timeOutput, readout))
	return
}

func simpleModel(s *op.Scope, input tf.Output, filter, weights tf.Output) (class tf.Output) {
	batchSize := input.Shape().Size(0)
	// ugly hack
	if batchSize == -1 {
		batchSize = 1
	}
	fmt.Println("batch_size:", batchSize)
	conv1 := op.Conv2D(s.SubScope("conv1"),
		input,
		filter,
		[]int64{1, 7, 7, 1},
		"VALID",
	)
	//fmt.Println("conv1", conv1.Shape())
	flatInput := op.Reshape(s, conv1, op.Const(s.SubScope("shape"), []int64{batchSize, -1}))
	class = op.Softmax(s, op.MatMul(s, flatInput, weights))
	return
}

var simpleModelParamDefs = []models.ParamDef{
	models.ParamDef{Name: "filter1", Shape: tf.MakeShape(23, 11, 1, 3)},
	models.ParamDef{Name: "weights", Shape: tf.MakeShape(14148, 2)},
}

func makeSimpleModel(input, target tf.Output) (
	lossFunc descend.LossFunc,
	size int64,
	makeFinalizeAccuracy func(*op.Scope, tf.Output, tf.Output, tf.Output) tf.Output,
) {
	unflatten, size := models.MakeUnflatten(simpleModelParamDefs)

	lossFunc = func(s *op.Scope, params tf.Output) (loss tf.Output) {
		layerParams := unflatten(s.SubScope("unflatten"), params)
		output := simpleModel(s.SubScope("model"), input, layerParams[0], layerParams[1])
		loss = op.Mean(s,
			op.Sum(s,
				op.Square(s, op.Sub(s, output, target)),
				op.Const(s.SubScope("one"), int64(1)),
			),
			op.Const(s.SubScope("zero"), int64(0)),
		)
		return
	}
	makeFinalizeAccuracy = func(s *op.Scope,
		params tf.Output,
		testInputs, testTargets tf.Output,
	) (
		accuracy tf.Output,
	) {
		layerParams := unflatten(s.SubScope("accuracy_unflatten"), params)
		actual := simpleModel(s.SubScope("model"), testInputs, layerParams[0], layerParams[1])
		actualLabels := op.ArgMax(s, actual, op.Const(s.SubScope("argmax_dim"), int32(-1)), op.ArgMaxOutputType(tf.Int32))
		targetLabels := op.ArgMax(s.SubScope("targets"), testTargets, op.Const(s.SubScope("argmax_dim"), int32(-1)), op.ArgMaxOutputType(tf.Int32))
		correct := op.Reshape(s.SubScope("correct"), op.Equal(s, actualLabels, targetLabels), op.Const(s.SubScope("all"), []int32{-1}))
		accuracy = op.Mean(s, op.Cast(s.SubScope("accuracy"), correct, tf.Float), op.Const(s.SubScope("mean_dim"), int32(0)))
		return
	}

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
				500,
				500,
			),
		),
		op.Const(s.SubScope("two"), int64(2)),
	)
	return
}

func makeConvModel(input, target tf.Output) (
	lossFunc descend.LossFunc,
	size int64,
	makeFinalizeAccuracy func(*op.Scope, tf.Output, tf.Output, tf.Output) tf.Output,
) {
	paramDefs := []models.ParamDef{
		models.ParamDef{Name: "filter1", Shape: tf.MakeShape(7, 5, 1, 3)},
		models.ParamDef{Name: "filter2", Shape: tf.MakeShape(7, 5, 3, 5)},
		models.ParamDef{Name: "filter3", Shape: tf.MakeShape(7, 5, 5, 5)},
		models.ParamDef{Name: "fc", Shape: tf.MakeShape(845, 7)},
		models.ParamDef{Name: "readout", Shape: tf.MakeShape(7, 2)},
	}
	unflatten, size := models.MakeUnflatten(paramDefs)

	lossFunc = func(s *op.Scope, params tf.Output) (loss tf.Output) {
		layerParams := unflatten(s.SubScope("unflatten"), params)
		output := convModel(s.SubScope("model"), input, layerParams[0], layerParams[1], layerParams[2], layerParams[3], layerParams[4])
		loss = op.Mean(s,
			op.Sum(s,
				op.Square(s, op.Sub(s, output, target)),
				op.Const(s.SubScope("one"), int64(1)),
			),
			op.Const(s.SubScope("zero"), int64(0)),
		)
		return
	}
	makeFinalizeAccuracy = func(s *op.Scope,
		params tf.Output,
		testInputs, testTargets tf.Output,
	) (
		accuracy tf.Output,
	) {
		layerParams := unflatten(s.SubScope("accuracy_unflatten"), params)
		actual := convModel(s.SubScope("model"), testInputs, layerParams[0], layerParams[1], layerParams[2], layerParams[3], layerParams[4])
		actualLabels := op.ArgMax(s, actual, op.Const(s.SubScope("argmax_dim"), int32(-1)), op.ArgMaxOutputType(tf.Int32))
		targetLabels := op.ArgMax(s.SubScope("targets"), testTargets, op.Const(s.SubScope("argmax_dim"), int32(-1)), op.ArgMaxOutputType(tf.Int32))
		correct := op.Reshape(s.SubScope("correct"), op.Equal(s, actualLabels, targetLabels), op.Const(s.SubScope("all"), []int32{-1}))
		accuracy = op.Mean(s, op.Cast(s.SubScope("accuracy"), correct, tf.Float), op.Const(s.SubScope("mean_dim"), int32(0)))
		return
	}

	return
}

func initDataQueue(s *op.Scope,
	preprocess func(*op.Scope, tf.Output) tf.Output,
	n int32,
) (
	initLoadDatum func(*tf.Session) (func(string, bool) error, error),
	closeQueue *tf.Operation,
	dequeueFFTs, dequeueLabels tf.Output,
) {
	queue := op.FIFOQueueV2(s, []tf.DataType{tf.Float, tf.Float}, op.FIFOQueueV2Shapes([]tf.Shape{tf.MakeShape(938, 257, 1), tf.MakeShape(2)}))

	sizeVar := op.VarHandleOp(s.SubScope("n"), tf.Int32, tf.ScalarShape(), op.VarHandleOpContainer("n"))
	increment := op.AssignAddVariableOp(s, sizeVar, op.Const(s.SubScope("1i32"), int32(1)))
	reset := op.AssignVariableOp(s, sizeVar, op.Const(s.SubScope("0i32"), int32(0)))
	readSize := op.ReadVariableOp(s, sizeVar, tf.Int32)
	_ = readSize

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
	dequeueComponents := op.QueueDequeueManyV2(s, queue, op.Const(s.SubScope("n"), n), []tf.DataType{tf.Float, tf.Float})
	fmt.Println(s.Err())

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
	outputShapes := []tf.Shape{tf.MakeShape(n, 938, 257, 1), tf.MakeShape(n, 2)}
	preBatchOutputShapes := []tf.Shape{tf.MakeShape(938, 257, 1), tf.MakeShape(2)}
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

func varCache(s *op.Scope, input tf.Output, shape tf.Shape, name string) (init *tf.Operation, output tf.Output) {
	variable := op.VarHandleOp(s, input.DataType(), shape, op.VarHandleOpSharedName(name))
	init = op.AssignVariableOp(s, variable, input)
	output = op.ReadVariableOp(s, variable, input.DataType())
	return
}

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

func getOP(graph *tf.Graph, name string) (operation *tf.Operation, err error) {
	operation = graph.Operation(name)
	if operation == nil {
		err = errors.New("can't find operation " + name)
		return
	}
	return
}

const gobalSeed int64 = 0

func main() {
	const subSize = 30
	const globalSeed = 42
	const batchSize = 150
	const searchSize float32 = 0.0003
	const gradsScale float32 = 0.005
	const dataSize int64 = 1100
	fmt.Println(tf.Version())

	s := op.NewScope()
	initLoadDatum, closeQueue, dequeueFFTs, dequeueLabels := initDataQueue(s.SubScope("queue"), preprocessAudio, int32(dataSize))
	initFFTcache, readFFTs := varCache(s.SubScope("fft_cache"), dequeueFFTs, tf.MakeShape(dataSize, 938, 257, 1), "ffts")
	initLabelsCache, readLabels := varCache(s.SubScope("labels_cache"), dequeueLabels, tf.MakeShape(dataSize, 2), "labels")
	fmt.Println("shape:", readFFTs.Shape())

	scalarSeed := op.Const(s.SubScope("scalar_seed"), int64(gobalSeed))
	fftsBatch, labelsBatch, initOP := nextBatch(s.SubScope("dataset"), readFFTs, readLabels, scalarSeed, batchSize)

	step := op.Const(s.SubScope("search_size"), searchSize)
	lossFunc, size, makeFinalizeAccuracy := makeSimpleModel(fftsBatch, labelsBatch)
	//lossFunc, size, makeFinalizeAccuracy := makeConvModel(fftsBatch, labelsBatch)
	fmt.Println("size:", size)
	updatesPH := op.Placeholder(s.SubScope("updates"), tf.Float, op.PlaceholderShape(tf.MakeShape(subSize)))
	randomExpand := descend.MakeRandomExpand(size, 42)
	initSM, createObserveGrads, incGeneration, generation, params, perturb := descend.NewDynamicSubDimSM(s.SubScope("sm"), updatesPH, randomExpand, size) // make the state machine.
	_ = params
	_ = generation

	loss := lossFunc(s.SubScope("loss"), params)

	grads := createObserveGrads(lossFunc, step)
	updates := op.Mul(s.SubScope("scale_grads"), grads, op.Const(s.SubScope("grads_scale"), gradsScale))
	// We are reusing the training data for test. This is bad practice.
	accuracyOP := makeFinalizeAccuracy(s.SubScope("accuracy"), params, readFFTs, readLabels)

	unflatten, _ := models.MakeUnflatten(simpleModelParamDefs)
	layerParams := unflatten(s.SubScope("unflatten"), params)

	graph, err := s.Finalize()
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
		return loadDatum(fileName, false)
	})
	loadClass("label/nongood", func(fileName string) error {
		return loadDatum(fileName, true)
	})

	_, err = sess.Run(nil, nil, []*tf.Operation{initFFTcache, initLabelsCache})
	if err != nil {
		panic(err)
	}

	_, err = sess.Run(nil, nil, []*tf.Operation{initOP})
	if err != nil {
		panic(err)
	}
	// now we can close the queue
	_, err = sess.Run(nil, nil, []*tf.Operation{closeQueue})
	if err != nil {
		panic(err)
	}

	err = initSM(sess)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10000; i++ {
		observedGrads, err := sess.Run(nil, []tf.Output{updates, loss}, nil)
		if err != nil {
			panic(err)
		}
		//fmt.Println(observedGrads[0].Value())
		fmt.Println("loss:", observedGrads[1].Value())
		_, err = sess.Run(map[tf.Output]*tf.Tensor{updatesPH: observedGrads[0]}, nil, []*tf.Operation{perturb})
		if err != nil {
			panic(err)
		}
		_, err = sess.Run(nil, nil, []*tf.Operation{incGeneration})
		if err != nil {
			panic(err)
		}
		if i%1 == 0 {
			acc, err := sess.Run(nil, []tf.Output{accuracyOP}, nil)
			if err != nil {
				panic(err)
			}
			fmt.Println(i, acc[0].Value().(float32)*100.0, "%")
		}
	}

	results, err := sess.Run(nil, layerParams, nil)
	if err != nil {
		panic(err)
	}

	s = op.NewScope()
	filter := op.Const(s.SubScope("filter"), results[0])
	weights := op.Const(s.SubScope("weights"), results[1])

	dataPH := op.Placeholder(s.SubScope("input"), tf.String, op.PlaceholderShape(tf.ScalarShape()))
	ffts := preprocessAudio(s.SubScope("preprocess"), dataPH)
	fmt.Println("ffts:", ffts.Shape())
	expandedFfts := op.ExpandDims(s, ffts, op.Const(s.SubScope("one"), int64(0)))
	fmt.Println("expandedFfts:", expandedFfts.Shape())
	output := simpleModel(s.SubScope("model"), expandedFfts, filter, weights)
	label := namedIdentity(s, op.Squeeze(s.SubScope("remove_dim"), output), "output")
	_ = label
	fmt.Println(output.Shape())
	graph, err = s.Finalize()
	if err != nil {
		panic(err)
	}

	file, err := os.Create("conv1.pb")
	if err != nil {
		panic(err)
	}
	_, err = graph.WriteTo(file)
	if err != nil {
		panic(err)
	}
}

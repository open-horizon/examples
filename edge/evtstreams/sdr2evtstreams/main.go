package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/open-horizon/examples/edge/evtstreams/sdr2evtstreams/audiolib"
	rtlsdr "github.com/open-horizon/examples/edge/services/sdr/rtlsdrclientlib"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/viert/lame"
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

// model holds the session, the input placeholder and output.
type model struct {
	Sess    *tf.Session
	InputPH tf.Output
	Output  tf.Output
}

// goodness takes a chunk of raw audio with no headers and returns a value between 0 and 1.
// 1 for good (in this case speech), 0 for nongood (in this case nonspeech).
// the audio must be exactly 32 seconds long.
func (m *model) goodness(audio []byte) (value float32, err error) {
	// first we must convert the audio to a string tensor.
	inputTensor, err := tf.NewTensor(string(audio))
	if err != nil {
		return
	}
	// then feed the input into the input placeholder while pulling on the output.
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

type evtstreamsConn struct {
	Producer sarama.SyncProducer
	Topic    string
}

// taken from cloud/sdr/data-ingest/example-go-clients/util/util.go
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

func connect(topic string) (conn evtstreamsConn, err error) {
	conn.Topic = topic
	apiKey := getEnv("EVTSTREAMS_API_KEY")
	username := "token"
	password := apiKey
	brokerStr := getEnv("EVTSTREAMS_BROKER_URL")
	brokers := strings.Split(brokerStr, ",")
	config := sarama.NewConfig()
	err = populateConfig(config, username, password, apiKey)
	if err != nil {
		return
	}
	fmt.Println("now connecting to evtstreams")
	conn.Producer, err = sarama.NewSyncProducer(brokers, config)
	fmt.Println("done trying to connect")
	if err != nil {
		return
	}
	return
}

func (conn *evtstreamsConn) publishAudio(audioMsg *audiolib.AudioMsg) (err error) {
	// as AudioMsg implements the sarama.Encoder interface, we can pass it directly to ProducerMessage.
	msg := &sarama.ProducerMessage{Topic: conn.Topic, Key: nil, Value: audioMsg}
	partition, offset, err := conn.Producer.SendMessage(msg)
	if err != nil {
		log.Printf("FAILED to send message: %s\n", err)
	} else {
		log.Printf("> message sent to partition %d at offset %d\n", partition, offset)
	}
	return
}

// read env vars from system with fall back.
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

// Copy pasted from github.com/open-horizon/examples/edge/services/gps/src/hgps to workaround package import issues.
type sourceType string

const (
	MANUAL    sourceType = "Manual"
	ESTIMATED sourceType = "Estimated"
	SEARCHING sourceType = "Searching"
	GPS       sourceType = "GPS"
)

// JSON struct for location data
type locationData struct {
	Latitude   float64    `json:"latitude" description:"Location latitude"`
	Longitude  float64    `json:"longitude" description:"Location longitude"`
	ElevationM float64    `json:"elevation" description:"Location elevation in meters"`
	AccuracyKM float64    `json:"accuracy_km" description:"Location accuracy in kilometers"`
	LocSource  sourceType `json:"loc_source" description:"Location source (one of: Manual, Estimated, GPS, or Searching)"`
	LastUpdate int64      `json:"loc_last_update" description:"Time of most recent location update (UTC)."`
}

func getGPS() (location locationData, err error) {
	resp, err := http.Get("http://" + gpshostname + ":8080/v1/gps/location")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = errors.New("bad resp")
		return
	}
	jsonByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonByte, &location)
	return
}

func rawToB64Mp3(rawBytes []byte) (b64Bytes string) {
	reader := bytes.NewReader(rawBytes)
	mp3Buff := bytes.Buffer{}

	wr := lame.NewWriter(&mp3Buff)
	wr.Encoder.SetBitrate(30)
	wr.Encoder.SetQuality(1)
	wr.Encoder.SetInSamplerate(16000)
	wr.Encoder.SetNumChannels(1)
	// IMPORTANT!
	wr.Encoder.InitParams()
	reader.WriteTo(wr)

	b64Buff := bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, &b64Buff)
	encoder.Write(mp3Buff.Bytes())
	encoder.Close()

	b64Bytes = string(b64Buff.Bytes())
	return
}

// the default hostname if not overridden
var hostname string = "ibm.sdr"
var gpshostname string = "ibm.gps"

func main() {
	alt_addr := os.Getenv("RTLSDR_ADDR")
	// if no alternative address is set, use the default.
	if alt_addr != "" {
		fmt.Println("connecting to remote rtlsdr:", alt_addr)
		hostname = alt_addr
	}
	gps_alt_addr := os.Getenv("GPS_ADDR")
	// if no alternative address is set, use the default.
	if gps_alt_addr != "" {
		fmt.Println("connecting to remote gps:", gps_alt_addr)
		gpshostname = gps_alt_addr
	}
	use_gps := os.Getenv("USE_GPS") != "false"
	if !use_gps {
		fmt.Println("not using GPS because USE_GPS=false")
	}
	verbose := os.Getenv("VERBOSE") == "1"
	if verbose {
		fmt.Println("verbose logging enabled")
	}
	devID := getEnv("HZN_ORG_ID", "HZN_ORGANIZATION") + "/" + getEnv("HZN_DEVICE_ID")
	// load the graph def from FS
	m, err := newModel("model.pb")
	if err != nil {
		panic(err)
	}
	fmt.Println("model loaded")
	topic := getEnv("EVTSTREAMS_TOPIC")
	fmt.Printf("using topic %s\n", topic)
	conn, err := connect(topic)
	if err != nil {
		panic(err)
	}
	fmt.Println("connected to evtstreams")
	// create a map to hold the goodness for each station we have ever oberved.
	// This map will grow as long as the program lives
	stationGoodness := map[float32]float32{}
	lastStationsRefresh := time.Time{}

	// make it fail sooner.
	if use_gps {
		_, err = getGPS()
		if err != nil {
			fmt.Println(err)
			panic("can't get location from GPS")
		}
	}
	sdr_origin := ""
	var hasCapturedFirstClip = false
	var hasSentFirstClip = false
	for {
		// if it has been over 5 minuts since we last updated the list of strong stations,
		if time.Now().Sub(lastStationsRefresh) > (5 * time.Minute) {
			fmt.Println("fetching new list of stations")
			// for ever, we aquire a list of stations,
			freqs, err := rtlsdr.GetFreqs(hostname)
			if err != nil {
				panic(err)
			}
			fmt.Println("got", len(freqs.Freqs), "freqs from sdr")
			sdr_origin = freqs.Origin
			for _, station := range freqs.Freqs {
				_, prs := stationGoodness[station]
				if !prs {
					// only if the station is not already in our map, do we add it, with an initial value of 0.5
					fmt.Println("found new station: ", station)
					stationGoodness[station] = 0.5
				}
			}
			// if no stations can be found, we can't do anything, so panic.
			if len(stationGoodness) < 1 {
				panic("No FM stations. Move the antenna?")
			}
			fmt.Println("found", len(freqs.Freqs), "stations from", freqs.Origin)
			fmt.Println(stationGoodness)
			lastStationsRefresh = time.Now()
		}
		for station, goodness := range stationGoodness {
			// if our goodness is less then a random number between 0 and 1.
			if rand.Float32() < goodness {
				audio, err := rtlsdr.GetAudio(hostname, int(station))
				if err != nil {
					panic(err)
				}
				if !hasCapturedFirstClip {
					fmt.Println("Captured first clip")
					hasCapturedFirstClip = true
				}
				val, err := m.goodness(audio)
				if err != nil {
					panic(err)
				}
				// if the value is close to 1, the goodness of that station will increase, if the value is small, the goodness will decrease.
				stationGoodness[station] = stationGoodness[station]*(val+0.3) + 0.05
				if verbose {
					fmt.Println(station, "observed value:", val, "updated goodness:", stationGoodness[station])
				}
				// if the value is over 0.5, it is worth sending to the cloud.
				if val > 0.5 {
					var location = locationData{}
					if use_gps {
						location, err = getGPS()
						if err != nil {
							fmt.Println(err)
							continue
						}
					}
					// construct the message,
					msg := &audiolib.AudioMsg{
						Audio:         rawToB64Mp3(audio),
						Ts:            time.Now().Unix(),
						Freq:          station,
						ExpectedValue: val,
						DevID:         devID,
						Lat:           float32(location.Latitude),
						Lon:           float32(location.Longitude),
						ContentType:   "audio/mpeg",
						Origin:        sdr_origin,
					}
					// and publish it to evtstreams
					err = conn.publishAudio(msg)
					if err != nil {
						fmt.Println(err)
					}
					if !hasSentFirstClip {
						fmt.Println("Sent first clip")
						hasSentFirstClip = true
					}
				} else {
					if verbose {
						fmt.Println("Not sending sample from", station, "becouse value is", val)
					}
				}
			}
		}
	}
}

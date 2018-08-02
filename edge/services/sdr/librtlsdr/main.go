package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	rtlsdr "github.com/open-horizon/examples/edge/services/sdr/librtlsdr/rtlsdrclientlib"
)

func captureAudio(freq int) (audio []byte, err error) {
	cmd := exec.Command("rtl_fm", "-M", "fm", "-s", "170k", "-o", "4", "-A", "fast", "-r", "16k", "-l", "0", "-E", "deemp", "-f", strconv.Itoa(freq))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(30 * time.Second)
	err = cmd.Process.Kill()
	if err != nil {
		err = errors.New(string(stderr.Bytes()))
		return
	}
	audio = stdout.Bytes()
	if len(audio) < 900000 {
		err = errors.New("for some reason, audio is too short")
	}
	// if the audio is too long, trim it.
	if len(audio) > 938496 {
		audio = audio[:938496]
	}
	if len(audio) < 938496 {
		audio = append(audio, make([]byte, 938496-len(audio))...)
	}
	return
}

const ROWS int = 18
const COLS int = 411

func stringListToFloat(stringList []string) (floatList []float32) {
	for _, val := range stringList {
		num, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
		if err != nil {
			panic(err)
		}
		floatList = append(floatList, float32(num))
	}
	return
}

func capturePower() (power rtlsdr.PowerDist, err error) {
	start := 70000000
	end := 110000000
	power.Low = float32(start)
	power.High = float32(end)
	cmd := exec.Command("rtl_power", "-e", "10", "-c", "20%", "-f", strconv.Itoa(start)+":"+strconv.Itoa(end)+":10000")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = errors.New(string(stderr.Bytes()))
		return
	}
	r := csv.NewReader(bytes.NewReader(stdout.Bytes()))
	recordList, err := r.ReadAll()
	if err != nil {
		return
	}
	if len(recordList) != ROWS {
		err = errors.New("expected " + strconv.Itoa(ROWS) + " rows, got " + strconv.Itoa(len(recordList)) + " rows")
		return
	}
	for _, row := range recordList {
		if len(row[6:]) != COLS {
			err = errors.New("expected " + strconv.Itoa(COLS) + " elems, got " + strconv.Itoa(len(row[6:])) + " elems")
			return
		}
		power.Dbm = append(power.Dbm, stringListToFloat(row[6:])...)
	}
	return
}

func audioHandler(w http.ResponseWriter, r *http.Request) {
	freq, err := strconv.Atoi(r.URL.Path[7:])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	audio, err := captureAudio(freq)
	if (err != nil) && !(os.Getenv("MOCK_IF_YOU_MUST") == "false") {
		fmt.Println("using mock audio")
		audio, err = ioutil.ReadFile("mock_audio.raw")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.Write(audio)
}

func powerHandler(w http.ResponseWriter, r *http.Request) {
	power, err := capturePower()
	if (err != nil) && !(os.Getenv("MOCK_IF_YOU_MUST") == "false") {
		fmt.Println("using mock power data:", err.Error())
		err = nil
		power = rtlsdr.PowerDist{
			Low:  float32(70000000),
			High: float32(110000000),
			Dbm:  make([]float32, ROWS*COLS),
		}
	}
	for i := range power.Dbm {
		if math.IsNaN(float64(power.Dbm[i])) {
			power.Dbm[i] = -1234
		}
	}
	jsonBytes, err := json.Marshal(power)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}

func main() {
	fmt.Println("starting sdr daemon")
	http.HandleFunc("/audio/", audioHandler)
	http.HandleFunc("/power", powerHandler)
	log.Fatal(http.ListenAndServe(":5427", nil))
}

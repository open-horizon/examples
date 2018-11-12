package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/open-horizon/examples/edge/services/sdr/bbcfake"
	rtlsdr "github.com/open-horizon/examples/edge/services/sdr/rtlsdrclientlib"
)

func captureAudio(freq int) (audio []byte, err error) {
	cmd := exec.Command("rtl_fm", "-M", "fm", "-s", "170k", "-o", "4", "-A", "fast", "-r", "16k", "-l", "0", "-E", "deemp", "-f", strconv.Itoa(freq))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = append(cmd.Env, "RTLSDR_RPC_IS_ENABLED=1", "RTLSDR_RPC_SERV_ADDR=localhost")
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
	power.Origin = "sdr_hardware"
	power.Low = float32(start)
	power.High = float32(end)
	// rtl_power -e 10 -c 20% -f 70000000:110000000:10000
	cmd := exec.Command("rtl_power", "-e", "10", "-c", "20%", "-f", strconv.Itoa(start)+":"+strconv.Itoa(end)+":10000")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = append(cmd.Env, "RTLSDR_RPC_IS_ENABLED=1", "RTLSDR_RPC_SERV_ADDR=localhost")
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

func makeAudioHandler(fake *bbcfake.FakeRadio) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		freq, err := strconv.Atoi(r.URL.Path[7:])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var audio []byte
		if freq == 0 {
			audio = fake.GetNextChunk()
			time.Sleep(30 * time.Second)
		} else {
			if rtlRpcdIsAlive {
				audio, err = captureAudio(freq)
			} else {
				err = errors.New("freq != 0 but rtl_rpcd is dead")
			}
		}
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(audio[:938496])
	}
}

func powerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/power is deprecated! Please use /freqs")
	power, err := capturePower()
	if (err != nil) && !(os.Getenv("MOCK_IF_YOU_MUST") == "false") {
		fmt.Println("using mock power data:", err.Error())
		err = nil
		power = rtlsdr.PowerDist{
			Origin: "mock_file",
			Low:    float32(70000000),
			High:   float32(110000000),
			Dbm:    make([]float32, ROWS*COLS),
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

func FreqToIndex(freq float32, data rtlsdr.PowerDist) int {
	percentPos := (freq - data.Low) / (data.High - data.Low)
	index := int(float32(len(data.Dbm)) * percentPos)
	return index
}

func getCeilingSignals(celling float32) (freqs []float32, origin string) {
	fmt.Println("begin capture power")
	data, err := capturePower()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("sending fake freq")
		origin = "fake"
		freqs = []float32{0.0}
		return
	}
	for i := range data.Dbm {
		if math.IsNaN(float64(data.Dbm[i])) {
			data.Dbm[i] = -1234
		}
	}

	for i := float32(85900000); i < data.High; i += 200000 {
		dbm := data.Dbm[FreqToIndex(i, data)]
		if dbm > celling {
			freqs = append(freqs, i)
		}
	}
	origin = "sdr_hardware"
	return
}

func freqsHandler(w http.ResponseWriter, r *http.Request) {
	freqs, origin := getCeilingSignals(-8)
	jsonBytes, err := json.Marshal(rtlsdr.Freqs{Origin: origin, Freqs: freqs})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}

var rtlRpcdIsAlive = true

func startRtlrpcd(onFail func(error)) {
	cmd := exec.Command("rtl_rpcd")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	startTime := time.Now()
	fmt.Println("starting rtl_rpcd")
	err := cmd.Run()
	fmt.Println("run err:", err)
	rtlRpcdIsAlive = false
	fmt.Println("rtl_rpcd stopped after", startTime.Sub(time.Now()), "seconds")
	onFail(err)
}

func main() {
	fake := bbcfake.NewFakeRadio()
	go startRtlrpcd(func(err error) {
		fmt.Println("rtl_rpcd has died:", err.Error())
		fmt.Println("Requests for audio will now be fulfilled using fake data.")
	})
	http.HandleFunc("/audio/", makeAudioHandler(&fake))
	http.HandleFunc("/power", powerHandler)
	http.HandleFunc("/freqs", freqsHandler)
	log.Fatal(http.ListenAndServe(":5427", nil))
}

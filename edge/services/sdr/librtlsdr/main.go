package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func captureAudio(freq int) (audio []byte, err error) {
	cmd := exec.Command("rtl_fm", "-M", "fm", "-s", "170k", "-o", "4", "-A", "fast", "-r", "16k", "-l", "0", "-E", "deemp", "-f", strconv.Itoa(freq))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	fmt.Println("starting command")
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(30 * time.Second)
	err = cmd.Process.Kill()
	if err != nil {
		panic(err)
	}
	audio = stdout.Bytes()
	//errStr := string(stderr.Bytes())
	//fmt.Println(errStr)
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

func capturePower() (power PowerDist, err error) {
	start := 70000000
	end := 110000000
	power.Low = float32(start)
	power.High = float32(end)
	cmd := exec.Command("rtl_power", "-e", "10", "-c", "20%", "-f", strconv.Itoa(start)+":"+strconv.Itoa(end)+":10000")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	fmt.Println("starting command")
	err = cmd.Run()
	if err != nil {
		panic(err)
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
	//fmt.Println(recordList)
	return
}

// PowerDist is the distribution of power of frequency.
type PowerDist struct {
	Low  float32   `json:"low"`
	High float32   `json:"high"`
	Dbm  []float32 `json:"dbm"`
}

func audioHandler(w http.ResponseWriter, r *http.Request) {
	freq, err := strconv.Atoi(r.URL.Path[7:])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	audio, err := captureAudio(freq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(audio)
}

func powerHandler(w http.ResponseWriter, r *http.Request) {
	power, err := capturePower()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonBytes, err := json.Marshal(power)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonBytes)
}

func main() {
	http.HandleFunc("/audio/", audioHandler)
	http.HandleFunc("/power", powerHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//78

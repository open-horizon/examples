package main

import (
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/open-horizon/examples/cloud/sdr/sdr_data_processing/watson/stt"
)

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
	if len(audio) < 100 {
		panic("audio is too short")
	}
	return
}

// FreqToIndex converts a frequnecy to a list index.
func FreqToIndex(freq float32, data PowerDist) int {
	percentPos := (freq - data.Low) / (data.High - data.Low)
	index := int(float32(len(data.Dbm)) * percentPos)
	return index
}

func GetCeilingSignals(data PowerDist, celling float32) (stationFreqs []float32) {
	for i := float32(85900000); i < data.High; i += 200000 {
		dbm := data.Dbm[FreqToIndex(i, data)]
		if dbm > celling && dbm != 0 {
			stationFreqs = append(stationFreqs, i)
		}
	}
	return
}

// PowerDist is the distribution of power of frequency.
type PowerDist struct {
	Low  float32   `json:"low"`
	High float32   `json:"high"`
	Dbm  []float32 `json:"dbm"`
}

func getPower() (power PowerDist, err error) {
	resp, err := http.Get("http://localhost:8080/power")
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
	err = json.Unmarshal(jsonByte, &power)
	return
}

func totalText(transcript stt.TranscribeResponse) (sum int) {
	if len(transcript.Results) == 0 {
		return
	}
	for _, result := range transcript.Results {
		for _, alt := range result.Alternatives {
			sum += len(alt.Transcript)
		}
	}
	return
}

func main() {
	username := os.Getenv("STT_USERNAME")
	if username == "" {
		panic("STT_USERNAME not set")
	}
	password := os.Getenv("STT_PASSWORD")
	if password == "" {
		panic("STT_USERNAME not set")
	}
	var i = 0
	for {
		power, err := getPower()
		if err != nil {
			panic(err)
		}
		stations := GetCeilingSignals(power, -13)
		for _, station := range stations {
			fmt.Println("starting freq", station)
			audio, err := getAudio(int(station))
			if err != nil {
				panic(err)
			}
			transcript, err := stt.Transcribe(audio, username, password)
			if err != nil {
				panic(err)
			}
			fmt.Println(totalText(transcript), transcript)
			hash := sha256.Sum256(audio)
			name := base32.StdEncoding.EncodeToString(hash[:])
			if totalText(transcript) > 20 {
				err = ioutil.WriteFile("good/"+name+".raw", audio, 0644)
			} else {
				err = ioutil.WriteFile("nongood/"+name+".raw", audio, 0644)
			}
			if err != nil {
				panic(err)
			}
			i++
		}
	}
}

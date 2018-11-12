package rtlsdr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// GetAudio fetches a 30 second chunk of raw audio.
func GetAudio(hostname string, freq int) (audio []byte, err error) {
	resp, err := http.Get("http://" + hostname + ":5427/audio/" + strconv.Itoa(freq))
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

// GetCeilingSignals fetches the signal power distribution and samples it down to a list of frequencies at which there (probably) exist strong FM signals.
func GetCeilingSignals(hostname string, celling float32) (stationFreqs []float32, origin string, err error) {
	data, err := getPower(hostname)
	if err != nil {
		return
	}
	for i := float32(85900000); i < data.High; i += 200000 {
		dbm := data.Dbm[FreqToIndex(i, data)]
		if dbm > celling {
			stationFreqs = append(stationFreqs, i)
		}
	}
	origin = data.Origin
	return
}

// Freqs stores a list of frequencies of stations
type Freqs struct {
	Origin string    `json:"origin"`
	Freqs  []float32 `json:"freqs"`
}

// PowerDist is the distribution of power of frequency.
type PowerDist struct {
	Origin string    `json:"origin"`
	Low    float32   `json:"low"`
	High   float32   `json:"high"`
	Dbm    []float32 `json:"dbm"`
}

func GetFreqs(hostname string) (freqs Freqs, err error) {
	timeout := time.Duration(40 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get("http://" + hostname + ":5427/freqs")
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
	err = json.Unmarshal(jsonByte, &freqs)
	fmt.Println("done with freq")
	return
}

func getPower(hostname string) (power PowerDist, err error) {
	resp, err := http.Get("http://" + hostname + ":5427/power")
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

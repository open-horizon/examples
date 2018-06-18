package rtlsdr

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

// GetAudio fetches a 30 second chunk of raw audio.
func GetAudio(hostname string, freq int) (audio []byte, err error) {
	resp, err := http.Get("http://" + hostname + ":8080/audio/" + strconv.Itoa(freq))
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
func GetCeilingSignals(hostname string, celling float32) (stationFreqs []float32, err error) {
	data, err := getPower(hostname)
	if err != nil {
		return
	}
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

func getPower(hostname string) (power PowerDist, err error) {
	resp, err := http.Get("http://" + hostname + ":8080/power")
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

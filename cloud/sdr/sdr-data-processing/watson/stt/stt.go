package stt

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type TranscribeResponse struct {
	Results []Result `json:"results"`
	Index   int      `json:"results_index"`
}

type Result struct {
	Alternatives []Alternative `json:"alternatives"`
	Final        bool          `json:"final"`
}

type Alternative struct {
	Confidence float32 `json:"confidence"`
	Transcript string  `json:"transcript"`
}

// AppendWAVheader for a 16k s16le wav file
// This is ugly.
// Don't do it if you have a better way.
func appendWAVheader(rawAudio []byte) (wavAudio []byte) {
	hexHeader := "5249464646520e0057415645666d74201000000001000100803e0000007d0000020010004c4953541a000000494e464f495346540e0000004c61766635382e31322e313030006461746100520e00"
	header, err := hex.DecodeString(hexHeader)
	if err != nil {
		panic("bad hex")
	}
	wavAudio = append(header, rawAudio...)
	ioutil.WriteFile("demo.wav", wavAudio, 0644)
	return
}

// Transcribe a chunk of raw audio
func Transcribe(rawAudio []byte, username, password string) (response TranscribeResponse, err error) {
	wavAudio := appendWAVheader(rawAudio)
	apiURL := "https://stream.watsonplatform.net/speech-to-text/api/v1/recognize"

	u, _ := url.ParseRequestURI(apiURL)
	urlStr := fmt.Sprintf("%v", u)
	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewReader(wavAudio))
	r.SetBasicAuth(username, password)
	r.Header.Add("Content-Type", "audio/wav")
	resp, err := client.Do(r)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("got status" + strconv.Itoa(resp.StatusCode))
		return
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(bytes, &response); err != nil {
		return
	}
	return
}

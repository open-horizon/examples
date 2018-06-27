package main

import (
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/open-horizon/examples/cloud/sdr/sdr-data-processing/watson/stt"
	rtlsdr "github.com/open-horizon/examples/edge/services/sdr/librtlsdr/rtlsdrclientlib"
)

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
		panic("STT_PASSWORD not set")
	}
	var i = 0
	for {
		stations, err := rtlsdr.GetCeilingSignals("localhost", -5)
		if err != nil {
			panic(err)
		}
		fmt.Println("found", len(stations), "stations")
		for _, station := range stations {
			fmt.Println("starting freq", station)
			audio, err := rtlsdr.GetAudio("localhost", int(station))
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

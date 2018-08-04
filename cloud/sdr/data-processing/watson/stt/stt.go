package stt

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/open-horizon/examples/cloud/sdr/data-ingest/example-go-clients/util"
	"github.com/open-horizon/examples/cloud/sdr/data-processing/wutil"
)

// TranscribeResponse is the top level struct which Watson speech to text gives us.
type TranscribeResponse struct {
	Results []Result `json:"results"`
	Index   int      `json:"results_index"`
}

// Result is just a list of Alternatives with a final bool. For non streaming, we can probaly ignore Final.
type Result struct {
	Alternatives []Alternative `json:"alternatives"`
	Final        bool          `json:"final"`
}

// Alternative holds the actual text along with a Confidence
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
	wavFileName := "/tmp/stt-demo.wav"
	util.Verbose("converting raw audio to a wav file, and storing in %s\n", wavFileName)
	wavAudio = append(header, rawAudio...)
	ioutil.WriteFile(wavFileName, wavAudio, 0644)
	return
}

// Transcribe a chunk of raw audio
// takes raw auidio without a header
func Transcribe(rawAudio []byte, username, password string) (response TranscribeResponse, err error) {
	// we need to add a header so that watson will know the specs of the audio.
	wavAudio := appendWAVheader(rawAudio)
	fmt.Println("using Watson STT to convert the audio to text...")
	// Watson STT API: https://www.ibm.com/watson/developercloud/speech-to-text/api/v1/curl.html?curl#recognize-sessionless
	apiURL := "https://stream.watsonplatform.net/speech-to-text/api/v1/recognize"
	headers := []wutil.Header{{Key: "Content-Type", Value: "audio/wav"}}
	err = wutil.HTTPPost(apiURL, username, password, headers, bytes.NewReader(wavAudio), &response)
	return
}

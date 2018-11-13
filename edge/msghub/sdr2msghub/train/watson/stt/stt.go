package stt

import (
	"bytes"
	"fmt"

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

// Transcribe a chunk of audio
func Transcribe(audioBytes []byte, contentType string, username, password string) (response TranscribeResponse, err error) {
	fmt.Println("using Watson STT to convert the audio to text...")
	// Watson STT API: https://www.ibm.com/watson/developercloud/speech-to-text/api/v1/curl.html?curl#recognize-sessionless
	apiURL := "https://stream.watsonplatform.net/speech-to-text/api/v1/recognize"
	headers := []wutil.Header{{Key: "Content-Type", Value: contentType}}
	err = wutil.HTTPPost(apiURL, username, password, headers, bytes.NewReader(audioBytes), &response)
	return
}

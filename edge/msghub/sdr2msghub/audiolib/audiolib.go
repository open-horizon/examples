package audiolib

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	"github.com/viert/lame"
)

// AudioMsg holds the metadata and audio that we send to IBM Message Hub
type AudioMsg struct {
	Audio         string  `json:"audio"`
	Ts            int64   `json:"ts"`
	Freq          float32 `json:"freq"`
	ExpectedValue float32 `json:"expectedValue"`
	DevID         string  `json:"devID"`
	Lat           float32 `json:"lat"`
	Lon           float32 `json:"lon"`
	ContentType   string  `json:"contentType"`
}

// Encode implemented for the https://godoc.org/github.com/Shopify/sarama#Encoder interface
func (msg *AudioMsg) Encode() (serialized []byte, err error) {
	serialized, err = json.Marshal(msg)
	return
}

// Length implemented for the https://godoc.org/github.com/Shopify/sarama#Encoder interface
// This is an ugly hack becouse I can't easily calculate the length without actualy serializing the object.
func (msg *AudioMsg) Length() int {
	serialized, _ := msg.Encode()
	return len(serialized)
}

func RawToB64Mp3(rawBytes []byte) (b64Bytes string) {
	reader := bytes.NewReader(rawBytes)
	mp3Buff := bytes.Buffer{}

	wr := lame.NewWriter(&mp3Buff)
	wr.Encoder.SetBitrate(30)
	wr.Encoder.SetQuality(1)
	wr.Encoder.SetInSamplerate(16000)
	wr.Encoder.SetNumChannels(1)
	// IMPORTANT!
	wr.Encoder.InitParams()
	reader.WriteTo(wr)

	b64Buff := bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, &b64Buff)
	encoder.Write(mp3Buff.Bytes())
	encoder.Close()

	b64Bytes = string(b64Buff.Bytes())
	return
}

package audiolib

import (
	"encoding/json"
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

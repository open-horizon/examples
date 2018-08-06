package audiolib

import (
	"bytes"
	"encoding/gob"
	"time"
)

// AudioMsg holds a 30 second chunk of raw audio and metadata
type AudioMsg struct {
	Audio         []byte // the chunk of raw audio. No headers
	Ts            time.Time
	Freq          float32 // the frequancy on the FM spectrum at which it was captured.
	ExpectedValue float32 // the goodness of the clip, between 0 and 1, and no less then the threshold, currently 0.5
	Lat           float32
	Lon           float32
	DevID         string // id of the device which captured the audio.
}

// Encode implemented for the https://godoc.org/github.com/Shopify/sarama#Encoder interface
func (msg *AudioMsg) Encode() (serialized []byte, err error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	if err != nil {
		panic(err)
	}
	err = enc.Encode(msg)
	if err != nil {
		panic(err)
	}
	serialized = buff.Bytes()
	return
}

// Length implemented for the https://godoc.org/github.com/Shopify/sarama#Encoder interface
// This is an ugly hack becouse I can't easily calculate the length without actualy serializing the object.
func (msg *AudioMsg) Length() int {
	serialized, _ := msg.Encode()
	return len(serialized)
}

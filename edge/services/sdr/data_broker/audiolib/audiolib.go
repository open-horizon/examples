package audiolib

import (
	"bytes"
	"encoding/gob"
	"time"
)

// AudioMsg holds a 30 second chunk of raw audio and metadata
type AudioMsg struct {
	Audio         []byte
	Ts            time.Time
	Freq          float32
	ExpectedValue float32
	Lat           float32
	Lon           float32
	DevID         string
}

// Encode implemented for sarama
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

// Length implemented for sarama
func (msg *AudioMsg) Length() int {
	serialized, _ := msg.Encode()
	return len(serialized)
}

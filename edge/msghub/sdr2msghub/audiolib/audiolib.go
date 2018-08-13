package audiolib

import "github.com/gogo/protobuf/proto"

// Encode implemented for the https://godoc.org/github.com/Shopify/sarama#Encoder interface
func (msg *AudioMsg) Encode() (serialized []byte, err error) {
	serialized, err = proto.Marshal(msg)
	return
}

// Length implemented for the https://godoc.org/github.com/Shopify/sarama#Encoder interface
// This is an ugly hack becouse I can't easily calculate the length without actualy serializing the object.
func (msg *AudioMsg) Length() int {
	serialized, _ := msg.Encode()
	return len(serialized)
}

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type intentMsg struct {
	Name string    `json:"name"`
	Prob float32   `json:"prob"`
	TS   time.Time `json:"ts"`
}

func main() {
	// connect to the audio stream
	conn, err := net.Dial("tcp", "aural2:49610")
	fmt.Println("failed to connect to aural2:49610, trying localhost")
	if err != nil {
		conn, err = net.Dial("tcp", "localhost:49610")
		if err != nil {
			panic(err)
		}
	}
	var intent intentMsg
	decoder := json.NewDecoder(conn)
	for {
		err = decoder.Decode(&intent)
		if err != nil {
			panic(err)
		}
		fmt.Println(intent)
		if intent.Name == "play0.5" {
			// pre fetch something, but don't actually do anything.
		}
		if intent.Name == "play0.9" {
			// actually do something
		}
	}
}

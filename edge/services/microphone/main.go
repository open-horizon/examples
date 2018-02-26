package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

// try to connect to an audio device.
func tryToConnect(hardwareNum string) (reader io.Reader, err error) {
	devName := "plughw:" + hardwareNum
	cmd := exec.Command("arecord", "-D", devName, "-r", "16000", "-t", "raw", "-f", "S16_LE", "-c", "1", "-")
	cmd.Stderr = os.Stdout
	reader, err = cmd.StdoutPipe()
	if err != nil {
		fmt.Println("can't open stdout on " + devName + ": " + err.Error())
		cmd.Process.Kill()
		return
	}
	if err = cmd.Start(); err != nil {
		fmt.Println("arecord failed on "+devName+" with error: "+err.Error(), "trying another device")
		cmd.Process.Kill()
		return
	}
	buff := make([]byte, 2)
	i, err := reader.Read(buff)
	if err != nil {
		cmd.Process.Kill()
		return
	}
	if i != 2 {
		err = errors.New("can't read a sample")
		cmd.Process.Kill()
		return
	}
	return
}

var version string

func main() {
	fmt.Println("starting microphone ms version", version)
	args := os.Args
	var reader io.Reader
	var err error
	if len(args) == 2 {
		reader, err = tryToConnect(args[1])
	} else {
		for _, i := range []int{5, 4, 3, 2, 1, 0} {
			fmt.Println("reading data from device", i)
			reader, err = tryToConnect(strconv.Itoa(i))
			if err != nil {
				fmt.Println("can't read from device", i, "trying next device")
				continue
			}
			break
		}
	}
	if reader == nil {
		fmt.Println("fail to find device from which arecord could get audio")
		os.Exit(1)
	}
	buffsMap := map[int]chan []byte{}
	var buffsIndex int
	buffsMutex := sync.Mutex{}
	go func() {
		buff := make([]byte, 4000)
		for {
			i, err := reader.Read(buff)
			if err != nil {
				panic(err)
			}
			_ = i
			buffsMutex.Lock()
			for i, connChan := range buffsMap {
				if len(connChan) < 99 {
					connChan <- buff
				} else {
					fmt.Println("conn", i, "is too slow, dropping samples")
				}
			}
			buffsMutex.Unlock()
		}
	}()
	listenAddr := "microphone:48926"
	fmt.Println("attempting to listen on", listenAddr)
	l, err := net.Listen("tcp", listenAddr)
	fmt.Println("done trying to resolve addr")
	if err != nil {
		fmt.Println("Error listening on", listenAddr, err.Error(), ", trying localhost")
		listenAddr = "localhost:48926"
		l, err = net.Listen("tcp", listenAddr)
		if err != nil {
			fmt.Println("Error listening:", err.Error())
			os.Exit(1)
		}
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening for tcp connections on", listenAddr)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
		} else {
			go func() {
				buffsMutex.Lock()
				buffsIndex++
				index := buffsIndex
				buffsMap[buffsIndex] = make(chan []byte, 100)
				buffsMutex.Unlock()
				fmt.Println("starting to transfer samples")
				for {
					buffsMutex.Lock()
					connChan := buffsMap[index]
					buffsMutex.Unlock()
					samples := <-connChan
					if err != nil {
						fmt.Println(err)
						break
					}
					_, err = conn.Write(samples)
					if err != nil {
						fmt.Println(err)
						break
					}
				}
				conn.Close()
				buffsMutex.Lock()
				fmt.Println("closing connection", index)
				delete(buffsMap, index)
				buffsMutex.Unlock()
			}()
		}
	}
}

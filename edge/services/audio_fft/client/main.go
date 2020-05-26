// Copyright (c) 2020 SoftServe Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package main

import (
	"encoding/binary"
	"go-fft-client/shared"
	"math"
	"syscall"

	"github.com/gordonklaus/portaudio"
	"github.com/sirupsen/logrus"

	"os"
	"os/signal"
)

func main() {
	config := NewConfig()

	l, err := logrus.ParseLevel(config.LogLevel)
	if err == nil {
		logrus.SetLevel(l)
	}

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	c := shared.NewMQTTClient(config.MQTT)

	err = portaudio.Initialize()
	if err != nil {
		logrus.Fatal(err)
	}
	defer portaudio.Terminate()

	defaultDevice, err := portaudio.DefaultInputDevice()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Debugf("Suggested rate is %f", defaultDevice.DefaultSampleRate)

	bufferLen := config.RecordFrame * config.SampleRate
	buffer := make([]byte, 4*bufferLen)

	stream, err := portaudio.OpenDefaultStream(
		1, 0, float64(config.SampleRate), bufferLen, func(in []float32) {
			for i, x := range in {
				binary.BigEndian.PutUint32(buffer[i*4:], math.Float32bits(x))
			}

			logrus.Debug("Sample is ready")

			c.SendSample(buffer)
		})

	if err != nil {
		logrus.Fatal(err)
	}

	err = stream.Start()
	if err != nil {
		logrus.Fatal(err)
	}

	for range sig {
		c.Close()
		logrus.Info("Aborting sio")
		stream.Abort()
		stream.Close()
		os.Exit(0)
	}
}

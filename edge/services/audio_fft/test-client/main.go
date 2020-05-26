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
	"errors"
	"fmt"
	"go-fft-client/shared"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/zenwerk/go-wave"
)

func main() {
	ch := make(chan *shared.Response, 1)

	config := NewConfig()
	logrus.SetLevel(logrus.InfoLevel)

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: true,
		FullTimestamp:             true,
	})

	c := shared.NewMQTTClient(config.MQTT)
	c.Subscribe(config.MQTT.ResultTopic, func(bytes []byte) {
		r, err := shared.FromJson(bytes)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Infof("Triggered: %t", r.Trigger)
		ch <- r
	})

	folders, err := ioutil.ReadDir("./sets")

	if err != nil {
		logrus.Fatal(err)
	}

	type sample struct {
		name   string
		sample []byte
	}

	errorsCount := 0

	for _, f := range folders {
		if !f.IsDir() {
			continue
		}

		logrus.Infof("Using %s sets", f.Name())

		files, err := ioutil.ReadDir(fmt.Sprintf("./sets/%s/", f.Name()))
		if err != nil {
			logrus.Fatal(err)
		}

		samples := make([]*sample, 0)

		for _, s := range files {
			if s.IsDir() {
				continue
			}

			if filepath.Ext(s.Name()) != ".wav" {
				logrus.Warnf("Non-wav file found %s", s.Name())
				continue
			}

			fSample, err := getSample(fmt.Sprintf("./sets/%s/%s", f.Name(), s.Name()))

			if err != nil {
				logrus.Error("Failed to read sample", err)
				continue
			}

			samples = append(samples, &sample{
				name:   s.Name(),
				sample: fSample,
			})
			logrus.Infof("Found sample %s", s.Name())
		}

		lastSend := -1
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))

		for ii := 0; ii < 10; ii++ {
			index := rng.Intn(len(samples))
			shouldTrigger := false
			if index != lastSend {
				shouldTrigger = true
			}

			lastSend = index
			logrus.Infof("Sending sample %s, should trigger: %t", samples[index].name, shouldTrigger)
			err = sendSample(c, samples[index].sample, shouldTrigger, ch)
			if err != nil {
				logrus.Error("Failed to send sample", err)
				errorsCount++
			}
		}
	}

	if 0 == errorsCount {
		logrus.Info("All passed")
	} else {
		logrus.Errorf("Got %d errors", errorsCount)
	}
}

func getSample(fName string) ([]byte, error) {
	goodReader, err := wave.NewReader(fName)

	if err != nil {
		return nil, err
	}

	buf := make([]byte, 0)

	for {
		b, err := goodReader.ReadSample()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		buffer := make([]byte, 4*len(b))

		for i, x := range b {
			binary.BigEndian.PutUint32(buffer[i*4:], math.Float32bits(float32(x)))
		}

		buf = append(buf, buffer...)
	}

	return buf, nil
}

func sendSample(client shared.MQTTClient, sample []byte, expectedResult bool, ch chan *shared.Response) error {
	err := client.SendSample(sample)
	if err != nil {
		logrus.Fatal(err)
	}

	select {
	case r := <-ch:
		if r.Trigger != expectedResult {
			return errors.New("got unexpected trigger response")
		}
		break

	case <-time.After(5 * time.Second):
		return errors.New("got timeout while waiting on response")
	}

	return nil
}

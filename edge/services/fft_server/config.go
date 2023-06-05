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
	"go-fft-client/shared"

	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel   string `short:"l" long:"log_level" desction:"Log level" default:"info"`
	SampleRate int    `short:"r" long:"sample_rate" description:"Recording sample rate" default:"48000"`
	MQTT       *shared.MQTTConfig

	Analyzer struct {
		NFFT           int     `short:"n" long:"nfft" description:"Number of data points for FFT" default:"8192"`
		PeaksLimit     int     `short:"i" long:"peaks_limit" description:"Peaks limit" default:"2"`
		PeaksThreshold float64 `short:"e" long:"peaks_threshold" description:"Peaks threshold" default:"0.25"`
		FreqsThreshold float64 `short:"f" long:"freqs_threshold" description:"Frequencies threshold" default:"0.5"`
	}
}

func NewConfig() *Config {
	config := new(Config)
	_, err := flags.Parse(config)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	return config
}

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
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-fft-client/shared"

	"github.com/sirupsen/logrus"
)

type Server struct {
	config   *Config
	analyzer *Analyzer
}

func NewServer(config *Config, analyzer *Analyzer) *Server {
	return &Server{
		config:   config,
		analyzer: analyzer,
	}
}

func (s *Server) Process(req *SampleRequest) *shared.Response {
	r := shared.NewResponse(req.Timestamp)
	r.Trigger = s.analyzer.Analyze(req)
	return r
}

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

	analyzer := NewAnalyzer(config)
	server := NewServer(config, analyzer)

	c := shared.NewMQTTClient(config.MQTT)
	c.Subscribe(config.MQTT.RequestsTopic, func(in []byte) {
		req, err := NewRequest(in)
		if err != nil {
			logrus.Warn(err)
			return
		}

		go func() {
			response := server.Process(req)
			if response != nil {
				if response.Trigger {
					logrus.Infof("Triggered at %s", response.Timestamp.Format(time.Kitchen))
				}
				err := c.SendResults(response)
				if err != nil {
					logrus.Warn(err)
				}
			}
		}()
	})

	for range sig {
		c.Close()
		os.Exit(0)
	}
}

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
	"math"
	"time"
)

type SampleRequest struct {
	Timestamp time.Time
	Samples   []float64
}

func float64FromBytes(bytes []byte) float64 {
	bits := binary.BigEndian.Uint32(bytes)
	return float64(math.Float32frombits(bits))
}

func NewRequest(payload []byte) (*SampleRequest, error) {
	ts := time.Unix(int64(binary.LittleEndian.Uint64(payload[0:8])), 0)

	samplesPayload := payload[8:]
	samplesLen := len(samplesPayload) / 4
	samples := make([]float64, samplesLen)
	for i := 0; i < samplesLen; i++ {
		bits := samplesPayload[i*4 : i*4+4]
		samples[i] = float64FromBytes(bits)
	}

	return &SampleRequest{
		Timestamp: ts,
		Samples:   samples,
	}, nil
}

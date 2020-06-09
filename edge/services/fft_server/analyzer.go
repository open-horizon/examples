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
	"sort"
	"time"

	"github.com/mjibson/go-dsp/spectral"
	"github.com/sirupsen/logrus"
)

type Peaks map[float64]float64

type Analyzer struct {
	peaksLimit     int
	sampleRate     float64
	options        *spectral.PwelchOptions
	previous       Peaks
	previousTS     time.Time
	peaksThreshold float64
	freqsThreshold float64
}

func NewAnalyzer(config *Config) *Analyzer {
	return &Analyzer{
		peaksLimit:     config.Analyzer.PeaksLimit,
		peaksThreshold: config.Analyzer.PeaksThreshold,
		freqsThreshold: config.Analyzer.FreqsThreshold,
		sampleRate:     float64(config.SampleRate),
		options: &spectral.PwelchOptions{
			NFFT: config.Analyzer.NFFT,
		},
	}
}

func (a *Analyzer) getPeaks(samples []float64) Peaks {
	// TODO: change Pwelch to work with float32 instead float64
	peaks, frequencies := spectral.Pwelch(samples, a.sampleRate, a.options)
	peakFrequencies := make(Peaks)
	for i := 0; i < len(peaks); i++ {
		px := peaks[i] * 1000
		freq := frequencies[i]
		if px > 0 && freq > 0 {
			peakFrequencies[px] = freq
		}
	}
	return peakFrequencies
}

func (a *Analyzer) getTopPeaks(peaks Peaks) Peaks {
	// TODO: need better idea to sort by peak
	topPeaks := make(Peaks)
	sortedPeaks := make([]float64, 0, len(peaks))
	for peak := range peaks {
		sortedPeaks = append(sortedPeaks, peak)
	}
	sort.Float64s(sortedPeaks)
	for i := len(sortedPeaks) - 1; i > 0 && i > len(sortedPeaks)-int(a.peaksLimit)-1; i-- {
		peak := sortedPeaks[i]
		frequency := peaks[peak]
		topPeaks[peak] = frequency
		logrus.Debugf("Freq: %f\tPeak: %f\n", frequency, peak)
	}
	logrus.Debug("\n\n")

	return topPeaks
}

func (a *Analyzer) Analyze(req *SampleRequest) bool {
	peaks := a.getPeaks(req.Samples)
	topPeaks := a.getTopPeaks(peaks)

	// TODO: improve time complexity to at least O(n), implement ranges matrix
	for pPeak, pFreq := range a.previous {
		for peak, freq := range topPeaks {
			pPeakMax := pPeak + a.peaksThreshold
			pPeakMin := pPeak - a.peaksThreshold
			pFreqMax := pFreq + a.freqsThreshold
			pFreqMin := pFreq - a.freqsThreshold
			if pPeakMin <= peak && peak <= pPeakMax && pFreqMin <= freq && freq <= pFreqMax {
				delete(a.previous, pPeak)
			}
		}
	}

	a.previousTS = req.Timestamp
	if len(a.previous) > 0 {
		a.previous = topPeaks
		return true
	}
	a.previous = topPeaks
	return false
}

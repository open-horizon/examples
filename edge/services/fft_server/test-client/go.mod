module go-fft-test

go 1.23

require (
	github.com/jessevdk/go-flags v1.6.1
	github.com/sirupsen/logrus v1.9.3
	github.com/zenwerk/go-wave v0.0.0-20190102022600-1be84bfef50c
	go-fft-client v0.0.0
)

replace go-fft-client => ../fft_client


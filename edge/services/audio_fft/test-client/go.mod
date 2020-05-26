module go-fft-test

go 1.14

require (
	github.com/jessevdk/go-flags v1.4.0
	github.com/sirupsen/logrus v1.5.0
	github.com/zenwerk/go-wave v0.0.0-20190102022600-1be84bfef50c
	go-fft-client v0.0.0
)

replace go-fft-client => ../client

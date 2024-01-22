module go-fft-analyzer

go 1.14

require (
	github.com/jessevdk/go-flags v1.4.0
	github.com/mjibson/go-dsp v0.0.0-20180508042940-11479a337f12
	github.com/sirupsen/logrus v1.5.0
	go-fft-client v0.0.0
	golang.org/x/net v0.7.0 // indirect
)

replace go-fft-client => ../fft_client

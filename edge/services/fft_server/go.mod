module go-fft-analyzer

go 1.14

require (
	github.com/eclipse/paho.mqtt.golang v1.2.0 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/mjibson/go-dsp v0.0.0-20180508042940-11479a337f12
	github.com/sirupsen/logrus v1.5.0
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0 // indirect
	go-fft-client v0.0.0
)

replace go-fft-client => ../fft_client
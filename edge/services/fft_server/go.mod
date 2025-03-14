module go-fft-analyzer

go 1.23.0

toolchain go1.23.4

require (
	github.com/jessevdk/go-flags v1.6.1
	github.com/mjibson/go-dsp v0.0.0-20180508042940-11479a337f12
	github.com/sirupsen/logrus v1.9.3
	go-fft-client v0.0.0
)

require (
	github.com/eclipse/paho.mqtt.golang v1.5.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
)

replace go-fft-client => ../fft_client

module go-fft-analyzer

go 1.24.0

require (
	github.com/jessevdk/go-flags v1.6.1
	github.com/mjibson/go-dsp v0.0.0-20180508042940-11479a337f12
	github.com/sirupsen/logrus v1.9.3
	go-fft-client v0.0.0
)

require (
	github.com/eclipse/paho.mqtt.golang v1.5.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
)

replace go-fft-client => ../fft_client

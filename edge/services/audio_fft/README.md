# fft-example

This example is using Fast Fourier transform (FFT) [method](https://en.wikipedia.org/wiki/Fast_Fourier_transform) to compare input samples and send a "trigger" if thresholds are met.

Implementation consists of two parts: client and server.

Client is responsible for capturing a short audio sample using [PortAudio](http://www.portaudio.com) and sending it through MQTT topic to analyzer service.

Server is responsible for analyzing received audio sample and sending a trigger notification to a (possibly) different MQTT topic.  

Two horizon services used: one for the client, another -- for volantMQ and server. 

[VolantMQ](https://volantmq.io) is used as MQTT broker. Dockerfile and build script located in `/volanmq` folder.

## Configuration params

Both applications share the same MQTT config.

* `broker` -- MQTT broker location.
* `client` -- client ID to use. Depending on broker configuration, you'll need to use different IDs.
* `username` -- login to use.
* `password` -- password to use.
* `topic` -- topic name for the audio samples. 
* `result_topic` -- topic name for trigger alerts. If left empty, server will output to console only.
* `qos` -- quality of service to use with messages.

### Client params

* `log_level` -- [logrus](https://github.com/sirupsen/logrus) log level.
* `sample_rate` -- audio capturing frame rate. If you're not sure which rate to use, start with `debug` log level and check console output.
* `record_frame` -- number of seconds to record.

### Server params 

* `log_level` - [logrus](https://github.com/sirupsen/logrus) log level.
* `sample_rate` - frame rate used for received audio sample.
* `nfft` - number of data points for FFT.
* `peaks_limit` -- peaks limit.
* `peaks_threshold` -- peaks threshold.
* `freqs_threshold` -- frequencies threshold.

Please note: `sample_rate` must be the same for correct processing. 

## Installing portaudio on Mac

1. Install `pkg-config`

```bash
brew install pkg-config 
```

Make sure `/usr/local/bin` is in the `$PATH` otherwise override env variable `PKG_CONFIG=/usr/local/bin/pkg-config`.

2. Install `PortAudio`

```bash
brew install portaudio
```

Note: It's not possible to run `client` container on a mac, since there's no `/dev/snd` analogue. Binary itself works.

## Publishing services

`/horizon` folder contains readme files for preparing MQTT and publishing services, as well as all pre-requirements for configuring and registering RPi node. 

### Test client

`/test-client` folder contains a simple test client which will sand 10 random samples from `/test-client/sets` and compare expected results. 

Please note: if you just started the server, first sample will have an incorrect response, since the server is comparing current sample with previous; and when server starts there's no initial state.

To run the test-client first start dev server

```bash
cd server
hzn dev service start -S -v             
```   

Place you `.wav` files under `/sets/<set-name>` folder and start the test:

```bash
./test-client -b localhost:1883 -u fft-client -p client-pass -c fft-test --result_topic results
```

Two samples provided: 
1. Motors -- custom recording of a rotor motors during a normal operation and "malfunction".
2. Killer [Whale](https://commons.wikimedia.org/w/index.php?title=File%3AKiller_whale_residents_broadband.ogg#) -- recording made by National Park Service, using a hydrophone that is anchored near the mouth of Glacier Bay, Alaska for the purpose of monitoring ambient noise. Available under [Creative Commons CC0 License](https://creativecommons.org/publicdomain/zero/1.0/).

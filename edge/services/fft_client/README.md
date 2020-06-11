# fft-example

This example is using Fast Fourier transform (FFT) [method](https://en.wikipedia.org/wiki/Fast_Fourier_transform) to compare input samples and send a "trigger" if thresholds are met.

Client is responsible for capturing a short audio sample using [PortAudio](http://www.portaudio.com) and sending it through MQTT topic to analyzer service.

## Building

You will need a quemu to build a multi-arch image. For Mac just update to edge version.

```
gmake build
```

## Configuration params

* `broker` -- MQTT broker location.
* `client` -- client ID to use. Depending on broker configuration, you'll need to use different IDs.
* `username` -- login to use.
* `password` -- password to use.
* `topic` -- topic name for the audio samples. 
* `qos` -- quality of service to use with messages.
* `log_level` -- [logrus](https://github.com/sirupsen/logrus) log level.
* `sample_rate` -- audio capturing frame rate. If you're not sure which rate to use, start with `debug` log level and check console output.
* `record_frame` -- number of seconds to record.

Please note: `sample_rate` must be the same on server and client for correct processing. 

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
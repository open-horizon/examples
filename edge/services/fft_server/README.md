# fft-example

This example is using Fast Fourier transform (FFT) [method](https://en.wikipedia.org/wiki/Fast_Fourier_transform) to compare input samples and send a "trigger" if thresholds are met.

Server is responsible for analyzing received audio sample and sending a trigger notification to a (possibly) different MQTT topic. 

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
* `result_topic` -- topic name for trigger alerts. If left empty, server will output to console only.
* `qos` -- quality of service to use with messages.
* `log_level` - [logrus](https://github.com/sirupsen/logrus) log level.
* `sample_rate` - frame rate used for received audio sample.
* `nfft` - number of data points for FFT.
* `peaks_limit` -- peaks limit.
* `peaks_threshold` -- peaks threshold.
* `freqs_threshold` -- frequencies threshold.

Please note: `sample_rate` must be the same on server and client for correct processing. 

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
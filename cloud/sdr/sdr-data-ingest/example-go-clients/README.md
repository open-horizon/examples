# IBM Message Hub Producer and Consumer Client Examples in Go

## Setup

```
go get github.com/Shopify/sarama
go get github.com/bsm/sarama-cluster
export MSGHUB_API_KEY='abcdefg'
export MSGHUB_BROKER_URL="kafka01-prod123.messagehub.services.us-south.bluemix.net:9093,kafka02-prod456.messagehub.services.us-south.bluemix.net:9093"
export MSGHUB_TOPIC="sdr-audio"
export STT_USERNAME="123456789-abcd"
export STT_PASSWORD="ABCDEFGHIJKLMNOP"
```

You can optionally use your own certificate, which can be created like this:
```
openssl genrsa -out server.key 2048
openssl req -new -x509 -key server.key -out server.pem -days 3650
```

## Build All Examples

```
make
```

## Produce Messages to IBM Message Hub

```
msghub-producer 'hello world'
msghub-producer -t <topic> 'hello world'   # produce to a different topic
msghub-producer -v 'hello world'     # see verbose output
msghub-producer    # will publish several generated msgs
msghub-producer -h     # see all of the flags and environment variables
```

## Produce Messages to IBM Message Hub

```
msghub-consumer
msghub-consumer -t <topic>   # consume from a different topic
msghub-consumer -v     # see verbose output
msghub-consumer -h     # see all of the flags and environment variables
```

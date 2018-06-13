# IBM Message Hub Producer and Consumer Client Examples in Go

## Setup

```
go get github.com/Shopify/sarama
go get github.com/bsm/sarama-cluster
openssl genrsa -out server.key 2048
openssl req -new -x509 -key server.key -out server.pem -days 3650
export MSGHUB_API_KEY='abcdefg'
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

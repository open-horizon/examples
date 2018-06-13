# IBM Message Hub Publish and Consume Client Examples in Go

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

## Publish Synchronously

```
msghub-pubsync 'hello world'
msghub-pubsync -v 'hello world'     # see verbose output
msghub-pubsync    # will publish several generated msgs
msghub-pubsync -h     # see all of the flags and environment variables
```

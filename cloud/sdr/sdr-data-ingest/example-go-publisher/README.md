# IBM Message Hub Publish and Consume Client Examples in Go

## Setup

```
go get github.com/Shopify/sarama
openssl genrsa -out server.key 2048
openssl req -new -x509 -key server.key -out server.pem -days 3650
export MSGHUB_API_KEY='abcdefg'
```

## Publish Synchronously

```
go run msghubpub.go
```

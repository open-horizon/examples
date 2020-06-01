FROM golang:alpine as build-env
ADD . /src
WORKDIR /src/test-client
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /src/fft-test

FROM alpine
ENV MQTT_BROKER=host.docker.internal:1883 \
	MQTT_SERVER_USER=fft-client \
	MQTT_SERVER_PASS=client-pass \
	MQTT_SERVER_CLIENT=fft-test \
	MQTT_RESULTS_TOPIC=results

COPY --from=build-env /src/fft-test /fft-test
COPY test-client/start.sh .
COPY test-client/sets/ ./sets/
RUN chmod +x start.sh
ENTRYPOINT ["./start.sh"]
ARG BUILD_IMAGE
ARG RUN_IMAGE
FROM $BUILD_IMAGE AS build

ADD . /src
WORKDIR /src/server

ENV GOOS=linux

ARG GOARCH
ENV GOARCH=$GOARCH

RUN apk --no-cache add build-base git bzr mercurial gcc && \
	go get -d -v && \
	go build -o /server

##################################################################################################

FROM $RUN_IMAGE

ENV MQTT_BROKER=volantmq:1883 \
	MQTT_SERVER_USER=fft-server \
	MQTT_SERVER_PASS=server-pass \
	MQTT_SERVER_CLIENT=fft-server \
	MQTT_RESULTS_TOPIC=results \
	SAMPLE_RATE=48000

COPY --from=build /server /server
ADD ./server/start.sh .

RUN chmod +x start.sh

ENTRYPOINT ["./start.sh"]
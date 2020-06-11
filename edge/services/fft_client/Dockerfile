ARG BUILD_IMAGE
ARG RUN_IMAGE
FROM $BUILD_IMAGE AS build

ADD . /src
WORKDIR /src

ENV GOOS=linux

ARG GOARCH
ENV GOARCH=$GOARCH

RUN apk --no-cache add build-base git bzr mercurial gcc portaudio portaudio-dev pkgconfig && \
	go get -d -v && \
	go build -o /client


##################################################################################################

FROM $RUN_IMAGE

ENV MQTT_BROKER=volantmq:1883 \
	MQTT_CLIENT_USER=fft-client \
	MQTT_CLIENT_PASS=client-pass \
	MQTT_CLIENT_CLIENT=fft-client \ 
	MQTT_QOS=2 \
	LOG_LEVEL=info \
	SAMPLE_RATE=48000 \ 
	RECORD_FRAME=3

COPY --from=build /client /client
ADD start.sh .

RUN apk --no-cache add portaudio && \
	chmod +x start.sh

ENTRYPOINT ["./start.sh"]

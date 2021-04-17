FROM alpine:latest

RUN apk update \
    && apk add --no-cache bash mosquitto mosquitto-clients \
    && rm -fr /tmp/*

COPY mqtt.sh /
WORKDIR /
RUN ["chmod", "+x", "./mqtt.sh"]
CMD ["./mqtt.sh"]

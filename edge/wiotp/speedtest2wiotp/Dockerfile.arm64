FROM arm64v8/python:3-alpine

RUN apk --no-cache add curl jq mosquitto-clients
COPY *.sh /
WORKDIR /
CMD sh /speedtest2wiotp.sh


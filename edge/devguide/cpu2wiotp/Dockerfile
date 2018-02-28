FROM alpine:latest
RUN apk --no-cache --update add curl mosquitto-clients bc jq
COPY *.sh /
WORKDIR /
CMD /workload.sh

FROM alpine:latest
RUN apk --no-cache --update add curl ca-certificates wget jq bc

# Install kafka to publish msgs to IBM Message Hub
RUN wget --quiet --output-document=/etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub && wget https://github.com/sgerrand/alpine-pkg-kafkacat/releases/download/1.3.1-r0/kafkacat-1.3.1-r0.apk && apk --no-cache add kafkacat-1.3.1-r0.apk && rm kafkacat-1.3.1-r0.apk

COPY *.sh /
WORKDIR /
CMD /service.sh

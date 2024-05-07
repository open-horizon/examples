FROM alpine:latest

RUN apk --no-cache --update add jq bash

COPY *.sh /
WORKDIR /
CMD /service.sh
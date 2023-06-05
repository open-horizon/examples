FROM alpine:3.16

RUN apk --no-cache --update add jq 'curl<7.88.0' bash

COPY *.sh /
WORKDIR /
CMD /service.sh

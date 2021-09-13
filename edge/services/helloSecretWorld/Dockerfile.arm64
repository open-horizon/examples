FROM arm64v8/alpine:latest

RUN apk --no-cache --update add jq curl bash

COPY *.sh /
WORKDIR /
CMD /service.sh
FROM alpine:latest

RUN apk --no-cache --update add curl

COPY start.sh /
COPY service /usr/local/bin/

RUN mkdir -p /ess-store

WORKDIR /
CMD /start.sh
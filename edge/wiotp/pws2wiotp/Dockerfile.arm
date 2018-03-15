FROM alpine:latest
MAINTAINER dyec@us.ibm.com

ENV ARCH=arm

RUN apk --no-cache add curl mosquitto-clients
COPY *.sh /
COPY *.pem /

## For development
#RUN mkdir -p devenv
#COPY devenv/*.pem /devenv/

WORKDIR /
CMD /workload.sh

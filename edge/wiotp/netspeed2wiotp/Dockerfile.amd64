FROM alpine:latest
MAINTAINER Chris Dye <dyec@us.ibm.com>

ENV ARCH=x86

# Need to do this on a single line so this docker image layer will have the pkgs removed
RUN apk --no-cache --update add python python-dev py-pip && pip install --upgrade pip paho-mqtt && apk del python-dev py-pip

COPY *.py /
COPY *.pem /
WORKDIR /

CMD python netspeed_edge.py --verbose --mqtt --policy

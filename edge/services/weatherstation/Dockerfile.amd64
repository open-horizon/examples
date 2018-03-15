# Set base image to alpine linux (very lightweight)
FROM alpine:latest

# File Author / Maintainer 
MAINTAINER dyec@us.ibm.com

ENV ARCH=amd64
RUN apk --no-cache add python py-pip && apk --no-cache add py-configobj libusb-compat libusb-compat-dev libusb-dev && rm -rf /var/cache/apk/* && pip install six flask multiprocess

# Copy weewx files over
RUN mkdir -p /tmp/weather
COPY weewx/*.tar.gz /tmp/weather/

# Unpack / build some dependencies from source
WORKDIR /tmp/weather
RUN mkdir pu
RUN tar xvf pyusb* -C pu --strip-components=1
WORKDIR /tmp/weather/pu
RUN python setup.py install

## Unpack / build-install weather software
WORKDIR /tmp/weather
RUN mkdir wwx
RUN tar xvf weewx* -C wwx --strip-components=1
WORKDIR /tmp/weather/wwx
RUN ./setup.py build
COPY weewx/answers.txt /tmp/weather/wwx/
# An initial set of Station params: Loc/Altitude/Lat/Lon (SF, CA USA)
RUN cat ./answers.txt | ./setup.py install
# The above step installs weewx to /home/weewx/bin/

# Copy microservice scripts: start.py, weewx_mod.py, flask_server.py
COPY weewx/*.py /home/weewx/bin/
# Copy edited weewx engine file (we added a shared dict)
COPY weewx/bin-weewx/engine.py /home/weewx/bin/weewx/
WORKDIR /home/weewx/bin

# Install stuff for container dev/inspection (comment for production)
#RUN apk --no-cache add vim dropbear-scp curl && rm -rf /var/cache/apk/*

# Remove temp files, cleanup
RUN rm -rf /tmp/weather && apk --no-cache del py-pip

# Update config file, Run weewx and flask server processes
CMD python start.py ../weewx.conf

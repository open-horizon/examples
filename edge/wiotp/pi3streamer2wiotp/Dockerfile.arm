FROM armhf/alpine:latest
ENV ARCH=armhf
RUN apk --no-cache --update add curl mosquitto-clients jq bc
COPY *.sh /
COPY *.pem /


WORKDIR /
CMD /workload.sh

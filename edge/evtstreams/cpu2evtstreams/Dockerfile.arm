FROM arm32v6/alpine:latest
RUN apk --no-cache --update add curl jq bc

# We build kafkacat outside of this Dockerfile using: https://github.com/sgerrand/alpine-pkg-kafkacat and the instructions in README.md
# Install kafkacat (the Makefile already copied these files from ../../../tools/kafkacat/$ARCH) to tmp/$ARCH
#todo: i do not know how to do these 2 COPYs and the RUN in the same layer, so removing the apk pkg is effective
# Note: https://docs.docker.com/engine/reference/builder/#copy says "you cannot COPY ../something /something, because the first step of a docker build is to send the context directory (and subdirectories) to the docker daemon."
COPY tmp/arm/*.rsa.pub /etc/apk/keys/
COPY tmp/arm/kafkacat-*.apk /
RUN apk --no-cache add /kafkacat-*.apk && rm kafkacat-*.apk

COPY *.sh /
WORKDIR /
CMD /service.sh

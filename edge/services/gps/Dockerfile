# Horizon GPS Location Service: gps
#
# This server provides REST access to gps receiver location and satellite data
# (or gps location estimated from the IP address if hardware is not available).
#
# More precise documentation of the behavior of this container may be found
# in the src/main.go source code.
#
# To build this server container, run the following command in this directory:
#   $ make
#
# Written by Glen Darling, November 2016.
# Updated to 2-stage build, and modified to target arm64v8, May 2019.

# Build stage 0: Go compilation
FROM golang:1.15-alpine as go_build

RUN apk --no-cache update && apk add git

RUN mkdir -p /build/bin
COPY src /build/src

ARG ARCH=${ARCH}

WORKDIR /build
RUN env GOPATH=/build GOOPTIONS_${ARCH}="CGO_ENABLED=0 GOOS=linux GOARCH=$ARCH" go get github.com/kellydunn/golang-geo
RUN env GOPATH=/build GOOPTIONS_${ARCH}="CGO_ENABLED=0 GOOS=linux GOARCH=$ARCH" go build -o /build/bin/${ARCH}_gps /build/src/main.go

FROM alpine:latest

# Install the gpsd daemon, and the certs needed to use https services
RUN apk update && apk add gpsd curl --no-cache ca-certificates

# Copy in the server binary from stage 0 of the build (above)
ARG ARCH=${ARCH}
COPY --from=0 /build/bin/${ARCH}_gps /gps

# Create hzngroup and hznuser
RUN addgroup -S hzngroup && adduser -S hznuser -G hzngroup

# Run container as hznuser user
USER hznuser

# The gps service uses this port to respond to REST requests
EXPOSE 8080

# Set the default command to be the go executable to start everything
CMD /gps
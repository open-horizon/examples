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

# Build stage 0: Go compilation
FROM golang:1.19-alpine as go_build
LABEL stage=builder

RUN apk --no-cache update && apk add git

RUN mkdir -p /build/bin
COPY src /build/src

ARG ARCH=${ARCH}

WORKDIR /build                                                                    
RUN env GO111MODULE=off env GOPATH=/build GOOPTIONS_${ARCH}="CGO_ENABLED=0 GOOS=linux GOARCH=$ARCH" go get github.com/kellydunn/golang-geo
RUN env GO111MODULE=off env GOPATH=/build GOOPTIONS_${ARCH}="CGO_ENABLED=0 GOOS=linux GOARCH=$ARCH" go build -o /build/bin/gps /build/src/main.go

FROM alpine:latest

# Install the gpsd daemon, and the certs needed to use https services
RUN apk update && apk add gpsd curl --no-cache ca-certificates

# Copy in the server binary of the build (above)
ARG ARCH=${ARCH}
COPY --from=go_build /build/bin/gps /gps

# Create hzngroup and hznuser
RUN addgroup -S hzngroup && adduser -S hznuser -G hzngroup

# Run container as hznuser user
USER hznuser

# The gps service uses this port to respond to REST requests
EXPOSE 8080

# Set the default command to be the go executable to start everything
CMD /gps
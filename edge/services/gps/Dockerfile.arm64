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

FROM arm64v8/golang:alpine

RUN apk --no-cache update && apk add git
RUN apk add gcc musl musl-dev

RUN mkdir -p /build/bin
COPY src /build/src

WORKDIR /build
RUN env GOPATH=/build GOOPTIONS_ARM64='CGO_ENABLED=0 GOOS=linux GOARCH=arm64' go get github.com/kellydunn/golang-geo
RUN env GOPATH=/build GOOPTIONS_ARM64='CGO_ENABLED=0 GOOS=linux GOARCH=arm64 ' go build -o /build/bin/arm64_gps /build/src/main.go



# Build stage 1: The final container (including armv6_gps binary from above)

FROM aarch64/alpine:latest

# Install the gpsd daemon, and the certs needed to use https services
RUN apk update && apk add gpsd --no-cache ca-certificates

# Copy in the server binary from stage 0 of the build (above)
COPY --from=0 /build/bin/arm64_gps /gps

# The gps service uses this port to respond to REST requests
EXPOSE 80

# Set the default command to be the go executable to start everything
CMD /gps

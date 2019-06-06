#
# Blue Horizon Firmware Device API: gps
#
# This server provides REST access to gps receiver location and satellite data
#
# More precise documentation of the behavior of this container may be found
# in the src/main.go source code.
#
# To build this server container, run the following command in this directory:
#   $ make
#
# To run the firmware container as a daemon process (e.g., so you can test it):
#   $ make daemon
#
# To run the firmware container in dev mode (normally used for development):
#   $ make develop
#
# Written by Glen Darling, November 2016.
#


# Build stage 0: Go compilation

FROM golang:1.10.0-alpine as go_build

RUN apk --no-cache update && apk add git

RUN mkdir -p /build/bin
COPY src /build/src

WORKDIR /build
RUN env GOPATH=/build GOOPTIONS_AMD64='CGO_ENABLED=0 GOOS=linux GOARCH=amd64' go get github.com/kellydunn/golang-geo
RUN env GOPATH=/build GOOPTIONS_AMD64='CGO_ENABLED=0 GOOS=linux GOARCH=amd64' go build -o /build/bin/amd64_gps /build/src/main.go



# Build stage 1: The final container (including armv6_gps binary from above)

FROM alpine:latest

# Install the gpsd daemon, and the certs needed to use https services
RUN apk update && apk add gpsd --no-cache ca-certificates

# Copy in the server binary from stage 0 of the build (above)
COPY --from=0 /build/bin/amd64_gps /gps

# The gps service uses this port to respond to REST requests
EXPOSE 80

# Set the default command to be the go executable to start everything
CMD /gps

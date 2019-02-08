#
# Blue Horizon Firmware (Device API) test container (amd64)
#
# Run this on an amd64 host to test this firmware container
#
# To build the firmware container for testing:
#   $ make
#
# To run the firmware container as a daemon process (so you can test it):
#   $ make daemon
#
# To build the test container:
#   $ make build_test
#
# To run the test container:
#   $ make test
#
# Written by Glen Darling, November 2016.
#

FROM alpine:latest
MAINTAINER glendarling@us.ibm.com

# Install tools we may need during debug
#RUN apk --no-cache add musl-dev bash vim curl jq

COPY bin/gps-test.amd64 /gps-test

CMD /gps-test

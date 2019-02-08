#
# Blue Horizon Firmware (Device API) test container (arm)
#
# Run this on an arm host to test this firmware container
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

FROM arm32v6/alpine
MAINTAINER glendarling@us.ibm.com

# Install mosquitto-clients and other useful tools we may want for test/debug
RUN apk --no-cache add mosquitto-clients musl-dev bash vim curl wget jq

# Copy over the test code
COPY bin/gps-test.arm /gps-test

# Default command will run the tests
CMD /gps-test

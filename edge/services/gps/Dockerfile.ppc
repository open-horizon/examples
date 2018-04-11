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

FROM ppc64le/debian:latest
MAINTAINER bp@us.ibm.com

# Install gpsd
RUN apt-get update && apt-get -y install gpsd

# Instruct docker to make this port available to clients like the location container
EXPOSE 31779

# Copy over the ppc server binary (done 2nd to last because it changes the most)
COPY bin/gps.ppc /gps

# Set the default command to be the go executable to start everything
CMD /gps

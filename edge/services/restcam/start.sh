#!/bin/sh

echo "Starting ibm.restcam..."
socat TCP4-LISTEN:80,fork EXEC:./cam.sh


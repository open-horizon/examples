#!/bin/sh

socat TCP4-LISTEN:80,fork EXEC:./cam.sh


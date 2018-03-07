#!/bin/bash

# Setting defaults
VF=""
HF=""
RES="-x 640 -y 480"
FPS="20"

# Check settings (horiz/vertical flip, resolution)
if [[ "$VERT_FLIP" == "1" ]]; then
  VF="-vf"
  echo "pi3streamer: start.sh: VERT_FLIP=1, setting flag: '${VF}'"
fi

if [[ "$HORZ_FLIP" == "1" ]]; then
  HF="-hf"
  echo "pi3streamer: start.sh: HORZ_FLIP=1, setting flag: '${HF}'"
fi

# Dangerous - no value check
if [[ ! -z "$RESOLUTION" ]]; then
  RES="$RESOLUTION"
  echo "pi3streamer: start.sh: RESOLUTION='${RESOLUTION}' (format: '-x <width> -y <height>')"
fi

if [[ ! -z "$FRAMERATE" ]]; then
  FPS="$FRAMERATE"
  echo "pi3streamer: start.sh: FRAMERATE=${FPS} fps"
fi

## Run the picam streamer microservice
./mjpg_streamer -o "output_http.so -w ./www" -i "input_raspicam.so $RES -fps $FPS -ex night $VF $HF"

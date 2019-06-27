#!/bin/bash

MESSAGE="Hello!"
while true; do
  # echo "Checking ESS..."
  NEW=`/ess-read-all.sh example-type 2>/dev/null | jq '.[].data' | head -n 1`
  if [ ! -z "$NEW" ]; then
    MESSAGE="${NEW}"
  fi
  echo "${HZN_DEVICE_ID} says: \"${MESSAGE}\""
  sleep 3
done

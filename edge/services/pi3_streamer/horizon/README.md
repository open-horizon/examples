This directory contains an input file template to create a Horizon Exchange microservice definition for the pi3streamer docker image.

Fill in the values of the variables in the template with commands like:

```
export DOCKER_HUB_ID=openhorizon   # or your own docker hub id
export ARCH2=arm     # Not available yet for amd64 or arm64 (requires RPi Camera)
export PI3STREAMER_VERSION=1.0.0
export HZN_ORG_ID=abcdef
export WIOTP_DOMAIN=internetofthings.ibmcloud.com

envsubst < pi3streamer-template.json > pi3streamer.json
```

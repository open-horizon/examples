This directory contains input file templates to create a Horizon Exchange pattern definitions and input files for the pi3streamer2wiotp pattern.

Fill in the values of the variables in the template with commands like:

```
export DOCKER_HUB_ID=openhorizon   # or your own docker hub id
export ARCH2=arm    # arm only
export PI3STREAMER2WIOTP_VERSION=1.0
export HZN_ORG_ID=abcdef
export WIOTP_DOMAIN=internetofthings.ibmcloud.com

envsubst < pi3streamer2wiotp-template.json > pi3streamer2wiotp.json
envsubst < insert-pi3streamer2wiotp-template.json > insert-pi3streamer2wiotp.json
```

This directory contains an input file template to create a Horizon Exchange workload definition for the cpu2wiotp docker image.

Fill in the values of the variables in the template with commands like:

```
export DOCKER_HUB_ID=openhorizon   # or your own docker hub id
export ARCH2=amd64    # or arm or arm64
export CPU2WIOTP_VERSION=1.2.1
export HZN_ORG_ID=abcdef

envsubst < cpu2wiotp-template.json > cpu2wiotp.json
```

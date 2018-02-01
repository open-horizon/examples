This directory contains an input file template to create a Horizon Exchange microservice definition for the microphone docker image.

Fill in the values of the variables in the template with commands like:

```
export DOCKER_HUB_ID=openhorizon   # or your own docker hub id
export ARCH1=x86     # or arm or arm64
export ARCH2=amd64    # or arm or arm64
export MIC_VERSION=1.2.4

envsubst < microphone-template.json > microphone.json
```

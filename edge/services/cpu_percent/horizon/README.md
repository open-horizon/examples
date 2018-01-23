This directory contains an input file template to create a Horizon Exchange microservice definition for the cpu docker image.

Fill in the values of the variables in the template with commands like:

```
export DOCKER_HUB_ID=openhorizon   # or your own docker hub id
export ARCH1=x86     # or arm or arm64
export ARCH2=amd64    # or arm or arm64
export CPU_VERSION=1.2.2

envsubst < cpu-template.json > cpu.json
```

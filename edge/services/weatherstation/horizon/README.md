This directory contains an input file template to create a Horizon Exchange microservice definition for the pwsms docker image.

Fill in the values of the variables in the template with commands like:

```
export DOCKER_HUB_ID=openhorizon   # or your own docker hub id
export ARCH2=amd64     # or arm or arm64
export VERSION=1.1.0

envsubst < pwsms-template.json > pwsms.json
```

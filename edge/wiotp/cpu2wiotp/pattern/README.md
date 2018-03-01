This directory contains input file templates to create a Horizon Exchange pattern definitions and input files for the cpu2wiotp pattern.

Fill in the values of the variables in the template with commands like:

```
export DOCKER_HUB_ID=openhorizon   # or your own docker hub id
export ARCH2=amd64    # or arm or arm64
export CPU2WIOTP_VERSION=1.2.1
export HZN_ORG_ID=abcdef
export WIOTP_DOMAIN=internetofthings.ibmcloud.com

envsubst < cpu2wiotp-template.json > cpu2wiotp.json
envsubst < insert-cpu2wiotp-template.json > insert-cpu2wiotp.json
```

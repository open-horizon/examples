This directory contains input file templates to create a Horizon Exchange pattern definitions and input files for the netspeed2wiotp pattern.

Fill in the values of the variables in the template with commands like:

```
export DOCKER_HUB_ID=openhorizon   # or your own docker hub id
export ARCH2=amd64    # or arm or arm64
export NETSPEED2WIOTP_VERSION=2.6
export HZN_ORG_ID=abcdef
export WIOTP_DOMAIN=internetofthings.ibmcloud.com

envsubst < netspeed2wiotp-template.json > netspeed2wiotp.json
envsubst < insert-netspeed2wiotp-template.json > insert-netspeed2wiotp.json
```

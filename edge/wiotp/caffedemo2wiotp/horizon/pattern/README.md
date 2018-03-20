This directory contains input file templates to create a Horizon Exchange pattern definitions and input files for the service pattern.

Fill in the values of the variables in the template with commands like:

```
export DOCKER_HUB_ID=openhorizon   # or your own docker hub id
export ARCH2=amd64    # or arm or arm64
export WL_VERSION=1.0
export WIOTP_ORG_ID=abcdef
export WIOTP_TEST_ENV2=''
export WIOTP_EDGE_MQTT_IP=10.1.2.3   # the private IP of your edge node
export WIOTP_CLASS_ID=g
export WIOTP_GW_TYPE=mygwtype
export WIOTP_GW_ID=mygw
export WIOTP_GW_TOKEN=mytok

envsubst < caffedemo-template.json > caffedemo.json
envsubst < caffedemo-input-template.json > caffedemo-input.json
envsubst < insert-caffedemo-template.json > insert-caffedemo.json
envsubst < hznEdgeCoreIoTInput.json.template > hznEdgeCoreIoTInput.json
```

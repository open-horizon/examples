This directory contains input file templates to create a Horizon Exchange pattern definitions and input files for the netspeed2wiotp pattern.

Fill in the values of the variables in the template with commands like:

```
export DOCKER_HUB_ID=openhorizon   # or your own docker hub id
export ARCH2=amd64    # or arm or arm64
export NETSPEED2WIOTP_VERSION=1.1.8
export WIOTP_ORG_ID=abcdef
export WIOTP_TEST_ENV2=''
export WIOTP_EDGE_MQTT_IP=10.1.2.3   # the private IP of your edge node
export WIOTP_CLASS_ID=g
export WIOTP_GW_TYPE=mygwtype
export WIOTP_GW_ID=mygw
export WIOTP_GW_TOKEN=mytok

envsubst < netspeed2wiotp-template.json > netspeed2wiotp.json
envsubst < netspeed2wiotp-input-template.json > netspeed2wiotp-input.json
envsubst < insert-netspeed2wiotp-template.json > insert-netspeed2wiotp.json
envsubst < hznEdgeCoreIoTInput.json.template > hznEdgeCoreIoTInput.json
```

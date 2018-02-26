This directory contains pre-signed input file templates to create a Horizon Exchange workload definitions for the cpu2wiotp docker image. These are used by the quick start guide to speed the process.

Fill in the values of the variables in the template with commands like:

```
export HZN_ORG_ID=abcdef

envsubst < cpu2wiotp-template-amd64.json > cpu2wiotp-amd64.json
envsubst < cpu2wiotp-template-arm.json > cpu2wiotp-arm.json
envsubst < cpu2wiotp-template-arm64.json > cpu2wiotp-arm64.json
```

# Horizon CPU To IBM Message Hub Service

Repeatedly queries the edge node's CPU percentage from the CPU Percent service and the GPS location from the GPS service, and then sends both to the IBM Message Hub. The topic is sends to is: `$HZN_ORG_ID.$HZN_DEVICE_ID`

## Input Values

The following input values can be given to this service in the input file given to `hzn register`:


| Name | Required? | Type | Description |
| ---- | --------- | ---- | ---------------- |
| MSGHUB_API_KEY | yes | string | the API key of the IBM Message Hub instance you are sending data to |
| MOCK | no | boolean | default is false. If true, send fake data instead of querying the cpu and gps services |
| PUBLISH | no | boolean | default is true. If false, do not send data to message hub, only print it to the log |
| SAMPLE_INTERVAL | no | integer | default is 5. How often (in seconds) to query the cpu percent. (The gps location is queried every SAMPLE_INTERVAL * SAMPLE_SIZE seconds.)  |
| VERBOSE | no | integer | default is 10. The number of cpu samples to read before calculating and publishing the cpu average and gps coordinates |


#### Example:
A sample `services` section of the input file given to `hzn register`:
```
    "services": [
        {
            "org": "$HZN_ORG_ID",
            "url": "https://$MYDOMAIN/service-$CPU2MSGHUB_NAME",
            "versionRange": "[0.0.0,INFINITY)",
            "variables": {
                "MSGHUB_API_KEY": "$MSGHUB_API_KEY",
                "MOCK": false,
                "PUBLISH": true,
                "SAMPLE_INTERVAL": 2,
                "SAMPLE_SIZE": 5,
                "VERBOSE": "1"
            }
        }
    ]
```

## Deployment Values

These values can be specified in the `deployment` section of the service definition, or be overridden in the `deployment_overrides` section of your pattern.


| Name | Type | Description |
| ---- | ---- | ---------------- |
| MSGHUB_BROKER_URL | comma-separated string | One or more kafka bootstrap hostnames for a client to communicate to the instance of IBM Message Hub |


#### Example:
```
    "deployment": {
        "services": {
            "cpu2msghub": {
                "environment": [
                    "MSGHUB_BROKER_URL=$MSGHUB_BROKER_URL"
                ],
                "image": "$DOCKER_HUB_ID/${ARCH}_$CPU2MSGHUB_NAME:$CPU2MSGHUB_VERSION"
            }
        }
    },
```

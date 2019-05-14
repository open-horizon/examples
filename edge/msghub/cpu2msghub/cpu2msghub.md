# Horizon CPU To IBM Message Hub Service

Repeatedly queries the edge node's CPU percentage from the CPU Percent service and the GPS location from the GPS service, and then sends both to the IBM Message Hub. The topic is sends to is: `cpu2msghub` (can be overridden)

## Input Values

The following input values **must** be given to this service in the input file given to `hzn register`:


| Name | Required? | Type | Description |
| ---- | --------- | ---- | ---------------- |
| MSGHUB_API_KEY | yes | string | the API key of the IBM Message Hub instance you are sending data to |
| MSGHUB_BROKER_URL | yes | string | The comma-separated list of URLs to use when sending messages to your instance of IBM Message Hub |


These **optional** input values can be overridden:


| Name | Required? | Type | Description |
| ---- | --------- | ---- | ---------------- |
| MSGHUB_TOPIC | yes | string | The topic to use when sending messages to your instance of IBM Message Hub |
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
                "MSGHUB_BROKER_URL": "$MSGHUB_BROKER_URL",
                "VERBOSE": "1"
            }
        }
    ]
```

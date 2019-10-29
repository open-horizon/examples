# Horizon CPU To IBM Event Streams Service

Repeatedly queries the edge node's CPU percentage from the CPU Percent service and the GPS location from the GPS service, and then sends both to the IBM Event Streams. The topic is sends to is: `cpu2evtstreams` (can be overridden)

## Input Values

The following input values **must** be given to this service in the input file given to `hzn register`:


| Name | Required? | Type | Description |
| ---- | --------- | ---- | ---------------- |
| EVTSTREAMS_API_KEY | yes | string | the API key of the IBM Event Streams instance you are sending data to |
| EVTSTREAMS_BROKER_URL | yes | string | The comma-separated list of URLs to use when sending messages to your instance of IBM Event Streams |


These **optional** input values can be overridden:


| Name | Required? | Type | Description |
| ---- | --------- | ---- | ---------------- |
| EVTSTREAMS_TOPIC | yes | string | Default is `cpu2evtstreams`. The topic to use when sending messages to your instance of IBM Event Streams |
| EVTSTREAMS_CERT_ENCODED | no | string | Default is `-`. The base64-encoded self-signed certificate to use when sending messages to your ICP instance of IBM Event Streams. Not needed for IBM Cloud Event Streams. |
| MOCK | no | boolean | default is false. If true, send fake data instead of querying the cpu and gps services |
| PUBLISH | no | boolean | default is true. If false, do not send data to Event Streams, only print it to the log |
| SAMPLE_INTERVAL | no | integer | default is 5. How often (in seconds) to query the cpu percent. (The gps location is queried every SAMPLE_INTERVAL * SAMPLE_SIZE seconds.)  |
| VERBOSE | no | integer | default is 10. The number of cpu samples to read before calculating and publishing the cpu average and gps coordinates |


#### Example:
A sample `services` section of the input file given to `hzn register`:
```
    "services": [
        {
            "org": "$HZN_ORG_ID",
            "url": "$SERVICE_NAME",
            "variables": {
                "EVTSTREAMS_API_KEY": "$EVTSTREAMS_API_KEY",
                "EVTSTREAMS_BROKER_URL": "$EVTSTREAMS_BROKER_URL",
                "VERBOSE": "1"
            }
        }
    ]
```

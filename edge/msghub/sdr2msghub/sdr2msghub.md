# Horizon SDR To IBM Message Hub Service

Sample Horizon service that sends 30 second clips of FM radio rich in speech to IBM Message Hub. It contains Edge Node software that requires SDR hardware (but it can simulate that hardware when it is not present, which is especially useful during development). The Edge Node software receives radio signals, does some local analysis, and sends lower-volume, higher-value data to the cloud. The SDR example also contains a powerful cloud back end implementation for the application. The back end receives data from the Edge Nodes, presents a web UI with a map upon which your Edge Nodes appear. It also performs deeper data analysis by using IBM Watson APIs.The topic is sends to is: `sdr-audio` (can be overridden)

## Input Values

The following input values **must** be given to this service in the input file given to `hzn register`:


| Name | Required? | Type | Description |
| ---- | --------- | ---- | ---------------- |
| MSGHUB_API_KEY | yes | string | the API key of the IBM Message Hub instance you are sending data to |
| MSGHUB_BROKER_URL | yes | string | The comma-separated list of URLs to use when sending messages to your instance of IBM Message Hub |
| MSGHUB_TOPIC | yes | string | The topic to use when sending messages to your instance of IBM Message Hub |

These **optional** input values can be overridden:

| Name | Required? | Type | Description |
| ---- | --------- | ---- | ---------------- |
| VERBOSE | no | integer | default is 0. Set to 1 to log everything that happens. |


#### Example:
A sample `services` section of the input file given to `hzn register`:
```
    "services": [
        {
            "org": "$HZN_ORG_ID",
            "url": "$SERVICE_NAME",
            "variables": {
                "MSGHUB_API_KEY": "$MSGHUB_API_KEY",
                "MSGHUB_BROKER_URL": "$MSGHUB_BROKER_URL",
            }
        }
    ]
```

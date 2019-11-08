# Horizon MQTT to IBM Event Streams Service

This services the lower level service mqtt. When something is published to the mqtt topic this service is subscribed to, it will receive it and send it to IBM Event Streams.

## Input Values

| Name | Required? | Type | Description |
| ---- | --------- | ---- | ---------------- |
| MQTT_WST_EVST | yes | string | Default is "mqtt_wst_evst". MQTT topic this service is subscribed to | 
| EVTSTREAMS_API_KEY | yes | string | The API key of the IBM Event Streams instance you are sending data to |
| EVTSTREAMS_BROKER_URL | yes | string | The comma-separated list of URLs to use when sending messages to your instance of IBM Event Streams |
| EVTSTREAMS_TOPIC | yes | string | Default is "myeventstreams". The topic to use when sending messages to your instance of IBM Event Streams |
| EVTSTREAMS_CERT_ENCODED | no | string | Default is "-". The base64-encoded self-signed certificate to use when sending messages to your ICP instance of  IBM Event Streams. Not needed for IBM Cloud Event Streams |

#### Example:
A sample `services` section of the input file given to `hzn register`:
```
{
    "services": [
        {
            "org": "$HZN_ORG_ID",
            "url": "$SERVICE_NAME",
            "variables": {
                "MQTT_WST_EVST": "$MQTT_WST_EVST",
                "EVTSTREAMS_API_KEY": "$EVTSTREAMS_API_KEY",
                "EVTSTREAMS_BROKER_URL": "$EVTSTREAMS_BROKER_URL",
                "EVTSTREAMS_TOPIC": "$EVTSTREAMS_TOPIC",
                "EVTSTREAMS_CERT_ENCODED": "$EVTSTREAMS_CERT_ENCODED"
            }
        }
    ]
}
```


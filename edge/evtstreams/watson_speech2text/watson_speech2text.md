# Horizon IBM Watson Speech to Text to IBM Event Streams Service

Service listens for hot work "Watson," once detected it captures an audio clip that is sent to an instance of Speech to Text, optionally removing stop words if user input is se to "true," then the transcribed text is send to the IBM Event Streams.

This services depends on four lower level services: `mqtt`, `mqtt2kafka`, `hotword_detection`, and `stopword_removal`. The hotword detection service is constantly listening for "Watson" and upon detection it will record a clip of audio and send it via the mqtt broker to the `watsons2text` service. 

The `watsons2text` service relies on the IBM Speech to Text service which requires an API Key and url to convert the audio clip to text. Once the audio clip is converted to text it will be passed to `stopword_removal` (if the environment variable `REMOVE_SW` is set to "true"), which is running as a WSGI server, where common stop words are removed and sent back to watson_speecdh2text. Finally, again via the `mqtt` broker, the text is sent to `mqtt2kafka` where it will be sent to IBM Event Streams. 

## Hardware Requirements 

This service was developed on, and designed for use with, a Raspberry Pi. For best results it is recommended to use a TROND External USB Audio Adapter Sound Card.

## Input Values for watsons2text service

The following input values **must** be given to this service in the input file given to `hzn register`:


| Name | Required? | Type | Description |
| ---- | --------- | ---- | ---------------- |
| STT_IAM_APIKEY | yes | string | The API key of the IBM Speech to Text instance |
| STT_URL | yes | string | Default is Dallas endpoint. Service endpoint url for Speech to Text service |

These **optional** input values can be overridden:


| Name | Required? | Type | Description |
| ---- | --------- | ---- | ---------------- |
| MQTT_HOST | yes | string | Default is ibm.mqtt. MQTT broker name | 
| MQTT_PORT | yes | string | Default is 1883. MQTT broker port | 
| MQTT_HWD_WST | yes | string | Default is mqtt_hwd_wst. Communication from hotworddetect to watsons2text | 
| MQTT_WST_EVST | yes | string | Default is mqtt_wst_evst. Communication from watsons2text to mqtt2kafka | 
| REMOVE_SW | yes | string | Flag to enable stopwordremoval service | 
| SW_HOST | yes | string | Default is ibm.stopwordremoval - Name of stop word removal endpoint | 
| SW_PORT | yes | string | Default is 5002. Port stop word removal WSGI server is listening on | 

## Input Values mqtt2kafka service

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
            "url": "ibm.watsons2text",
            "variables": {
                "MQTT_HOST": "$MQTT_HOST",
                "MQTT_PORT": "$MQTT_PORT",
                "MQTT_HWD_WST": "$MQTT_HWD_WST",
                "MQTT_WST_EVST": "$MQTT_WST_EVST",
                "REMOVE_SW": "$REMOVE_SW",
                "STT_IAM_APIKEY": "$STT_IAM_APIKEY",
                "STT_URL": "$STT_URL",
                "SW_HOST": "$SW_HOST",
                "SW_PORT": "$SW_PORT"
            }
        },
        {
            "org": "$HZN_ORG_ID",
            "url": "ibm.mqtt2kafka",
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


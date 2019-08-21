# Horizon IBM Watson Speech to Text to IBM Event Streams Service

Service listens for hot work "Watson," once detected it captures an audio clip that is sent to an instance of Speech to Text, optionally removing stop words if user input is se to "true," then the transcribed text is send to the IBM Event Streams. The topic is sends to is: `myeventstreams` (can be overridden)

This services depends on four lower level services: mqtt, mqtt2kafka, hotword_detection, and stopword_removal. The hotword detection service is constantly listening for "Watson" and upon detection it will record a clip of audio and send it via the mqtt broker to the watson_speecdh2text service. 

Watson_speecdh2text relies on the IBM Speech to Text service which requires an API Key and url to convert the audio clip to text. Once the audio clip is converted to text it will be passed to stopword_removal (if the environment variable REMOVE_SW is set to "true"), which is running as a WSGI server, where common stop words are removed and sent back to watson_speecdh2text. Finally, again via the mqtt broker, the text is sent to mqtt2kafka where it will be sent to IBM Event Streams. 

## Input Values

The following input values **must** be given to this service in the input file given to `hzn register`:


| Name | Required? | Type | Description |
| ---- | --------- | ---- | ---------------- |
| MSGHUB_API_KEY | yes | string | The API key of the IBM Event Streams instance you are sending data to |
| MSGHUB_BROKER_URL | yes | string | The comma-separated list of URLs to use when sending messages to your instance of IBM Event Streams |
| STT_IAM_APIKEY | yes | string | The API key of the IBM Speech to Text instance |
| STT_URL | yes | string | Default is Dallas endpoint. Service endpoint url for Speech to Text service |

These **optional** input values can be overridden:


| Name | Required? | Type | Description |
| ---- | --------- | ---- | ---------------- |
| MQTT_HOST | yes | string | Default is ibm.mqtt. MQTT broker name | 
| MQTT_PORT | yes | string | Default is 1883. MQTT broker port | 
| MQTT_HWD_WST | yes | string | Default is mqtt_hwd_wst. Communication from hotworddetct to watsons2text | 
| MQTT_WST_EVST | yes | string | Default is mqtt_wst_evst. Communication from watsons2text to mqtt2kafka | 
| REMOVE_SW | yes | string | Flag to enable stopwordremoval service | 
| SW_HOST | yes | string | Default is 127.0.0.1 - Name of stop word removal endpoint (ibm.stopwordremoval) |  
| SW_PORT | yes | string | Default is 5002. Port stop word removal WSGI server is listening on | 
| MQTT_WST_EVST | yes | string | Communication from watsons2text to mqtt2kafka | 
| MSGHUB_TOPIC | yes | string | Defauly is "myeventstreams." The topic to use when sending messages to your instance of IBM Event Streams |


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
                "MSGHUB_API_KEY": "$MSGHUB_API_KEY",
                "MSGHUB_BROKER_URL": "$MSGHUB_BROKER_URL",
                "MSGHUB_TOPIC": "$MSGHUB_TOPIC"
            }
        }
    ]
}
```


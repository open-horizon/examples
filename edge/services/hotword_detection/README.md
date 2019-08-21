# Horizon Hotword Detection Service

Service listens for hot work "Watson," once detected it will record a clip of audio and publish it to the mqtt broker.

## Input Values

| Name | Required? | Type | Description |
| ---- | --------- | ---- | ---------------- |
| MQTT_HOST | yes | string | Default is ibm.mqtt. MQTT broker name | 
| MQTT_PORT | yes | string | Default is 1883. MQTT broker port | 
| MQTT_AUDIO_TOPIC | yes | string | Default is mqtt_hwd_wst. The MQTT topic to publish the audio data |
| AUDIO_FORMAT | yes | string | Default is flac. This is the format the audio clip will be saved as |


#### Example:
A sample `services` section of the input file given to `hzn register`:
```
{
    "services": [
        {
            "org": "$HZN_ORG_ID",
            "url": "$SERVICE_NAME",
            "variables": {
                "MQTT_HOST": "$MQTT_HOST",
                "MQTT_PORT": "$MQTT_PORT",
                "MQTT_AUDIO_TOPIC": "$MQTT_AUDIO_TOPIC",
                "AUDIO_FORMAT": "$AUDIO_FORMAT"
            }
        }
    ]
}
```


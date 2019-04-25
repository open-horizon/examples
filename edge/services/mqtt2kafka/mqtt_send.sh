#!/bin/bash

# this is for local testing

echo "--- sending message: ${MQTT_MESSAGE} to topic: ${MQTT_TOPIC}"
mosquitto_pub -d -t ${MQTT_TOPIC} -m "${MQTT_MESSAGE}"
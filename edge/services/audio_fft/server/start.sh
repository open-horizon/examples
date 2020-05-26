#!/bin/sh

./server -b ${MQTT_BROKER} -u ${MQTT_SERVER_USER} -p ${MQTT_SERVER_PASS} -c ${MQTT_SERVER_CLIENT} -r ${SAMPLE_RATE} --result_topic ${MQTT_RESULTS_TOPIC}
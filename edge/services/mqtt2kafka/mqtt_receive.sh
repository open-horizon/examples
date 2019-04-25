#!/bin/bash

# kafka envs
# MSGHUB_API_KEY
# MSGHUB_BROKER_URL
# MSGHUB_TOPIC

if [[ -z "${MSGHUB_API_KEY}" ]]; then
    	echo "MSGHUB_API_KEY not set: ${MSGHUB_API_KEY}"
    	exit 1
fi

if [[ -z "${MSGHUB_BROKER_URL}" ]]; then
    	echo "MSGHUB_BROKER_URL not set: ${MSGHUB_BROKER_URL}"
    	exit 1
fi

if [[ -z "${MSGHUB_TOPIC}" ]]; then
    	echo "MSGHUB_TOPIC not set: ${MSGHUB_TOPIC}"
    	exit 1
fi

if [[ -z "${MQTT_TOPIC}" ]]; then
    MQTT_TOPIC=zhangl_mqtt_topic
    echo "receive from topic: ${MQTT_TOPIC}"
fi

MSGHUB_USERNAME="${MSGHUB_API_KEY:0:16}"
MSGHUB_PASSWORD="${MSGHUB_API_KEY:16}"
echo "Will publish to kafka topic: ${MSGHUB_TOPIC}"

trap "exit 130" INT
mosquitto_sub -h zhangl.mqtt -p 1883 -t ${MQTT_TOPIC} | while read; do
	if [ ! -z "${REPLY}" ]; then
		echo "text received from MQTT: ${REPLY}"
	    	KAFKA_JSON='{"nodeID":"'$HZN_DEVICE_ID'","text":"'${REPLY}'"}'
	    	echo "json to send to kafka: ${KAFKA_JSON}"
	    	echo "echo ${KAFKA_JSON} | kafkacat -P -b $MSGHUB_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=$MSGHUB_USERNAME -X sasl.password=$MSGHUB_PASSWORD -t $MSGHUB_TOPIC"
	    	echo "${KAFKA_JSON}" | kafkacat -P -b $MSGHUB_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=$MSGHUB_USERNAME -X sasl.password=$MSGHUB_PASSWORD -t $MSGHUB_TOPIC
	    	sleep 1
    	fi
done

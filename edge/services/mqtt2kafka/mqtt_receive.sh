#!/bin/bash

# Check the exit status of the previously run command and exit if nonzero (unless 'continue' is passed in)
checkrc() {
  if [[ $1 -ne 0 ]]; then
    echo "Error: exit code $1 from $2"
    # Sometimes it is useful to not exit on error, because if you do the container restarts so quickly it is hard to get in it a debug
    if [[ "$3" != "continue" ]]; then
      exit $1
    fi
  fi
}


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

if [[ -z "${MQTT_WST_EVST}" ]]; then
    echo "receive from topic: ${MQTT_WST_EVST}"
    exit 1
fi

MSGHUB_USERNAME="${MSGHUB_API_KEY:0:16}"
MSGHUB_PASSWORD="${MSGHUB_API_KEY:16}"
echo "Will publish to kafka topic: ${MSGHUB_TOPIC}"

if [[ -n "$MSGHUB_CERT_ENCODED" && "$MSGHUB_CERT_ENCODED" != "-" ]]; then
    # They are using an instance of Event Streams deployed in ICP, because it needs a self-signed cert
    MSGHUB_CERT_FILE=/tmp/es-cert.pem
    echo "$MSGHUB_CERT_ENCODED" | base64 -d > $MSGHUB_CERT_FILE
    checkrc $? "decode cert"
fi


trap "exit 130" INT
mosquitto_sub -h ibm.mqtt -p 1883 -t ${MQTT_WST_EVST} | while read; do
	if [ ! -z "${REPLY}" ]; then
		echo "text received from MQTT: ${REPLY}"
	    KAFKA_JSON='{"nodeID":"'$HZN_DEVICE_ID'","text":"'${REPLY}'"}'
	    echo "json to send to kafka: ${KAFKA_JSON}"

        if [[ -n "$MSGHUB_CERT_FILE"  ]]; then
            # Send data to ICP Event Streams
            echo "echo ${KAFKA_JSON} | kafkacat -P -b $MSGHUB_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$MSGHUB_API_KEY -X ssl.ca.location=$MSGHUB_CERT_FILE -t $MSGHUB_TOPIC"
            echo "${KAFKA_JSON}" | kafkacat -P -b $MSGHUB_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$MSGHUB_API_KEY -X ssl.ca.location=$MSGHUB_CERT_FILE -t $MSGHUB_TOPIC

        else
            # Send data to IBM Cloud Event Streams
            echo "echo ${KAFKA_JSON} | kafkacat -P -b $MSGHUB_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=$MSGHUB_USERNAME -X sasl.password=$MSGHUB_PASSWORD -t $MSGHUB_TOPIC"

            echo "${KAFKA_JSON}" | kafkacat -P -b $MSGHUB_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=$MSGHUB_USERNAME -X sasl.password=$MSGHUB_PASSWORD -t $MSGHUB_TOPIC
        fi

        sleep 1
    fi
done

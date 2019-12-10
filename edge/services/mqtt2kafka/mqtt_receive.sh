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

echo "Starting mqtt2kafka service..."
if [[ -z "${EVTSTREAMS_API_KEY}" ]]; then
    echo "EVTSTREAMS_API_KEY not set: ${EVTSTREAMS_API_KEY}"
    exit 1
fi

if [[ -z "${EVTSTREAMS_BROKER_URL}" ]]; then
    echo "EVTSTREAMS_BROKER_URL not set: ${EVTSTREAMS_BROKER_URL}"
    exit 1
fi

if [[ -z "${EVTSTREAMS_TOPIC}" ]]; then
    echo "EVTSTREAMS_TOPIC not set: ${EVTSTREAMS_TOPIC}"
    exit 1
fi

if [[ -z "${MQTT_WST_EVST}" ]]; then
    echo "receive from topic: ${MQTT_WST_EVST}"
    exit 1
fi

echo "Will publish to kafka topic: ${EVTSTREAMS_TOPIC}"

if [[ -n "$EVTSTREAMS_CERT_ENCODED" && "$EVTSTREAMS_CERT_ENCODED" != "-" ]]; then
    # They are using an instance of Event Streams deployed in ICP, because it needs a self-signed cert
    EVTSTREAMS_CERT_FILE=/tmp/es-cert.pem
    echo "$EVTSTREAMS_CERT_ENCODED" | base64 -d > $EVTSTREAMS_CERT_FILE
    checkrc $? "decode cert"
fi


trap "exit 130" INT
mosquitto_sub -h ibm.mqtt -p 1883 -t ${MQTT_WST_EVST} | while read; do
	if [ ! -z "${REPLY}" ]; then
		echo "text received from MQTT: ${REPLY}"
	    KAFKA_JSON='{"nodeID":"'$HZN_DEVICE_ID'","text":"'${REPLY}'"}'
	    echo "json to send to kafka: ${KAFKA_JSON}"

            # Send data to IBM Event Streams
            echo "echo ${KAFKA_JSON} | kafkacat -P -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$EVTSTREAMS_API_KEY -X ssl.ca.location=$EVTSTREAMS_CERT_FILE -t $EVTSTREAMS_TOPIC"
            echo "${KAFKA_JSON}" | kafkacat -P -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$EVTSTREAMS_API_KEY -X ssl.ca.location=$EVTSTREAMS_CERT_FILE -t $EVTSTREAMS_TOPIC

        sleep 1
    fi
done

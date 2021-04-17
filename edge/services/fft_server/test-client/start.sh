#!/bin/sh

set -e

echo "`/sbin/ip route|grep default | cut -d ' ' -f3` host.docker.internal" | tee -a /etc/hosts > /dev/null

./fft-test -b ${MQTT_BROKER} -u ${MQTT_SERVER_USER} -p ${MQTT_SERVER_PASS} -c ${MQTT_SERVER_CLIENT} --result_topic ${MQTT_RESULTS_TOPIC}
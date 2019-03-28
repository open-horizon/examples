#!/bin/sh

export SPEEDTEST_SERVICE_URL="172.17.0.2/v1/speedtest"
export SPEEDTEST_PAUSE_SEC=10

export WIOTP_ORG=`cat wiotp.config | jq -r '.config."wiotp-org"'`
export WIOTP_DEV_TYPE=`cat wiotp.config | jq -r '.config."wiotp-device-type"'`
export WIOTP_DEV_ID=`cat wiotp.config | jq -r '.config."wiotp-device-id"'`
export WIOTP_DEV_TOKEN=`cat wiotp.config | jq -r '.config."wiotp-device-token"'`


while true; do
  
  result=`curl -sS "${SPEEDTEST_SERVICE_URL}"`
  # echo "${result}"

  mosquitto_pub -h "${WIOTP_ORG}.messaging.internetofthings.ibmcloud.com" -p 8883 -i "d:${WIOTP_ORG}:${WIOTP_DEV_TYPE}:${WIOTP_DEV_ID}" -u "use-token-auth" -P "${WIOTP_DEV_TOKEN}" --capath /etc/ssl/certs -t iot-2/evt/status/fmt/json -m "${result}" -d

  sleep ${SPEEDTEST_PAUSE_SEC}
done


#!/bin/sh

export SPEEDTEST_SERVICE_URL="http://speedtest/v1/speedtest"
export SPEEDTEST_PAUSE_SEC=10

while true; do
  
  result=`curl -sS "${SPEEDTEST_SERVICE_URL}"`
  #echo "${result}"

  e=`echo "${result}" | jq '.error'`
  if [ "null" != "${e}" ]; then
    echo "Waiting for \"speedtest\" service to warm up..."
  else
    mosquitto_pub -h "${WIOTP_ORG}.messaging.internetofthings.ibmcloud.com" -p 8883 -i "d:${WIOTP_ORG}:${WIOTP_DEVICE_TYPE}:${WIOTP_DEVICE_ID}" -u "use-token-auth" -P "${WIOTP_DEVICE_TOKEN}" --capath /etc/ssl/certs -t iot-2/evt/status/fmt/json -m "${result}" -d
  fi

  sleep ${SPEEDTEST_PAUSE_SEC}
done


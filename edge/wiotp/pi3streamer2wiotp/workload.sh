#!/bin/sh

# Horizon sample workload to query the status from a sample microservice hosting a local web UI, and publish status to Watson IoT Platform
# This workload expects the microservice to be running, unless it is running in mock mode.

# Verify required environment variables are set
checkRequiredEnvVar() {
  varname=$1
  if [ -z $(eval echo \$$varname) ]; then
    echo "Error: Environment variable $varname must be set; exiting."
    exit 2
  else
    echo "  $varname=" $(eval echo \$$varname)
  fi
}

# Environment variables that can optionally be set, or default
WIOTP_DOMAIN="${WIOTP_DOMAIN:-internetofthings.ibmcloud.com}"     # set in the pattern deployment_overrides field if you need to override
WIOTP_PEM_FILE="${WIOTP_PEM_FILE:-messaging.pem}"     # the cert to verify the WIoTP MQTT broker
PUBLISH_INTERVAL="${PUBLISH_INTERVAL:-5}"    # reporting interval in seconds
PUBLISH="${PUBLISH:-true}"    # whether or not to actually send data to wiotp
MOCK="${MOCK:-false}"     # if "true", just pretend to call the microservice
VERBOSE="${VERBOSE:-0}"    # set to 1 for verbose output

echo "Optional environment variables (or default values): WIOTP_DOMAIN=$WIOTP_DOMAIN, WIOTP_PEM_FILE=$WIOTP_PEM_FILE, WIOTP_EDGE_MQTT_IP=$WIOTP_EDGE_MQTT_IP, PUBLISH_INTERVAL=$PUBLISH_INTERVAL, PUBLISH=$PUBLISH, MOCK=$MOCK"

# When this workload is running in standalone mode, there are no required env vars.
if [[ "$PUBLISH" == "true" ]]; then
  echo "Checking for required environment variables are set:"
  checkRequiredEnvVar "HZN_ORGANIZATION"      # automatically passed in by Horizon
  checkRequiredEnvVar "HZN_DEVICE_ID"      # automatically passed in by Horizon. Wiotp automatically gives this a value of: g@mygwtype@mygw
fi

# Parse the class id, device type, and device id from HZN_DEVICE_ID. It will have a value like 'g@mygwtype@mygw'
if [[ ! -z "$HZN_DEVICE_ID" ]]; then
  id="$HZN_DEVICE_ID"
  CLASS_ID=${id%%@*}   # the class id is not actually used anymore
  id=${id#*@}
  DEVICE_TYPE=${id%%@*}
  DEVICE_ID=${id#*@}
  if [[ -z "$DEVICE_TYPE" || -z "$DEVICE_ID" ]]; then
    echo 'Error: HZN_DEVICE_ID must have the format: g@mygwtype@mygw'
    exit 2
  fi

  if [[ "$VERBOSE" == 1 ]]; then echo "  DEVICE_TYPE=$DEVICE_TYPE, DEVICE_ID=$DEVICE_ID"; fi
fi

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

echo 'Starting infinite loop to read from microservice then publish...'

while true; do
  # Get data from a local microservice
  status="active"
  if [[ "$MOCK" == "true" ]]; then
    output='{"ts":'$(date +%S)'} 200'
    status="mock"
    curlrc=0
  else
    output=$(curl -sS -w %{http_code} "http://pi3streamer:8080/static_simple.html")
    curlrc=$?     # save this before it gets overwritten
  fi

  httpcode=$(echo "$output" | tail -c 4) # last 3 characters (digits)
  json=$(echo "$output" | grep MJPG)     # for the output, get the string referencing the streamer 

  echo "httpcode = $httpcode... payload = $json"

  if [[ "$curlrc" != 0 ]]; then
    echo "Warning: Curl command to the local microservice returned exit code $curlrc, will try again next interval."
    status="inactive"
    # Don't do anything else...
  elif [[ "$httpcode" != 200 ]]; then
    echo "Warning: HTTP code $httpcode from the local microservice, will try again next interval."
  else
    status="active"
    # Report status to WIoTP
    json='{"pi3streamer":"'$status'", "type":"'$json'"}'

    if [[ "$PUBLISH" == "true" ]]; then
      # Send a "status" event to the Watson IoT Platform containing the data
      clientId="a:$HZN_ORGANIZATION:$DEVICE_TYPE$DEVICE_ID"       # sending as an app
      topic="iot-2/evt/status/fmt/json"
      if [[ -n "$WIOTP_EDGE_MQTT_IP" ]]; then
        # Send to the local WIoTP Edge Connector microservice mqtt broker, so it can store and forward
        msgHost="$WIOTP_EDGE_MQTT_IP"
      else
        # Send directly to the WIoTP cloud mqtt broker
        msgHost="$HZN_ORGANIZATION.messaging.$WIOTP_DOMAIN"
      fi

      # Send as App, to local edge-connector
      echo mosquitto_pub -h "$msgHost" -p 8883 -i "$clientId" --cafile $WIOTP_PEM_FILE -q 1 -t "$topic" -m "$json"
      mosquitto_pub -h "$msgHost" -p 8883 -i "$clientId" --cafile $WIOTP_PEM_FILE -q 1 -t "$topic" -m "$json" >/dev/null
      checkrc $? "mosquitto_pub $msgHost" "continue"
    fi
  fi

  # Pause before looping again
  sleep $PUBLISH_INTERVAL
done
# Not reached

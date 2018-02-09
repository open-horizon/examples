#!/bin/sh

# Horizon workload to query a Microservice and publish to Watson IoT Platform

# This workload expects the CPU microservice to be running.  Run 'make' in the
# sibling directory "../microservice" to start that microservice running.  You
# can verify that microservcie is working by running 'make check' here.

# Check the exit status of the previously run command and exit if nonzero (unless 'continue' is passed in)
checkrc() {
  if [[ $1 -ne 0 ]]; then
    echo "ERROR: exit code $1 from $2"
    # Sometimes it is useful to not exit on error, because if you do the container restarts so quickly it is hard to get in it a debug
    if [[ "$3" != "continue" ]]; then
      exit $1
    fi
  fi
}

# Verify required configuration and credentials are in the process environment
checkRequiredEnvVar() {
  varname=$1
  if [ -z $(eval echo \$$varname) ]; then
    echo "ERROR: Environment variable $varname must be set; exiting."
    exit 1
  else
    echo "  $varname=" $(eval echo \$$varname)
  fi
}

echo "Checking for required configuration from the process environment:"

checkRequiredEnvVar "SAMPLE_SIZE"
checkRequiredEnvVar "SAMPLE_INTERVAL"
checkRequiredEnvVar "MOCK"
checkRequiredEnvVar "PUBLISH"
checkRequiredEnvVar "VERBOSE"

if test "$PUBLISH" = "true"; then
  checkRequiredEnvVar "HZN_ORGANIZATION"      # automatically passed in by Horizon
  checkRequiredEnvVar "WIOTP_DEVICE_AUTH_TOKEN"   # a userInput value, so must be set in the input file passed to 'hzn register'
  checkRequiredEnvVar "HZN_DEVICE_ID"      # automatically passed in by Horizon

  # For device type and device id, there are 3 cases, checked for below:
  #  1) The workload is run in wiotp/horizon and will be sending mqtt as the gw: HZN_DEVICE_ID contains class id, device type and id, and WIOTP_CLASS_ID, WIOTP_DEVICE_TYPE and WIOTP_DEVICE_ID are blank or '-'
  #  2) The workload is run in wiotp/horizon and will be sending mqtt as a device: WIOTP_CLASS_ID, WIOTP_DEVICE_TYPE and/or WIOTP_DEVICE_ID have real values
  #  3) The workload is run in a non-wiotp horizon instance: HZN_DEVICE_ID is the simple device id (use that), and WIOTP_CLASS_ID and WIOTP_DEVICE_TYPE have values
  # The way we will handle this cases is if WIOTP_CLASS_ID, WIOTP_DEVICE_TYPE and/or WIOTP_DEVICE_ID are set, they will override what we can parse from HZN_DEVICE_ID
  if [[ "$WIOTP_CLASS_ID" != "-" ]]; then
    CLASS_ID="$WIOTP_CLASS_ID"
  fi
  if [[ "$WIOTP_DEVICE_TYPE" != "-" ]]; then
    DEVICE_TYPE="$WIOTP_DEVICE_TYPE"
  fi
  if [[ "$WIOTP_DEVICE_ID" != "-" ]]; then
    DEVICE_ID="$WIOTP_DEVICE_ID"
  fi
  if [[ "$HZN_DEVICE_ID" != "${HZN_DEVICE_ID#*@}" ]]; then
    # When this workload is deployed by WIoTP-Horizon, HZN_DEVICE_ID will have a value like 'g@mygwtype@mygw'. Parse it.
    id="$HZN_DEVICE_ID"
    classId=${id%%@*}
    id=${id#*@}
    deviceType=${id%%@*}
    deviceId=${id#*@}
    if [[ -z "$CLASS_ID" ]]; then
      CLASS_ID="$classId"
    fi
    if [[ -z "$DEVICE_TYPE" ]]; then
      DEVICE_TYPE="$deviceType"
    fi
    if [[ -z "$DEVICE_ID" ]]; then
      DEVICE_ID="$deviceId"
    fi
  else
    # When this workload is run in a non-wiotp horizon instance, HZN_DEVICE_ID will be a simple device id
    if [[ -z "$DEVICE_ID" ]]; then
      DEVICE_ID="$HZN_DEVICE_ID"
    fi
  fi
  if [[ -z "$CLASS_ID" ]]; then
    echo "ERROR: class id not set in WIOTP_CLASS_ID or HZN_DEVICE_ID."
    exit 1
  elif [[ "$CLASS_ID" != "d" && "$CLASS_ID" != "g" ]]; then
    echo "ERROR: class id in WIOTP_CLASS_ID or HZN_DEVICE_ID can only have the value 'd' or 'g'."
    exit 1
  fi
  if [[ -z "$DEVICE_TYPE" ]]; then
    echo "ERROR: device type not set in WIOTP_DEVICE_TYPE or HZN_DEVICE_ID."
    exit 1
  fi
  if [[ -z "$DEVICE_ID" ]]; then
    echo "ERROR: device id not set in WIOTP_DEVICE_ID or HZN_DEVICE_ID."
    exit 1
  fi
  echo "Optional override environment variables:"
  echo "  WIOTP_CLASS_ID=$WIOTP_CLASS_ID"
  echo "  WIOTP_DEVICE_TYPE=$WIOTP_DEVICE_TYPE"
  echo "  WIOTP_DEVICE_ID=$WIOTP_DEVICE_ID"
  echo "Derived variables:"
  echo "  CLASS_ID=$CLASS_ID"
  echo "  DEVICE_TYPE=$DEVICE_TYPE"
  echo "  DEVICE_ID=$DEVICE_ID"

  # Environment variables that can optionally be set, or default
  WIOTP_DOMAIN="${WIOTP_DOMAIN:-internetofthings.ibmcloud.com}"     # set in the pattern deployment_overrides field if you need to override
  # WIOTP_API_KEY: API key to use to create the WIoTP device if it doesn't already exist
  # WIOTP_API_AUTH_TOKEN: API token to use to create the WIoTP device if it doesn't already exist
  WIOTP_PEM_FILE="${WIOTP_PEM_FILE:-messaging.pem}"     # the cert to verify the WIoTP MQTT cloud broker
  # WIOTP_EDGE_MQTT_IP: to a local IP or hostname to send mqtt msgs via the WIoTP Edge Connector microservice (enables severability)

  echo "Optional environment variables (or default values):"
  echo "  WIOTP_DOMAIN=$WIOTP_DOMAIN"
  echo "  WIOTP_API_KEY=$WIOTP_API_KEY"
  echo "  WIOTP_API_AUTH_TOKEN=$WIOTP_API_AUTH_TOKEN"
  echo "  WIOTP_PEM_FILE=$WIOTP_PEM_FILE"
  echo "  WIOTP_EDGE_MQTT_IP=$WIOTP_EDGE_MQTT_IP"

  # If Watson IoT Platform API credentials are not provided assume existence.
  if [[ -z "$WIOTP_API_KEY" || -z "$WIOTP_API_AUTH_TOKEN" ]]; then
    echo "Watson IoT Platfrom REST API credentials WIOTP_API_KEY and WIOTP_API_AUTH_TOKEN were not provided,"
    echo "assuming type \"$DEVICE_TYPE\" with ID \"$DEVICE_ID\" already exists in Watson IoT Platform."
  else
    # Both credentials provided; prepare for Watson IoT Platform REST API calls
    echo "API credentials successfully received from process environment."
    copts='-sS -w %{http_code}'
    wiotpApiAuth="$WIOTP_API_KEY:$WIOTP_API_AUTH_TOKEN"
    apiUrl="https://$HZN_ORGANIZATION.$WIOTP_DOMAIN/api/v0002"
    contentJson='Content-Type: application/json'

    # Verify the specified DEVICE_TYPE exists and if not, exit.
    echo "Checking whether specified WIoTP Device Type exists..."
    httpCode=$(curl $copts -u "$wiotpApiAuth" -o /dev/null $apiUrl/device/types/$DEVICE_TYPE)
    checkrc $? "curl $apiUrl/device/types/$DEVICE_TYPE"
    if [[ "$httpCode" == "404" ]]; then
      echo "Watson IoT device Type \"$DEVICE_TYPE\" does not exist."
      exit 1
    fi
    echo "Device Type \"$DEVICE_TYPE\" exists in Watson IoT Platform."

    # Does the specified DEVICE_ID exist?  If not, create it.
    echo "Checking whether specified WIoTP Device ID exists..."
    httpCode=$(curl $copts -u "$wiotpApiAuth" -o /dev/null $apiUrl/device/types/$DEVICE_TYPE/devices/$DEVICE_ID)
    checkrc $? "curl $apiUrl/device/types/$DEVICE_TYPE/devices/$DEVICE_ID"
    if [[ "$httpCode" == "404" ]]; then
      echo "Creating device \"$DEVICE_ID\" in Watson IoT Platform..."
      body='{"deviceId":"'$DEVICE_ID'", "authToken":"'$WIOTP_DEVICE_TOKEN'", "deviceInfo":{"description":"My edge device"}}, "metadata":{}}'
      output=$(curl $copts -u "$wiotpApiAuth" -X POST -H "$contentJson" -d "$body" $apiUrl/device/types/$DEVICE_TYPE/devices)
      checkrc $? "curl $apiUrl/device/types/$DEVICE_TYPE/devices"
      httpCode=${output:$((${#output}-3))} # last 3 chars are http status code
      if [[ "$httpCode" != "201" ]]; then
        echo "ERROR: Failed to create device $DEVICE_ID: $output"
        exit 1
      fi
    elif [[ "$httpCode" != "200" ]]; then
      echo "ERROR: HTTP code $httpCode was returned when trying to check for device \"$DEVICE_ID\". Exiting..."
      exit 1
    fi
    echo "Device \"$DEVICE_ID\" exists in Watson IoT Platform."
  fi
fi

echo 'Starting infinite "read-from-microservice-then-publish" loop...'
sum=0
samples="$SAMPLE_SIZE"
samplecount=0
while true; do

  samplecount=$((samplecount + 1))

  # Get data from a local microservice
  if [[ "$MOCK" == "true" ]]; then
    output='{"cpu":1.2} 200'
    curlrc=0
  else
    output=$(curl -sS -w %{http_code} -m 15 "http://cpu:8347/v1/cpu")
    curlrc=$?     # save this before it gets overwritten
  fi
  httpcode=${output:$((${#output}-3))}    # the last 3 chars are the http code
  json="${output%?[0-9][0-9][0-9]}"   # for the output, get all but the newline and 3 digits of http code

  if [[ "$curlrc" != 0 ]]; then
    echo "Warning: Curl command to the local cpu microservice returned exit code $curlrc, will try again next interval."
  elif [[ "$httpcode" != 200 ]]; then
    echo "Warning: HTTP code $httpcode from the local cpu microservice REST API, will try again next interval."
  else

    # Accumulate the CPU usage and calculate the average after obtaining all samples.
    cpuusage=$(echo $json | jq '.cpu')
    if test "$VERBOSE" = "1"; then
      echo -e "Interval $samplecount usage: $cpuusage"
    fi
    sum=$(echo $sum + $cpuusage | bc)

    if test "$samplecount" -eq "$samples"; then
      average=$(echo $sum/$samples | bc -l)

      json='{"cpu":'$average'}'
      echo -e "avg: $json"

      if test "$PUBLISH" = "true"; then
        # Send a "status" event to the Watson IoT Platform containing the data
        clientId="$CLASS_ID:$HZN_ORGANIZATION:$DEVICE_TYPE:$DEVICE_ID"
        if [[ "$CLASS_ID" == "g" ]]; then
          topic="iot-2/type/$DEVICE_TYPE/id/$DEVICE_ID/evt/status/fmt/json"
        else
          topic="iot-2/evt/status/fmt/json"
        fi
        if [[ -n "$WIOTP_EDGE_MQTT_IP" ]]; then
          # Send to the local WIoTP Edge Connector microservice mqtt broker, so it can store and forward
          msgHost="$WIOTP_EDGE_MQTT_IP"
        else
          # Send directly to the WIoTP cloud mqtt broker
          msgHost="$HZN_ORGANIZATION.messaging.$WIOTP_DOMAIN"
        fi

        if [[ "$VERBOSE" == "1" ]]; then
          echo mosquitto_pub -h "$msgHost" -p 8883 -i "$clientId" -u "use-token-auth" -P "$WIOTP_DEVICE_AUTH_TOKEN" --cafile $WIOTP_PEM_FILE -q 2 -t "$topic" -m "$json"
        fi
        mosquitto_pub -h "$msgHost" -p 8883 -i "$clientId" -u "use-token-auth" -P "$WIOTP_DEVICE_AUTH_TOKEN" --cafile $WIOTP_PEM_FILE -q 2 -t "$topic" -m "$json" >/dev/null
        checkrc $? "mosquitto_pub $msgHost" "continue"
      fi
      sum=0
      samplecount=0
    fi
  fi

  # Pause before looping again
  sleep "$SAMPLE_INTERVAL"
done
# Not reached

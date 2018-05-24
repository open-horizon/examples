#!/bin/sh
#
# Blue Horizon location workload main program.
#
# Original written by Dima Rekesh
# Modified by Glen Darling and Bruce Potter
#
# Note: since this runs in alpine, it is the busybox shell, which has subtle differences from bash, including:
#  - String variables in if [[ ]] statements must be surrounded by ""
#  - Variable/parameter expansion has slightly different syntax in some cases

# Import the support functions for the central WIoTP MQTT service. May bring this back, have not yet decided...
#. wiotp.sh

# Check the exit status of the previously run command and exit if nonzero (unless 'continue' is passed in)
checkrc() {
  if [[ $1 -ne 0 ]]; then
    echo "Error: exit code $1 from $2"
    # Sometimes it is useful to not exit on error, because if you do the container restarts so quickly it is hard to get in it a debug
    if [[ "$3" != "continue" ]]; then
      exit $1
    else
      error="Error: exit code $1 from $2"
    fi
  else
    error=''
  fi
}

# Function to verify that a required environment variable is set
checkRequiredEnvVar() {
  varname=$1
  if [ -z $(eval echo \$$varname) ]; then
    echo "Error: Environment variable $varname must be set; exiting."
    exit 2
  else
    echo "  $varname="$(eval echo \$$varname)
  fi
}

# Required environment variables (must be provided; no default values given)
echo "Checking that required environment variables are set."
checkRequiredEnvVar "HZN_ORGANIZATION"
checkRequiredEnvVar "HZN_DEVICE_ID"
checkRequiredEnvVar "HZN_AGREEMENTID"

# Optional common environment variables (these have default values)
GPS_HOST_PORT="${GPS_HOST_PORT:-gps:31779}"     # The hostname and port number we should contact the gps rest api with
MQTT_PORT="${MQTT_PORT:-8883}"
# Configure how long to pause between successive publications to central MQTT
export REPORTING_INTERVAL=${REPORTING_INTERVAL:-20}
# If the gps data is the same, only send it to verne after this many readings
export SKIP_NUM_REPEAT_LOC_READINGS=${SKIP_NUM_REPEAT_LOC_READINGS:-6}
export SKIP_NUM_REPEAT_SAT_READINGS=${SKIP_NUM_REPEAT_SAT_READINGS:-6}
echo "Common, optional, environment variables:"
echo "  GPS_HOST_PORT=$GPS_HOST_PORT"
echo "  MQTT_PORT=$MQTT_PORT"
echo "  REPORTING_INTERVAL=$REPORTING_INTERVAL"
echo "  SKIP_NUM_REPEAT_LOC_READINGS=$SKIP_NUM_REPEAT_LOC_READINGS"
echo "  SKIP_NUM_REPEAT_SAT_READINGS=$SKIP_NUM_REPEAT_SAT_READINGS"
echo "  Static HZN_LAT=$HZN_LAT"
echo "  Static HZN_LON=$HZN_LON"

# Parse HZN_DEVICE_ID (e.g., "g@gwtype@gwid") for device type, and device id
if [[ "$HZN_DEVICE_ID" != "${HZN_DEVICE_ID#?*@?*@?}" ]]; then
    # We matched that pattern, so now parse it
    id="$HZN_DEVICE_ID"
    CLASS_ID=${id%%@*} # (ignored) class id is not used here
    id=${id#*@}
    GW_TYPE=${id%%@*}
    GW_ID=${id#*@}
    if [[ "$VERBOSE" == 1 ]]; then echo "[VERBOSE] Parsed values: GW_TYPE='$GW_TYPE', GW_ID='$GW_ID'"; fi
else
    echo 'Error: HZN_DEVICE_ID must have the format: g@mygwtype@mygw'
    exit 2
fi

WIOTP_DOMAIN="${WIOTP_DOMAIN:-internetofthings.ibmcloud.com}"     # set in the pattern deployment_overrides field if you need to override
if [[ -n "$WIOTP_EDGE_MQTT_IP" ]]; then
  # Send to the local WIoTP Edge Connector, for store and forward
  MQTT_HOST="$WIOTP_EDGE_MQTT_IP"
  WIOTP_PEM_FILE="${WIOTP_PEM_FILE:-/var/wiotp-edge/persist/dc/ca/ca.pem}"     # The cert to verify the WIoTP MQTT broker we are sending to
  CLIENT_ID="a:${HZN_AGREEMENTID:0:36}"      # sending as an app - wiotp limit for app id is 36
  TOPIC="iot-2/evt/status/fmt/json"
  #MQTT_AUTH=""    # no auth needed for edge-connector if you are on the same machine as it
else
  # Send directly to the WIoTP cloud mqtt broker. The gw token is required to authenticate this send.
  if [[ -z "$WIOTP_GW_TOKEN" || "$WIOTP_GW_TOKEN" == "-" ]]; then
    echo "Error: WIOTP_GW_TOKEN must be set to send messages directly to the WIoTP cloud MQTT broker. Exiting."
    exit 2
  fi
  MQTT_HOST="$HZN_ORGANIZATION.messaging.$WIOTP_DOMAIN"
  WIOTP_PEM_FILE="${WIOTP_PEM_FILE:-messaging.pem}"     # The cert to verify the WIoTP MQTT broker we are sending to
  CLIENT_ID="$CLASS_ID:$HZN_ORGANIZATION:$GW_TYPE:$GW_ID"
  TOPIC="iot-2/type/$GW_TYPE/id/$GW_ID/evt/status/fmt/json"
  #MQTT_AUTH="-u use-token-auth -P $WIOTP_GW_TOKEN"  # <- cant do this because need to put "" around the token (special chars) when passing it to mosquitto_pub
fi
echo "Environment variables (or default values) for sending to WIoTP:" 
echo "  WIOTP_DOMAIN=$WIOTP_DOMAIN"
echo "  WIOTP_PEM_FILE=$WIOTP_PEM_FILE"
echo "  WIOTP_EDGE_MQTT_IP='$WIOTP_EDGE_MQTT_IP' (optional, with no default)"
echo "  MQTT_HOST=$MQTT_HOST"
echo "  CLIENT_ID=$CLIENT_ID"
echo "  TOPIC=$TOPIC"
echo "  WIOTP_GW_TOKEN=$WIOTP_GW_TOKEN"

# Invoke the curl cmd
docurl() {
    out=$(curl -sS -w "%{http_code}" $* 2>&1)
    rc=$?
    # Note: can not separate the http code into a global var because this function gets run in a sub-shell, so the caller has to separate 
    echo "$out"
    return $rc
}

echo "Sending data to Watson IoT Platform..."

# Configure for the Blue Horizon "gps" microservice
GPS_MICROSERVICE_BASE_URI="http://$GPS_HOST_PORT/v1/gps"
GPS_LOCATION_URI="$GPS_MICROSERVICE_BASE_URI/location"
GPS_SATELLITES_URI="$GPS_MICROSERVICE_BASE_URI/satellites"

# Avoid sending the data every read interval if it is not changing
previous_loc=''
num_repeat_loc_readings=0
previous_sat=''
num_repeat_sat_readings=0

# Stream the location data
echo "Streaming location (and satellite, if any) data..."
while [ true ]
do

    # Location data
    # -------------

    # Get the location data from the Blue Horizon "gps" REST microservice
    output=$(docurl $GPS_LOCATION_URI)
    curlrc=$?
    httpcode=${output:$((${#output}-3))}    # the last 3 chars are the http code
    location_rest="${output%[0-9][0-9][0-9]}"   # for the output, get all but the last 3 digits

    # Did the REST request from the gps microservice succeed? If so, get the loc and ts vars
    echo "$location_rest" | grep -q -E '^\{"latitude":'
    if [[ $? -eq 0 && $curlrc -eq 0 && $httpcode -eq 200 ]]; then
        # GPS rest call succeeded, convert the JSON to the format expected by the central MQTT service
        # Ultimately we want the location json to look like this: {"t":1489681081,"r":{"lat":42.052346746,"lon":-73.960373724,"alt":52.662}}
        ts=$(echo $location_rest | jq -c -r '.loc_last_update')
        loc=$(echo $location_rest | jq -c '{ lat: .latitude, lon: .longitude, alt: .elevation }')
        if [[ "$VERBOSE" == 1 ]]; then echo "[VERBOSE] from GPS: ts='$ts', loc='$loc'"; fi
    else
        # Log the gps REST service failure, and then try to continue with the lat/lon env vars
        echo "ERROR: REST request to '$GPS_LOCATION_URI' microservice failed: rc=$curlrc, httpcode=$httpcode, output="$location_rest
        if [[ -n "$HZN_LAT" && -n "$HZN_LON" ]]; then
            echo "Continuing with location using the HZN_LAT and HZN_LON environment variables..."
            ts=`date +%s`
            loc="{\"lat\":$HZN_LAT,\"lon\":$HZN_LON,\"alt\":0}"
        else
            loc=''
        fi
    fi

    # Send the ts and loc, unless they have not changed recently
    if [[ -n "$loc" ]]; then
        if [[ "$loc" == "$previous_loc" && $num_repeat_loc_readings -lt $SKIP_NUM_REPEAT_LOC_READINGS ]]; then
            # Avoid constantly sending the same location. Wait until we have skipped SKIP_NUM_REPEAT_LOC_READINGS.
            if [[ "$VERBOSE" == 1 ]]; then echo "[VERBOSE] skipping this location reading, because it is the same and we have only skipped $num_repeat_loc_readings readings."; fi
            num_repeat_loc_readings=$(($num_repeat_loc_readings+1))
        else
            # This is a different location or we have waited long enough to report
            location="{\"t\":$ts,\"r\":$loc}"

            # Send to WIoTP
            if [[ -n "$WIOTP_EDGE_MQTT_IP" ]]; then
                # Sending via the local core-iot edge-connector, as an app
                echo mosquitto_pub -h "$MQTT_HOST" -p $MQTT_PORT -i "$CLIENT_ID" --cafile $WIOTP_PEM_FILE -q 1 -t "$TOPIC" -m "$location"
                mosquitto_pub -h "$MQTT_HOST" -p $MQTT_PORT -i "$CLIENT_ID" --cafile $WIOTP_PEM_FILE -q 1 -t "$TOPIC" -m "$location" >/dev/null
                checkrc $? "mosquitto_pub $MQTT_HOST" "continue"
            else
                # Sending straight to the wiotp cloud broker, as the gw 
                echo mosquitto_pub -h "$MQTT_HOST" -p $MQTT_PORT -i "$CLIENT_ID" -u 'use-token-auth' -P "$WIOTP_GW_TOKEN" --cafile $WIOTP_PEM_FILE -q 1 -t "$TOPIC" -m "$location"
                mosquitto_pub -h "$MQTT_HOST" -p $MQTT_PORT -i "$CLIENT_ID" -u 'use-token-auth' -P "$WIOTP_GW_TOKEN" --cafile $WIOTP_PEM_FILE -q 1 -t "$TOPIC" -m "$location" >/dev/null
                checkrc $? "mosquitto_pub $MQTT_HOST" "continue"
            fi
            #wiotp_pub $MQTT_HOST "$clientId" "$topic" $location 1 -r

            # NOTE: Sending is best effort only, if there is a problem send will display the error, but then we continue
            if [[ -z "$error" ]]; then
                # The send was successful, so reset the repeat variables
                previous_loc="$loc"
                num_repeat_loc_readings=0
            fi
        fi
    fi

    # Satellite data
    # -------------

    # Get the satellites data from the Blue Horizon "gps" REST microservice
    output=$(docurl $GPS_SATELLITES_URI)
    curlrc=$?
    httpcode=${output:$((${#output}-3))}    # the last 3 chars are the http code
    satellites_rest="${output%[0-9][0-9][0-9]}"   # for the output, get all but the last 3 digits

    # Did the REST request from the gps microservice succeed?
    echo "$satellites_rest" | grep -q -E '^\{"satellites":'
    if [[ $? -eq 0 && $curlrc -eq 0 && $httpcode -eq 200 ]]; then
        # GPS rest call succeeded, convert the JSON to the format expected by the central MQTT service
        # we want it to look like: {"t":1489681261,"d":[{"PRN":1,"el":10,"az":277,"ss":13,"used":false}]}
        ts=`date +%s`
        sats=$(echo $satellites_rest | jq -c '.satellites')
        if [[ "$VERBOSE" == 1 ]]; then echo "[VERBOSE] from GPS: ts='$ts', sats='$sats'"; fi

        # Only if there is actual satellite data, send it to the central Blue Horizon MQTT service
        if [[ "$sats" != 'null' && "$sats" != '[]' ]]; then
            if [[ "$sats" == "$previous_sat" && $num_repeat_sat_readings -lt $SKIP_NUM_REPEAT_SAT_READINGS ]]; then
                # Avoid constantly sending the same satellite. Wait until we have skipped SKIP_NUM_REPEAT_SAT_READINGS.
                if [[ "$VERBOSE" == 1 ]]; then echo "[VERBOSE] skipping this satellite reading, because it is the same and we have only skipped $num_repeat_sat_readings readings."; fi
                num_repeat_sat_readings=$(($num_repeat_sat_readings+1))
            else
                satellites="{\"t\":$ts,\"d\":$sats}"

                # Send to WIoTP
                if [[ -n "$WIOTP_EDGE_MQTT_IP" ]]; then
                    # Sending via the local core-iot edge-connector, as an app
                    echo mosquitto_pub -h "$MQTT_HOST" -p $MQTT_PORT -i "$CLIENT_ID" --cafile $WIOTP_PEM_FILE -q 1 -t "$TOPIC" -m "$satellites"
                    mosquitto_pub -h "$MQTT_HOST" -p $MQTT_PORT -i "$CLIENT_ID" --cafile $WIOTP_PEM_FILE -q 1 -t "$TOPIC" -m "$satellites" >/dev/null
                    checkrc $? "mosquitto_pub $MQTT_HOST" "continue"
                else
                    # Sending straight to the wiotp cloud broker, as the gw 
                    echo mosquitto_pub -h "$MQTT_HOST" -p $MQTT_PORT -i "$CLIENT_ID" -u 'use-token-auth' -P "$WIOTP_GW_TOKEN" --cafile $WIOTP_PEM_FILE -q 1 -t "$TOPIC" -m "$satellites"
                    mosquitto_pub -h "$MQTT_HOST" -p $MQTT_PORT -i "$CLIENT_ID" -u 'use-token-auth' -P "$WIOTP_GW_TOKEN" --cafile $WIOTP_PEM_FILE -q 1 -t "$TOPIC" -m "$satellites" >/dev/null
                    checkrc $? "mosquitto_pub $MQTT_HOST" "continue"
                fi
                #wiotp_pub $MQTT_HOST $HZN_DEVICE_ID "iot-2/evt/status/fmt/json" $satellites 0

                # NOTE: Sending is best effort only, if there is a problem send will display the error, but then we continue
                if [[ -z "$error" ]]; then
                    # The send was successful, so reset the repeat variables
                    previous_sat="$sats"
                    num_repeat_sat_readings=0
                fi
            fi
        fi
    else
        # Log the gps REST service failure
        echo "ERROR: REST request to '$GPS_SATELLITES_URI' microservice failed: rc=$curlrc, httpcode=$httpcode, output="$satellites_rest
    fi

    # Pause
    # -----

    # Wait for the specified number of seconds before going again
    sleep $REPORTING_INTERVAL

done

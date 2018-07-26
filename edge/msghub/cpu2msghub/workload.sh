#!/bin/sh

# Horizon sample workload to query the cpu load from a sample service, calculate a window average, and publish it to Watson IoT Platform
# This workload expects the CPU service to be running, unless it is running in mock mode.

# Verify required environment variables are set
checkRequiredEnvVar() {
  varname=$1
  if [ -z $(eval echo \$$varname) ]; then
    echo "Error: Environment variable $varname must be set; exiting."
    exit 2
  else
    echo "  $varname="$(eval echo \$$varname)
  fi
}

# Environment variables that can optionally be set, or default
SAMPLE_INTERVAL="${SAMPLE_INTERVAL:-5}"    # how often (in seconds) to query the cpu (the gps location is queried every SAMPLE_INTERVAL * SAMPLE_SIZE seconds)
SAMPLE_SIZE="${SAMPLE_SIZE:-10}"    # the number of cpu samples to read before calculating and publishing the cpu average and gps coordinates
PUBLISH="${PUBLISH:-true}"    # whether or not to actually send data to IBM Message Hub
MOCK="${MOCK:-false}"     # if "true", just pretend to call the cpu service REST API
VERBOSE="${VERBOSE:-0}"    # set to 1 for verbose output

echo "Optional environment variables (or default values): SAMPLE_INTERVAL=$SAMPLE_INTERVAL, SAMPLE_SIZE=$SAMPLE_SIZE, PUBLISH=$PUBLISH, MOCK=$MOCK"

# When this workload is running in standalone mode, there are no required env vars.
if [[ "$PUBLISH" == "true" ]]; then
  echo "Checking for required environment variables for publishing to IBM Message Hub:"
  checkRequiredEnvVar "HZN_ORGANIZATION"      # automatically passed in by Horizon
  checkRequiredEnvVar "HZN_DEVICE_ID"      # automatically passed in by Horizon
  checkRequiredEnvVar "MSGHUB_API_KEY"
  checkRequiredEnvVar "MSGHUB_BROKER_URL"
  MSGHUB_USERNAME="${MSGHUB_API_KEY:0:16}"
  MSGHUB_PASSWORD="${MSGHUB_API_KEY:16}"
  MSGHUB_TOPIC="$HZN_ORGANIZATION.$HZN_DEVICE_ID"  #todo: see if HZN_PATTERN is passed in and use that
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

echo 'Starting infinite loop to read from the cpu and gps services then publish...'
sum=0
samplesrequired="$SAMPLE_SIZE"
samplecount=0
while true; do
  samplecount=$((samplecount + 1))

  # Get data from a local service
  if [[ "$MOCK" == "true" ]]; then
    output='{"cpu":'$(date +%S)'} 200'
    curlrc=0
  else
    output=$(curl -sS -w %{http_code} "http://cpu:8347/v1/cpu")
    curlrc=$?     # save this before it gets overwritten
  fi
  httpcode=${output:$((${#output}-3))}    # the last 3 chars are the http code
  json="${output%?[0-9][0-9][0-9]}"   # for the output, get all but the newline and 3 digits of http code

  if [[ "$curlrc" != 0 ]]; then
    echo "Warning: Curl command to the local cpu service returned exit code $curlrc, will try again next interval."
  elif [[ "$httpcode" != 200 ]]; then
    echo "Warning: HTTP code $httpcode from the local cpu service REST API, will try again next interval."
  else
    # Accumulate the CPU usage and calculate the average after obtaining all samples.
    cpuusage=$(echo $json | jq '.cpu')
    if [[ "$VERBOSE" == 1 ]]; then echo " Interval $samplecount cpu: $cpuusage"; fi
    sum=$(echo $sum + $cpuusage | bc)

    if [[ "$samplecount" -eq "$samplesrequired" ]]; then
      # Have enough samples, ready to publish
      average=$(echo "scale=4; $sum/$samplesrequired" | bc -l)

      # Also get gps coordinates from the GPS service
      if [[ "$MOCK" == "true" ]]; then
        output='{"lat":0.0,"lon":0.0,"alt":0.0} 200'
        curlrc=0
      else
        output=$(curl -sS -w %{http_code} "http://gps:31779/v1/gps/location")
        curlrc=$?     # save this before it gets overwritten
      fi
      httpcode=${output:$((${#output}-3))}    # the last 3 chars are the http code
      if [[ "$curlrc" != 0 || "$httpcode" != 200 ]]; then
        echo "Warning: the gps curl cmd failed with exit code $curlrc, HTTP code $httpcode. Will try again next interval"
        loc='{"lat":0.0,"lon":0.0,"alt":0.0}'
      else
        json="${output%?[0-9][0-9][0-9]}"   # for the output, get all but the newline and 3 digits of http code
        loc=$(echo $json | jq -c '{ lat: .latitude, lon: .longitude, alt: .elevation }')
      fi

      json='{"cpu":'$average',"gps":'$loc'}'
      #echo "avg: $json"

      if [[ "$PUBLISH" == "true" ]]; then
        # Send data to IBM Message Hub
        echo "$json | kafkacat -P -b $MSGHUB_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=$MSGHUB_USERNAME -X sasl.password=$MSGHUB_PASSWORD -t $MSGHUB_TOPIC
        echo "$json" | kafkacat -P -b $MSGHUB_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=$MSGHUB_USERNAME -X sasl.password=$MSGHUB_PASSWORD -t $MSGHUB_TOPIC
        checkrc $? "kafkacat" "continue"
      else
        echo "$json"
      fi
      sum=0
      samplecount=0
    fi

  fi

  # Pause before looping again
  sleep $SAMPLE_INTERVAL
done
# never reached

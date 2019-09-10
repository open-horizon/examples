#!/bin/sh

# Horizon sample service to query the cpu load from a sample service, calculate a window average, query the gps coordinate, and publish them to Watson IoT Platform
# This service requires the CPU and GPS services to be running, unless it is running in mock mode.

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

# Environment variables that can optionally be set, or default
SAMPLE_INTERVAL="${SAMPLE_INTERVAL:-5}"    # how often (in seconds) to query the cpu (the gps location is queried every SAMPLE_INTERVAL * SAMPLE_SIZE seconds)
SAMPLE_SIZE="${SAMPLE_SIZE:-10}"    # the number of cpu samples to read before calculating and publishing the cpu average and gps coordinates
PUBLISH="${PUBLISH:-true}"    # whether or not to actually send data to IBM Event Streams
MOCK="${MOCK:-false}"     # if "true", just pretend to call the cpu service REST API
VERBOSE="${VERBOSE:-0}"    # set to 1 for verbose output
CPU_URL="${CPU_URL:-http://ibm.cpu:80/v1/ibm.cpu}"
GPS_URL="${GPS_URL:-http://ibm.gps:80/v1/gps/location}"

echo "Optional environment variables (or default values): SAMPLE_INTERVAL=$SAMPLE_INTERVAL, SAMPLE_SIZE=$SAMPLE_SIZE, PUBLISH=$PUBLISH, MOCK=$MOCK"

# When this service is running in standalone mode, there are no required env vars.
if [[ "$PUBLISH" == "true" ]]; then
  echo "Checking for required environment variables for publishing to IBM Event Streams:"
  checkRequiredEnvVar "HZN_ORGANIZATION"      # automatically passed in by Horizon
  checkRequiredEnvVar "HZN_DEVICE_ID"      # automatically passed in by Horizon
  checkRequiredEnvVar "EVTSTREAMS_TOPIC"
  checkRequiredEnvVar "EVTSTREAMS_API_KEY"
  checkRequiredEnvVar "EVTSTREAMS_BROKER_URL"
  EVTSTREAMS_USERNAME="${EVTSTREAMS_API_KEY:0:16}"
  EVTSTREAMS_PASSWORD="${EVTSTREAMS_API_KEY:16}"
  # The only special chars allowed in the topic are: -._
  EVTSTREAMS_TOPIC="${EVTSTREAMS_TOPIC//[@#%()+=:,<>]/_}"
  # Tranlating the slashes does not work in the above bash substitute in alpine, so use tr
  EVTSTREAMS_TOPIC=$(echo "$EVTSTREAMS_TOPIC" | tr / _)
  echo "Will publish to topic: $EVTSTREAMS_TOPIC"
  if [[ -n "$EVTSTREAMS_CERT_ENCODED" && "$EVTSTREAMS_CERT_ENCODED" != "-" ]]; then
    # They are using an instance of Event Streams deployed in ICP, because it needs a self-signed cert
    EVTSTREAMS_CERT_FILE=/tmp/es-cert.pem
    echo "$EVTSTREAMS_CERT_ENCODED" | base64 -d > $EVTSTREAMS_CERT_FILE
    checkrc $? "decode cert"
  fi
fi

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
    if [[ "$VERBOSE" == 1 ]]; then echo " Calling REST API ${CPU_URL}..."; fi
    output=$(curl -sS -w %{http_code} --connect-timeout 5 "$CPU_URL")
    curlrc=$?     # save this before it gets overwritten
  fi
  httpcode=${output:$((${#output}-3))}    # the last 3 chars are the http code

  if [[ "$curlrc" != 0  || "$httpcode" != 200 || ${#output} -lt 5 ]]; then
    echo "Warning: the cpu curl cmd failed with exit code $curlrc, HTTP code $httpcode, or returned no data. Will try again next interval"
  else
    # Accumulate the CPU usage and calculate the average after obtaining all samples.
    json="${output%?[0-9][0-9][0-9]}"   # for the output, get all but the 3 digits of http code
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
        if [[ "$VERBOSE" == 1 ]]; then echo " Calling REST API ${GPS_URL}..."; fi
        output=$(curl -sS -w %{http_code} --connect-timeout 5 "$GPS_URL")
        curlrc=$?     # save this before it gets overwritten
      fi
      httpcode=${output:$((${#output}-3))}    # the last 3 chars are the http code
      if [[ "$curlrc" != 0 || "$httpcode" != 200 || ${#output} -lt 5 ]]; then
        echo "Warning: the gps curl cmd failed with exit code $curlrc, HTTP code $httpcode, or returned no data. Will try again next interval"
        loc='{"lat":0.0,"lon":0.0,"alt":0.0}'
      else
        json="${output%?[0-9][0-9][0-9]}"   # for the output, get all but the 3 digits of http code
        loc=$(echo $json | jq -c '{ lat: .latitude, lon: .longitude, alt: .elevation }')
      fi

      json='{"nodeID":"'$HZN_DEVICE_ID'","cpu":'$average',"gps":'$loc'}'
      #echo "avg: $json"

      if [[ "$PUBLISH" == "true" ]]; then
        if [[ -n "$EVTSTREAMS_CERT_FILE" ]]; then
          # Send data to ICP Event Streams
          echo "echo $json | kafkacat -P -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$EVTSTREAMS_API_KEY -X ssl.ca.location=$EVTSTREAMS_CERT_FILE -t $EVTSTREAMS_TOPIC"
          echo "$json" | kafkacat -P -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$EVTSTREAMS_API_KEY -X ssl.ca.location=$EVTSTREAMS_CERT_FILE -t $EVTSTREAMS_TOPIC
          checkrc $? "kafkacat" "continue"
        else
          # Send data to IBM Cloud Event Streams
          echo "echo $json | kafkacat -P -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=$EVTSTREAMS_USERNAME -X sasl.password=$EVTSTREAMS_PASSWORD -t $EVTSTREAMS_TOPIC"
          echo "$json" | kafkacat -P -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=$EVTSTREAMS_USERNAME -X sasl.password=$EVTSTREAMS_PASSWORD -t $EVTSTREAMS_TOPIC
          checkrc $? "kafkacat" "continue"
        fi
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

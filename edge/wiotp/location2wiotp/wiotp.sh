#
# Central Watson IoT Platform MQTT service support routines.
#
# Not currently being used
#
# Written by Glen Darling
#

# NOTE: This library requires the following globals to be set:
#           WIOTP_DOMAIN - domain base for WIoTP MQTT URL
#           HZN_ORG_ID - WIoTP organization ID
#           WIOTP_GW_TYPE - WIoTP gateway type
#           WIOTP_GW_ID - WIoTP gateway instance ID
#           WIOTP_GW_TOKEN - WIoTP gateway auth token
#           WIOTP_API_KEY - WIoTP API key name
#           WIOTP_API_TOKEN - WIoTP API key auth token
#           HZN_DEVICE_ID - Horizon device ID (e.g., "g@gwtype@gwid")

# Diagnostic output
DEBUG=true
DEBUG="${DEBUG:-0}"      # set to 1 for debug output
VERBOSE="${VERBOSE:-0}"  # set to 1 for verbose output
OUTPUT_PREFIX="WIOTP "   # prefix for all normal output echoed by this module

# Simple debug utility
debug() {
    local message=$1
    if [[ -n "$DEBUG" ]]
    then
        echo "DEBUG: $message"
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

#
# WIoTP Public functions
#

# Function "wiotp_pub"
#     Publishes an MQTT message to WIoTP using the passed arguments
# Arguments:
#     $1: target host
#     $2: client ID
#     $3: topic
#     $4: message to be sent
#     $5: qos (quality of service) setting for this publication:
#           0: deliver message once, no confirmations
#           1: deliver message at least once, with confirmation required
#           2: deliver message exactly once, using four step handshake
#     $6: (optional) any other args you want to pass to mosquito_pub
# Side Effects:
#     in:   uses global variable WIOTP_PEM_FILE
#     out:  on return this sets global $error
wiotp_pub() {

    local host=$1
    local client=$2
    local topic=$3
    local json=$4
    local qos=$5
    local otheropts=$6

    debug "--> wiotp_pub(host=$host, client=$client, topic=$topic, msg=$json qos=$qos)"
    if [[ "$VERBOSE" == 1 ]]; then echo mosquitto_pub -h "$host" -p 8883 -i "$client" --cafile $WIOTP_PEM_FILE -q $qos -t "$topic" -m "$json" $otheropts --quiet 2>&1; fi

    # Actually publish the message to WIoTP using the passed arguments
    error=`mosquitto_pub -h "$host" -p 8883 -i "$client" --cafile $WIOTP_PEM_FILE -q $qos -t "$topic" -m "$json" $otheropts --quiet 2>&1`
    if [[ -n "$error" ]]; then
        debug "<-- wiotp_pub: error='$error'"
    fi
    # Any $error that occurs is just quietly propagated back to the caller
}

# NOTE: Appending this "$@" line enables any function in this module
#       to be invoked from the shell simply by passing it (followed
#       by any needed arguments) as arguments to this file.  Note also
#       that when no arguments are passed to this module when it is
#       sourced, then the line below has no effect,
#       E.g., this will invoke listen with its 3 arguments:
#           $ ./mqtt.sh listen username password topic
$@

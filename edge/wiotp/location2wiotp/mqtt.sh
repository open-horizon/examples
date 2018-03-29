#
# Central Blue Horizon MQTT service support routines.
#
# Written by Glen Darling
#

# NOTE: This library requires the following globals to be set:
#           HZN_BROKER -- the central Blue Horizon MQTT service URI
#           MQTT_PORT -- the central Blue Horizon MQTT service URI
#           HZN_AGREEMENTID -- the Blue Horizon contract ID (provided by Anax)
#           HZN_HASH -- the Blue Horizon nonce (provided by Anax)
#           CA_FILE -- the Blue Horizon cert file

#
# Configuration constants
#

# Uncomment this line (or set DEBUG in the environment) to turn on debug output
#DEBUG=true

# This prefix is used for all normal output echoed by this module
# OUTPUT_PREFIX="MQTT"

# The (central) Blue Horizon MQTT Broker (Verne) URI and port number
MQTT_BROKER_URI="$HZN_BROKER"

# The registration message consists of the Blue Horizon contract ID and
# the Blue Horizon configuration nonce (both of which are expected in the
# process environment when launched (because they are set by Anax).
REGISTRATION_MESSAGE="$HZN_AGREEMENTID $HZN_HASH"

# Global state of last message publication attempt
last_send_failed=false

# Simple debug utility
debug() {
    local message=$1
    if [[ -n "$DEBUG" ]]
    then 
        echo "DEBUG: $message"
    fi
}

#
# MQTT Public functions
#

# Function "mqtt_pub"
#     Publishes an MQTT message using the passed arguments
# Arguments:
#     $1: user name
#     $2: password
#     $3: destination topic path for the message
#     $4: message to be sent
#     $5: "quality of service" setting for this publication:
#        0: The broker/client will deliver the message once, with no confirmation.
#        1: The broker/client will deliver the message at least once, with confirmation required.
#        2: The broker/client will deliver the message exactly once by using a four step handshake.
# Side Effects:
#     sets global $error
mqtt_pub() {
    local username=$1
    local password=$2
    local topic=$3
    local message=$4
    local qos=$5
    local otheropts=$6
    debug "--> mqtt_pub( topic=$topic msg=$message qos=$qos )"
    # Actually publish the message to MQTT per the passed arguments
    error=`mosquitto_pub -h $MQTT_BROKER_URI -p $MQTT_PORT -t "$topic" -m "$message" -u "$username" -P "$password" --cafile $CA_FILE -q $qos $otheropts --quiet 2>&1`
    if [[ -n "$error" ]]; then
        debug "<-- mqtt_pub: error='$error'"
    fi
    # Any $error that occurs is just quietly propagated back to the caller
}

# Function "mqtt_sub"
#     Subscribes to an MQTT message using the passed arguments. Note: this is only used in listen(), which is only used in test scripts.
#     Note that this method subscribes using QoS=0 (fire-and-forget, or receive messages at most once).
# Arguments:
#     $1: user name
#     $2: password
#     $3: destination topic path for the message
# Side Effects:
#     sets global $error
mqtt_sub() {
    local username=$1
    local password=$2
    local topic=$3
    debug "--> mqtt_sub( topic=$topic )"
    # Actually subscribe to the MQTT topic per the passed arguments
    /usr/bin/mosquitto_sub -h $MQTT_BROKER_URI -p $MQTT_PORT -t "$topic" -u "$username" -P "$password" --cafile $CA_FILE -q 0
    # NOT REACHED (mosquitto_sub is expected to never exit)
    error="/usr/bin/mosquitto_sub unexpectedly exited: $?."
    echo "mqtt_sub: $error"
}

# Function "register"
#     Attempts to register with the central Blue Horizon MQTT service,
#     retrying until either successful or too many retries have occurred.
# Arguments:
#     none
# Side Effects:
#     sets global $error if retrying fails
register() {
    echo "Registering with the Blue Horizon MQTT service..."
    # debug "--> register()"
    # At most, try registration this many times
    local tries_remaining=$MAX_REGISTRATION_ATTEMPTS
    while [[ $tries_remaining -gt 0 ]]
    do
        debug "register: tries_remaining=$tries_remaining"

        # Decrement the counter
        # tries_remaining=`expr $tries_remaining - 1`
        tries_remaining=$(($tries_remaining-1))

        # Send the registration request to MQTT using MQTT QoS of 2 to only give it one registration request
        mqtt_pub "public" "public" "/registration" "$REGISTRATION_MESSAGE" 2

        # Did sending the MQTT registration request succeed?
        if [[ -z "$error" ]]; then
            # Succeeded
            # local tries_count=$(expr $MAX_REGISTRATION_ATTEMPTS - $tries_remaining)
            local tries_count=$(($MAX_REGISTRATION_ATTEMPTS - $tries_remaining))

            echo "register: MQTT registration request successfully sent (on try $tries_count). Sleeping for $REG_SECONDS_BEFORE_STREAMING seconds to allow registration to take effect..."
            # Pause before streaming begins
            sleep $REG_SECONDS_BEFORE_STREAMING
            break
        # Sending the request failed.  Should we try again (or have we tried too often)?
        elif [[ $tries_remaining -le 0 ]]; then
            # Number of registration retries was exceeded. Loop will exit.
            echo "ERROR: register: failed to send MQTT registration request after $MAX_REGISTRATION_ATTEMPTS attempts: $error. Giving up until next data send interval."
            break
        else
            # Still have tries remaining.  Pause before trying again...
            sleep $SECONDS_BETWEEN_REG_ATTEMPTS
        fi
    done
    if [[ -n "$error" ]]; then
        debug "<-- register: error='$error'"
    fi
    # Note: The only error exit case is in the elif above, so the $error is shown in the message echoed there.
}

# Function "send"
#     Sends a data update to the central Blue Horizon MQTT service.
#     Note that this is best effort.  If the send fails due to a registration
#     error, then an attempt is made to re-register.  Whether that attempt
#     succeeds or fails, the result is simply propagated back to the caller.
#     Under no circumstance is any attempt made to retry the send.  If the
#     caller wishes to retry they should do so slowly to avoid flooding
#     the central MQTT service.
# Arguments:
#     $1: username -- Blue Horizon contract ID (provided by Anax)
#     $2: password -- Blue Horizon contract configure nonce (provided by Anax)
#     $3: topic -- likely $SAT_TOPIC (satellites) or $LOC_TOPIC (location)
#     $4: message to be sent
# Side Effects:
#     sets global $error
send() {
    local username=$1
    local password=$2
    local topic=$3
    local message=$4
    local qos=$5
    local otheropts=$6
    # Publish the data to MQTT using MQTT QoS of 1 (deliver at least once, since the location info is critical for display on the map)
    echo "send: $topic $message $qos $otheropts"
    mqtt_pub "$username" "$password" "$topic" "$message" "$qos" "$otheropts"
    # Did the MQTT publication fail?
    if [[ -n "$error" ]]; then
        # Note that this publication failed
        debug "send: failed: $error"
        last_send_failed=true
        # Was the failure due to lack of registration?
        echo "$error" | grep -q -E "not authori"
        if [[ $? -eq 0 ]]; then
            echo "send: failed due to lack of registration: $error. Will send registration request again..."
            # Attempt re-registration (with retries)
            register
            # Afterward, whether successful or not, simply fall through
            # passing $error back to the caller.  Do not re-attempt the send.
            # As noted above, if the caller wishes to retry this send they
            # should do so slowly to avoid flooding the central MQTT service.
        fi
    else
        # Announce first successful publication after any failure
        if [[ $last_send_failed == true ]]
        then 
            echo "send: succeeded (after previous failure)"
        else 
            debug "send: succeeded"
        fi
        # Note that this publication was successful
        last_send_failed=false
    fi
    # debug "<-- send: error='$error'"
    # Any remaining $error is announced and then propagated back to the caller
    if [[ -n "$error" ]]; then
        echo "ERROR: send: $error"
    fi
}

# Function "listen"
#     Listens forever to an MQTT topic and echoes what it receives. Note: this is only used in test scripts.
# Arguments:
#     $1: username -- Blue Horizon contract ID (provided by Anax)
#     $2: password -- Blue Horizon contract configure nonce (provided by Anax)
#     $3: topic -- likely $SAT_TOPIC (satellites) or $LOC_TOPIC (location)
# Side Effects:
#     sets global $error
listen() {
    local username=$1
    local password=$2
    local topic=$3
    debug "--> listen( username=$username password=$password topic=$topic )"
    # Simply subscribe to the specified MQTT topic
    mqtt_sub "$username" "$password" "$topic"
    # NOT REACHED (unless an error occurs)
    echo "listen: $error"
}

# NOTE: Appending this "$@" line enables any function in this module        
#       to be invoked from the shell simply by passing it (followed          
#       by any needed arguments) as arguments to this file.  Note also          
#       that when no arguments are passed to this module when it is             
#       sourced, then the line below has no effect,                             
#       E.g., this will invoke listen with its 3 arguments:
#           $ ./mqtt.sh listen username password topic
$@                                               

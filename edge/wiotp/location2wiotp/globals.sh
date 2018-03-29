#
# Blue Horizon global variable settings.
# NOTE: this is only use when sending to bluehorizon (not wiotp), and will soon go away
#
# Original written by Dima Rekesh
# Modified by Glen Darling
#

# Configure the central Blue Horizon MQTT broker and cert details
# DEPL_ENV can be set to staging or prod. If its not set, use HZN_EXCHANGE_URL to determine a default
if [[ -n "$DEPL_ENV" ]]; then
    deplenv=$DEPL_ENV
elif [[ -n "$HZN_EXCHANGE_URL" ]]; then
    echo "$HZN_EXCHANGE_URL" | grep -q -E "staging"
    if [[ $? -eq 0 ]]; then
        deplenv="staging"
    else
        deplenv="prod"
    fi
fi
if [[ "$deplenv" == "staging" ]]; then
    if [[ "$USE_NEW_STAGING_URL" == "false" ]]; then
        # verne has NOT been switched over to the new staging hostname yet
        export HZN_BROKER=staging.bluehorizon.hovitos.engineering
        export CA_FILE=ca-staging.pem
    else
        export HZN_BROKER=staging.bluehorizon.network
        export CA_FILE=ca-new-staging.pem
    fi
elif [[ "$deplenv" == "prod" ]]; then
    export HZN_BROKER=bluehorizon.network
    export CA_FILE=ca-prod.pem
else
    # we must be running in a dev environment
    export HZN_BROKER=haproxy.rekesh.com
    export CA_FILE=ca-test.pem
fi

# HZN_DEVICE_ID should always come from the environment
# Configure the Blue Horizon device ID (if not already set in the environment)
#if [[ -z "$HZN_DEVICE_ID" ]]
#then
#  export HZN_DEVICE_ID=`cat /proc/cpuinfo | /bin/grep Serial | /usr/bin/awk '{print $3;}'`
#fi

# MQTT registration fails some times, so our functions will retry using these configuration values before announcing
# or propagating errors.
MAX_REGISTRATION_ATTEMPTS=${MAX_REGISTRATION_ATTEMPTS:-3}
SECONDS_BETWEEN_REG_ATTEMPTS=${SECONDS_BETWEEN_REG_ATTEMPTS:-3}
# Amount of time to pause after registration before starting to stream data
REG_SECONDS_BEFORE_STREAMING=${REG_SECONDS_BEFORE_STREAMING:-20}

echo "Blue Horizon location workload configuration:"
echo "    DEPL_ENV=$DEPL_ENV"
echo "    HZN_EXCHANGE_URL=$HZN_EXCHANGE_URL"
echo "    HZN_BROKER=$HZN_BROKER"
echo "    CA_FILE=$CA_FILE"
echo "    MAX_REGISTRATION_ATTEMPTS=$MAX_REGISTRATION_ATTEMPTS"
echo "    SECONDS_BETWEEN_REG_ATTEMPTS=$SECONDS_BETWEEN_REG_ATTEMPTS"
echo "    REG_SECONDS_BEFORE_STREAMING=$REG_SECONDS_BEFORE_STREAMING"

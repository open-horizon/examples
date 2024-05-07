#!/bin/sh

# Check env vars that we know should be set to verify that everything is working
function verify {
    if [ "$2" == "" ]
    then
        echo -e "Error: $1 should be set but is not."
        exit 2
    fi
}

verify "HZN_ESS_AUTH" $HZN_ESS_AUTH
verify "HZN_ESS_CERT" $HZN_ESS_CERT
verify "HZN_ESS_API_ADDRESS" $HZN_ESS_API_ADDRESS
verify "HZN_ESS_API_PORT" $HZN_ESS_API_PORT
verify "HZN_ESS_API_PROTOCOL" $HZN_ESS_API_PROTOCOL
verify "MMS_OBJECT_TYPES" $MMS_OBJECT_TYPES
echo -e "All ESS env vars verified."

/usr/local/bin/service -v 3 -logtostderr
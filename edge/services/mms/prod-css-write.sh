#!/bin/bash

# Create a new "objectID" (of a specified "objectType" at a specified "version)
# in the production (cloud) CSS.

# User must have `HZN_CSS_URL` set in their environment
if [ -z "$HZN_CSS_URL" ]; then
    echo "ERROR: \"HZN_CSS_URL\" must contain the cloud CSS URL."
    exit 1
fi

# Local CSS REST API endpoint (i.e., the base URL) -- NOTE: use `http`!
CSS_ENDPOINT="${HZN_CSS_URL}/api/v1"

# User must have `HZN_ORG_ID` and `HZN_EXCHANGE_USER_AUTH` set too
if [ -z "$HZN_ORG_ID" ]; then
    echo "ERROR: \"HZN_ORG_ID\" must contain your organization name."
    exit 1
fi
if [ -z "$HZN_EXCHANGE_USER_AUTH" ]; then
    echo "ERROR: \"HZN_EXCHANGE_USER_AUTH\" must contain your credentials."
    exit 1
fi

# Construct the auth string
AUTH="${HZN_ORG_ID}/${HZN_EXCHANGE_USER_AUTH}"

# Metadata for the object being created
objectType="$1"
objectID="$2"
version="1.0.0"
description=""

# Which Edge Nodes to target ("" is the wildcard for all IDs/types)
# Normally you will set the `destinationType` to the pattern name you
# used to register your nodes (without the organization prefix). That
# way the object will only be sent to those Edge Nodes that are
# registered with that specific pattern.
destinationID=""
destinationType=""

# Construct object metadata from the above variables
METADATA='{"meta":{"objectID":"'${objectID}'","objectType":"'${objectType}'","destinationID":"'${destinationID}'","destinationType":"'${destinationType}'","version":"'${version}'","description":"'${description}'"}}'

# Writing to the CSS is a 2-step process:
#  1. Setup the metadata, including the version
#  2. Sending the data of the object

# STEP 1: send the metadata to the CSS API to create the new object
HTTP_CODE=`echo "$METADATA" | \
  curl -sSL -X PUT -u "${AUTH}" -w '%{http_code}' \
    --header 'Content-Type:application/octet-stream' \
    --data-binary @- \
    "${CSS_ENDPOINT}/objects/${HZN_ORG_ID}/${objectType}/${objectID}"`
if [ ${HTTP_CODE:-0} -ne 204 ]; then
  echo "ERROR! Object metadata config failed: ${HTTP_CODE}" &> /dev/stderr
  exit 1
fi
echo "Object metadata configured."

# STEP 2: use `cat` to read data from stdin, then send it to the local CSS API
HTTP_CODE=`cat | \
  curl -sSL -X PUT -u "${AUTH}" -w '%{http_code}' \
    --header 'Content-Type:application/octet-stream' \
    --data-binary @- \
    "${CSS_ENDPOINT}/objects/${HZN_ORG_ID}/${objectType}/${objectID}/data"`
if [ ${HTTP_CODE:-0} -ne 204 ]; then
  echo "ERROR! Object data write failed: ${HTTP_CODE}" &> /dev/stderr
  exit 1
fi
echo "Object data written."


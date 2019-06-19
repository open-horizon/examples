#!/bin/bash

# Create a new "objectID" (of a specified "objectType" at a specified "version)
# in the development CSS that was started by `hzn dev service start`.

# Local CSS REST API endpoint (i.e., the base URL) -- NOTE: use `http`!
CSS_ENDPOINT="http://localhost:8580/api/v1"

# User must have `HZN_ORG_ID` set in their environment
if [ -z "$HZN_ORG_ID" ]; then
    echo "ERROR: \"HZN_ORG_ID\" must contain your organization name."
    exit 1
fi

# Construct the auth string
# Note that these well-known credentials are for development use only,
# and they only work on the dev CSS started by `hzn dev service start`
AUTH="${HZN_ORG_ID}/${HZN_ORG_ID}admin:${HZN_ORG_ID}adminpw"

# Metadata for the object being created
objectType="$1"
objectID="$2"
version="1.0.0"
description=""

# Which Edge Nodes to target ("" is the wildcard for all IDs/types)
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
    --trace-ascii dump.txt \
    --header 'Content-Type:application/octet-stream' \
    --data-binary @- \
    "${CSS_ENDPOINT}/objects/${HZN_ORG_ID}/${objectType}/${objectID}"`
if [ ${HTTP_CODE:-0} -ne 204 ]; then
  echo "ERROR! Object metadata config failed: ${HTTP_CODE}" &> /dev/stderr
  exit 1
fi
echo "Object metadata configured."

# STEP 2: use `cat` to read data from stdin, then send it to the dev CSS API
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


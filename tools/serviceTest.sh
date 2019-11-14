#!/bin/bash

########################### TESTING ############################
# For Linux and Mac:
# hzn service log -f $SERVICE_NAME
################################################################

# $1 == $SERVICE_NAME
# $2 == "key" being searched for to know if service is successfully runnning. If found, exit(0)
# $3 == timeout - if exceeded, service failed. exit(1)

name=$1
match=$2
timeOut=$3
START=$SECONDS

command="hzn service log -f $name"

####################### Loop until until either MATCH is found or TIMEOUT is exceeded #####################
$command | while read line; do
    # MATCH was found
    if grep -q -m 1 "$match" <<< "$line"; then
        exit 0

    # TIMEOUT was exceeded
    elif [ "$(($SECONDS - $START))" -ge "$timeOut" ]; then
        exit 1
    fi

    sleep 1;

done

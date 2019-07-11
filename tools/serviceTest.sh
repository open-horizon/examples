#!/bin/bash

########################### TESTING ############################
# For now on Mac:
# docker logs -f `docker ps -q --filter name=$SERVICE_NAME`
#
# For now on Linux:
# tail -f /var/log/syslog
################################################################

# $1 == $SERVICE_NAME
# $2 == "key" to know if service is successfully runnning. If found, exit(0)
# $3 == timeout - if exceeded, service failed. exit(1)

name=$1
match=$2
timeOut=$3
START=$SECONDS

##################################### Check the operating system #########################################
if [ $(uname -s) == "Darwin" ]; then
    # This is a MAC machine
    docker logs -f `docker ps -q --filter name=$name` > output.txt &
else
    # This is a LINUX machine
    tail -f /var/log/syslog > output.txt &
fi

####################### Loop until until either MATCH is found or TIMEOUT is exceeded #####################
while :
do
    curTime=$SECONDS

    # MATCH was found
    if grep "$match" output.txt; then
        rm output.txt
        exit 0

    # TIMEOUT was exceeded
    elif [ "$(($SECONDS - $START))" -ge "$timeOut" ]; then
        rm output.txt
        exit 1
    fi

    sleep 1;

done

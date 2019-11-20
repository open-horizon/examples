#!/bin/bash

########################### TESTING ############################
# For now on Mac:
# docker logs -f `docker ps -q --filter name=$SERVICE_NAME`
#
# For now on Linux:
# sudo tail -f /var/log/syslog
################################################################

# $1 == $SERVICE_NAME
# $2 == "key" being searched for to know if service is successfully runnning. If found, exit(0)
# $3 == timeout - if exceeded, service failed. exit(1)

name=$1
match=$2
timeOut=$3
START=$SECONDS

##################################### Check for hzn service log command and OS ###########################
logCommand=$( hzn service log 2>&1 )
if [ "$logCommand" == 'hzn: error: expected command but got "log", try --help' ]; then
  # Check the operating system
  if [ $(uname -s) == "Darwin" ]; then
      # This is a MAC machine
      command="docker logs -f `docker ps -q --filter name=$name`"
  else
      # This is a LINUX machine
      command="sudo tail -f /var/log/syslog"
  fi
else
  command="hzn service log -f $name"
fi

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

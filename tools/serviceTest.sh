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

##################################### Check the operating system #########################################
if [ $(uname -s) == "Darwin" ]; then
    # This is a MAC machine
    command="docker logs -f `docker ps -q --filter name=$name`"
else
    # This is a LINUX machine
    command="sudo grep -q -m 1 "$match" /var/log/syslog"
fi
docker ps -a

declare -i counter
counter=0
while :
do
  #sudo tail -n50 /var/log/syslog
  $command
  if [ $? -eq 0 ]; then
    echo "found it"
    exit 0
  elif [ $counter -ge $timeOut ]; then
    echo "timeout"
    exit 1
  fi

  counter=$counter+1
  #sudo tail -n1 /var/log/syslog
  sleep 1
done
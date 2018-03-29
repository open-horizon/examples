#!/bin/sh
#
# Blue Horizon location workload test listener
#
# Usage:
#
#     $ listener.sh [ interval [ quiet ] ]
#
# Where:
#
#     interval(integer)  report stats after this many reads
#     quiet(any)         suppresses the echoing of every line read
#
# Examples:
#
#     Run with full output of all messages read with no statistics reports
#         listener.sh
#
#     Run with full output of all messages read and report every 1000 messages
#         listener.sh 1000
#
#     Run quietly (no output of messages read) and report every 100000 messages:
#         listener.sh 100000 q
#
# Written by Glen Darling
#

# Configuration
interval=0
if [[ -n "$1" ]]
then
    interval=$1
fi
quiet=0
if [[ -n "$2" ]]
then
    quiet=1
fi
STATS_PREFIX="Location test: "

echo `/bin/date` "Listening to the location data stream..."
STARTTIME=$(/bin/date +%s)
last=$STARTTIME
count=0
mean=0
min=1000000
max=0
# Pipe the MQTT subscription to the "while read" loop below...
# Note that these HZN_ config variables are expanded in the subshell, not here.
/bin/sh -c '. /globals.sh > /dev/null; /mqtt.sh listen $HZN_AGREEMENTID $HZN_HASH $LOC_TOPIC' |
    while read line
    do
        # Echo the line (unless quiet was requested) and keep a running count
        if [[ $quiet -eq 0 ]]
        then
            echo $line
        fi
        count=$(expr $count + 1)

        # Compute delay since last message, and keep running min, mean, and max
        now=$(/bin/date +%s)
        delay=$(expr $now - $last)
        last=$now
        # Ignore meaningless first delay (i.e., only compute when count > 1)
        if [[ $count -gt 1 ]]
        then
            if [[ $delay -lt $min ]]
            then
                min=$delay
            fi
            if [[ $delay -gt $max ]]
            then
                max=$delay
            fi

            # Compute running arithmetic mean of delay times (except first one)
            mean=$(expr $(expr $(expr $mean \* $(expr $count - 2)) + $delay) / $(expr $count - 1))
        fi

        # Emit stats periodically, unless disabled
        if [[ $interval -ne 0 ]]
        then
            # Compute count modulo interval
            mod=`echo "$count % $interval" | /usr/bin/bc`
            if [[ $mod -eq 0 ]]
            then
                elapsed=$(expr $now - $STARTTIME)
                printf "%sUp %ds: %d msgs. Delays: min=%ds, mean=%ds, max=%ds.\n" "$STATS_PREFIX" $elapsed $count $min $mean $max
            fi
        fi

        # Watch for the error prefix, and exit when one is read
        echo "$line" | grep -q -E "^MQTT ERROR"
        if [[ $? -eq 0 ]]
        then
            break
        fi

    done


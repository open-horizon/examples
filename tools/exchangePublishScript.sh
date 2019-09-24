#!/bin/sh

# if the org id is set locally we don't want to override the IBM org of these samples
unset HZN_ORG_ID

# git branch/repository to clone
branch="-b master"
repository="https://github.com/open-horizon/examples.git"

# text file containing servies and patterns to publish
input="$(dirname $0)/examples/tools/blessedSamples.txt"

topDir=$(pwd)
error=0

git clone $branch $repository

# Check if EVTSTREAM_* env vars are empty, and give default values if so
if [ -z $EVTSTREAMS_ ]; then
    echo ""
    echo "EVTSTREAM_* variables for IBM Event Streams are not set. Providing default values."
    echo ""
    EVTSTREAMS_API_KEY="Some default value"
    echo $EVTSTREAMS_API_KEY
    echo ""
fi

# read in blessedSamples.txt which contains the services and patterns to publish
while IFS= read -r line
do
    # each $line contains the path to any service or pattern that needs to be published
    if cd $line; then
        echo `pwd`
        if make publish-only; then
            echo ""
        else
            echo "\n*** Error publishing $line *** \n"
            error=1
        fi
        cd $topDir

    else
        echo "\n*** Error finding service directory $line *** \n" 1>&2
        error=1
    fi

done < "$input"


# clean up if no errors
if [ $error != 0 ]; then
    echo "\n*** Errors were encountered when publishing, examples repo was not deleted *** \n"
else
    rm -f -r examples/
fi

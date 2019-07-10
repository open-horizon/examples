#!/bin/sh

# if the org id is set locally we don't want to override the IBM org of these samples
unset HZN_ORG_ID

# git repository to clone
repository="https://github.com/open-horizon/examples.git"

# text file containing servies and patterns to publish
curl https://raw.githubusercontent.com/open-horizon/examples/master/tools/blessedSamples.txt -O

# text file containing servies and patterns to publish
input="$(dirname $0)/blessedSamples.txt"

topDir=$(pwd)
error=0

git clone "$repository"

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
    rm blessedSamples.txt
fi

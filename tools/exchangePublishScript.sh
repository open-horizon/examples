#!/bin/sh

# git repository to clone
repository="https://github.com/open-horizon/examples.git"

# text file containing servies and patterns to publish
input="$(dirname $0)/blessedSamples.txt"

topDir=$(pwd)
echo $topDir

git clone "$repository"

# read in blessedSamples.txt which contains the services and patterns to publish
while IFS= read -r line
do
    # each $line contains the path to any service or pattern that needs to be published
    if cd $line; then
        echo `pwd`
        if echo `make publish-only`; then
            echo ""
        else
            echo "Error publishing."
        fi
        cd $topDir

    else
        echo "Error finding service dirsctory $line" 1>&2

    fi

done < "$input"


#!/bin/sh

# git repository to clone
repository="https://github.com/open-horizon/examples.git"

# text file containing servies and patterns to publish
input="samples.txt"

git clone "$repository"

# if x == 0 we are still working with lower level services
x=0

# read in samples.txt which contains the services and patterns to publish
while IFS= read -r line
do
    # when an empty line is encountered set x = 1 then begin publishing patterns
    if [ "$line" == "" ]; then
        x=1

    else
        # publishing lower level services
        if [ "$x" == "0" ]; then
            #if [ "$?" = "0" ]; then # check if the directory exists
            if cd examples/edge/services/$line; then
                echo `pwd`
                if echo `make publish-only`; then
                    echo ""
                else
                    echo "Error publishing."
                fi
                cd ../../../../

            else
                echo "Error finding service dirsctory $line" 1>&2
                #echo `pwd`

            fi

        else
            # publishing patterns
            if cd examples/edge/msghub/$line; then
                echo `pwd`
                if echo `make publish-only`; then
                    echo ""
                else
                    echo "Error publishing."
                fi
                cd ../../../../

            else
                echo "Error finding service dirsctory $line" 1>&2
                #echo `pwd`

            fi
        fi

    fi
    #echo "$line"

done < "$input"

echo `pwd`


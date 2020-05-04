#!/bin/sh

# if the org id is set locally we don't want to override the IBM org of these samples
unset HZN_ORG_ID

function scriptUsage () {
    cat << EOF

Usage: ./exchangePublishScript.sh [-c <cluster-name>]

Parameters:
  optional:
    -c              <cluster-name> set this flag to publish example deployment policy for the helloworld 
                      and cpu2evtstreams samples

Required Environment Variables:
  EXCHANGE_ROOT_PASS
  HZN_EXCHANGE_URL
  HZN_EXCHANGE_USER_AUTH

EOF
    exit 1
}

# parse any arguments
while (( "$#" )); do
    case "$1" in
        -h) # display script usage
            scriptUsage
            shift
            ;;
        -c) # cluster name to publish deployment policy
            ORG=$2
            shift 2
            ;;
    esac
done

# check if required environment variables are set
: ${HZN_EXCHANGE_URL:?} ${HZN_EXCHANGE_USER_AUTH:?}

# check the previous cmds exit code. 
checkexitcode () {   
    if [[ $1 == 0 ]]; then return; fi
    echo ""
    echo "Error: exit code $1 when $2"
    echo ""
    error=1
}

# publish deployment policy for helloworld and cpu2evtstreams if -c flag is used
function deployPolPublish () {
    if ([[ $line == *"cpu2evtstreams" ]] || [[ $line == *"helloworld" ]]); then 
        HZN_ORG_ID=$ORG make publish-business-policy
        checkexitcode $? "publishing deployment policy to the "$ORG" in the exchange"
    fi
}

# git branch/repository to clone
branch="-b master"
repository="https://github.com/open-horizon/examples.git"

# text file containing servies and patterns to publish
input="/tmp/examples/tools/blessedSamples.txt"

topDir=$(pwd)
error=0

git clone $branch $repository /tmp/examples

# read in blessedSamples.txt which contains the services and patterns to publish
while IFS= read -r line
do
    # each $line contains the path to any service or pattern that needs to be published
    cd /tmp/$line
    checkexitcode $? "finding service directory "$line""
    
    echo `pwd`
    make publish-only
    checkexitcode $? "publishing "$line" to the exchange"

    # check if an org was specified to publish sample deployment policy 
    if ! [ -z $ORG ]; then
        deployPolPublish
    fi

    cd $topDir

done < "$input"


# clean up if no errors
if [ $error != 0 ]; then
    echo -e "\n*** Errors were encountered when publishing, the cloned examples directory was not deleted *** \n"
else
    echo -e "\nSuccessfully published all content to the exchange. Removing examples directory...\n"
    rm -f -r /tmp/examples/
fi


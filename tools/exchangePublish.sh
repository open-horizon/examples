#!/bin/bash

# if the org id is set locally we don't want to override the IBM org of these samples
unset HZN_ORG_ID

scriptUsage() {
    cat << EOF
Usage: ${0##*/} [-h] [-v] [-c <cluster-name>]

Flag:
  -c <cluster-name>  Set this flag to publish example deployment policies to this org.
  -v                 Verbose output
  -h                 This usage

Required Environment Variables:
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
        -v) # verbose output
            VERBOSE=1
            shift
            ;;
        -c) # cluster name to publish deployment policy to
            ORG=$2
            shift 2
            ;;
    esac
done

# check if required environment variables are set
: ${HZN_EXCHANGE_URL:?} ${HZN_EXCHANGE_USER_AUTH:?}

# path to the cloned exmaples repo
PATH_TO_EXAMPLES=/tmp/open-horizon

# check the previous cmds exit code. 
checkexitcode() {
    local exitCode=$1
    local task=$2
    local dontExit=$3   # set to 'continue' to not exit for this error
    if [[ $exitCode == 0 ]]; then return; fi
    echo "Error: exit code $exitCode from: $task"
    if [[ $dontExit != 'continue' ]]; then
        exit $exitCode
    fi
}

# publish deployment policy for helloworld and cpu2evtstreams if -c flag is used
deployPolPublish() {
    if ([[ $line == *"cpu2evtstreams" ]] || [[ $line == *"helloworld" ]] || [[ $line == *"operator"* ]]); then 
        if [[ -n $VERBOSE ]]; then
            HZN_ORG_ID=$ORG make publish-deployment-policy
            checkexitcode $? "publishing deployment policy to the $ORG in the exchange"
        else
            output=$(HZN_ORG_ID=$ORG make publish-deployment-policy 2>&1)
            if [[ $? -ne 0 ]]; then
                echo "Error publishing deployment policy to the $ORG in the exchange: $output"
                exit 2
            fi
        fi
    fi
}

# git branch/repository to clone
branch="-b master"
repository="https://github.com/open-horizon/examples.git"

# text file containing servies and patterns to publish
input="$PATH_TO_EXAMPLES/examples/tools/blessedSamples.txt"

topDir=$(pwd)

if [[ -d "$PATH_TO_EXAMPLES/examples" ]]; then
    echo "Directory $PATH_TO_EXAMPLES/examples already exists, can not git clone into it. Will try to proceed assuming it was git cloned previously..."
else
    echo "Cloning $branch $repository to $PATH_TO_EXAMPLES/examples..."
    if [[ -n $VERBOSE ]]; then
        git clone $branch $repository $PATH_TO_EXAMPLES/examples
        checkexitcode $? "git clone $branch $repository $PATH_TO_EXAMPLES/examples"
    else
        output=$(git clone $branch $repository $PATH_TO_EXAMPLES/examples 2>&1)
        if [[ $? -ne 0 ]]; then
            echo "Error running git clone $branch $repository $PATH_TO_EXAMPLES/examples: $output"
            exit 2
        fi
    fi
fi

# read in blessedSamples.txt which contains the services, patterns, and policies to publish
while IFS= read -r line
do
    # each $line contains the path to any service/pattern/policy that needs to be published
    cd $PATH_TO_EXAMPLES/$line
    checkexitcode $? "finding service directory $line"
    
    echo "Publishing ${PWD}..."
    if [[ -n $VERBOSE ]]; then
        make publish-only
        checkexitcode $? "publishing services and patterns from $line to the exchange"
    else
        output=$(make publish-only 2>&1)
        if [[ $? -ne 0 ]]; then
            echo "Error publishing services and patterns from $line to the exchange: $output"
            exit 2
        fi
    fi

    # check if an org was specified to publish sample deployment policy 
    if [[ -n $ORG ]]; then
        deployPolPublish
    fi

    cd $topDir

done < "$input"


# clean up
echo -e "Successfully published all examples to the exchange. Removing $PATH_TO_EXAMPLES/examples directory."
rm -f -r $PATH_TO_EXAMPLES

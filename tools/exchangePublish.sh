#!/bin/bash

usage() {
    cat << EOF
Usage: ${0##*/} [-h] [-v] [-c <org-name>] [-X] [-e <examples-version>]

Flag:
  -c <org-name>    The exchange organization to publish example deployment policies to (the user's own org, not the IBM org).
  -X               Skip publishing patterns and services to the IBM org. Only valid in conjunction with -c <org-name> when publishing deployment policies to an additional org.
  -e <examples-version>   The branch of the examples repo to get the examples from, for example: v2.27 . Default: the first 2 numbers of the hzn version, preceded by 'v' .
  -v               Verbose output
  -h               This usage

Required Environment Variables:
  HZN_EXCHANGE_URL
  HZN_EXCHANGE_USER_AUTH

EOF
    exit $1
}

# Set global variables, that can be overridden by env vars
LOCAL_PATH_TO_EXAMPLES=${LOCAL_PATH_TO_EXAMPLES:-/tmp/open-horizon/examples}   # path to the cloned exmaples repo
EXCLUDE_IBM_PUBLISH=${EXCLUDE_IBM_PUBLISH:-false}
EXAMPLES_REPO=${EXAMPLES_REPO:-https://github.com/open-horizon/examples.git}

# other env vars that can be set (but default to blank)
# POLICY_ORG - org to publish deployment policy to
# EXAMPLES_PREVIEW_MODE - set to 'true' for testing/debugging this script

# if the org id is set locally we don't want to override the IBM org of these samples
unset HZN_ORG_ID

# Parse CLI arguments, overriding env vars where appropriate
while (( "$#" )); do
    case "$1" in
        -h) # display script usage
            usage
            shift
            ;;
        -v) # verbose output
            VERBOSE=1
            shift
            ;;
        -c) # org to publish deployment policy to
            POLICY_ORG=$2
            shift 2
            ;;
        -X) # exclude publishing to IBM org
            EXCLUDE_IBM_PUBLISH='true'
            shift
            ;;
        -e) # branch of the examples repo
            EXAMPLES_REPO_BRANCH=$2
            shift 2
            ;;
        -*) # invalid flag
            echo "ERROR: Unknow flag $1"
            usage 1
            ;;
        *) # there are no positional args
            echo "ERROR: Unknow argument $1"
            usage 1
            ;;
    esac
done

# Check input and requirements
: ${HZN_EXCHANGE_URL:?} ${HZN_EXCHANGE_USER_AUTH:?}

if [[ $EXCLUDE_IBM_PUBLISH == 'true' && -z $POLICY_ORG ]]; then
    echo "Error: if -X or EXCLUDE_IBM_PUBLISH is specified then -c or POLICY_ORG must also be specified"
    exit 1
fi

if ! command -v hzn >/dev/null 2>&1; then
    echo "Error: the 'hzn' command must be installed before running this script"
    exit 2
fi

# check the previous cmds exit code. 
chk() {
    local exitCode=$1
    local task=$2
    local dontExit=$3   # set to 'continue' to not exit for this error
    if [[ $exitCode == 0 ]]; then return; fi
    echo "Error: exit code $exitCode from: $task"
    if [[ $dontExit != 'continue' ]]; then
        exit $exitCode
    fi
}

# Run a command that does not have a quiet option, so we have to capture the output and only show if an error
runCmdQuietly() {
    # all of the args to this function are the cmd and its args
    if [[ -n $VERBOSE ]]; then
        $*
        chk $? "running: $*"
    else
        output=$($* 2>&1)
        local rc=$?
        if [[ $rc -ne 0 ]]; then
            echo "Error running $*: $output"
            exit $rc
        fi
    fi
}

# publish deployment policy for helloworld and cpu2evtstreams if -c flag is used
deployPolPublish() {
    local sample=${1:?}
    if ([[ $sample == *"cpu2evtstreams" ]] || [[ $sample == *"helloworld" ]] || [[ $sample == *"operator"* ]]); then 
        echo "Publishing deployment policy of $sample to $POLICY_ORG org..."
        if [[ $EXAMPLES_PREVIEW_MODE != 'true' ]]; then
            HZN_ORG_ID=$POLICY_ORG runCmdQuietly make publish-deployment-policy
        fi
    fi
}

origDir=$PWD

# Determine git branch to clone from
if [[ -z $EXAMPLES_REPO_BRANCH ]]; then
    hznVersion=$(hzn version | grep "^Horizon CLI")
    hznVersion=${hznVersion##* }   # remove all of the space-separated words, except the last one, so we are left with like: 2.27.0-123
    hznVersion=${hznVersion%.*}   # remove last dot and everything after, so left with like: 2.27
    EXAMPLES_REPO_BRANCH="v$hznVersion"
    echo "Using examples repo branch $EXAMPLES_REPO_BRANCH derived from the hzn version"
fi


# text file containing servies and patterns to publish
blessedSamples="$LOCAL_PATH_TO_EXAMPLES/tools/blessedSamples.txt"

# Clone the repo and switch to the branch
if [[ -d "$LOCAL_PATH_TO_EXAMPLES" ]]; then
    echo "Directory $LOCAL_PATH_TO_EXAMPLES already exists, can not git clone into it. Will try to proceed assuming it was git cloned previously..."
else
    echo "Cloning $EXAMPLES_REPO to $LOCAL_PATH_TO_EXAMPLES ..."
    #Adjust the default vaule of postBuffer for curl
    runCmdQuietly git config --global http.postBuffer 524288000
    #runCmdQuietly git clone $EXAMPLES_REPO_BRANCH $EXAMPLES_REPO $LOCAL_PATH_TO_EXAMPLES
    runCmdQuietly git clone $EXAMPLES_REPO $LOCAL_PATH_TO_EXAMPLES
fi
cd $LOCAL_PATH_TO_EXAMPLES
echo "Switching to branch $EXAMPLES_REPO_BRANCH ..."
if ! git checkout $EXAMPLES_REPO_BRANCH 2>/dev/null; then
    echo "Warning: examples branch '$EXAMPLES_REPO_BRANCH' does not exist, falling back to the master branch"
fi

if [[ $EXAMPLES_PREVIEW_MODE == 'true' ]]; then
    echo "Note: Running in preview mode, the samples will NOT actually be published..."
fi

# read in blessedSamples.txt which contains the services, patterns, and policies to publish
while IFS= read -r line
do
    # each $line contains the path to a service/pattern/policy that needs to be published
    sample=${line#examples/}   # in case we are using an older version of blessedSamples.txt
    cd $LOCAL_PATH_TO_EXAMPLES/$sample
    chk $? "finding service directory $sample"
    
    if [[ $EXCLUDE_IBM_PUBLISH != 'true' ]]; then
        echo "Publishing services and patterns of $sample to IBM org..."
        if [[ $EXAMPLES_PREVIEW_MODE != 'true' ]]; then
            runCmdQuietly make publish-only
        fi
    fi

    # check if an org was specified to publish sample deployment policy 
    if [[ -n $POLICY_ORG ]]; then
        deployPolPublish "$sample"
    fi

    cd $origDir

done < "$blessedSamples"


# clean up
echo -e "Successfully published all examples to the exchange. Removing $LOCAL_PATH_TO_EXAMPLES directory."
rm -f -r $LOCAL_PATH_TO_EXAMPLES

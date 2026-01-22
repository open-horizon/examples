#!/bin/bash

usage() {
    cat << EOF
Usage: ${0##*/} [-h] [-v] [-c <org-name>] [-X] [-e <examples-version>] [-a]

Flag:
  -c <org-name>             The exchange organization to publish example deployment policies to (the user's own org, not the IBM org).
  -X                        Skip publishing patterns and services to the IBM org. Only valid in conjunction with -c <org-name> when publishing deployment policies to an additional org.
  -a                        Use this flag to publish the example deployment policies in ALL available orgs.
  -e <examples-tag>         The tag of the examples repo to get the examples from, for example: v2.29.0-123. If you want the latest version of the examples, specify 'master'. Default: the CLI version returned by the 'hzn version' command, preceded by 'v'.
  -b <blessedSamples.txt>   Set this flag and provide your own blessedSamples file containing a list of sample service GitHub repositories.
  -v                        Verbose output
  -h                        This usage

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
# EXAMPLES_KEEP_LOCAL_REPO - set to 'true' to not remove the cloned/local repo at the end

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
        -a) # publish example policies to ALL orgs
            PUBLISH_ALL_ORGS='true'
            shift
            ;;
        -e) # tag of the examples repo
            EXAMPLES_REPO_TAG=$2
            shift 2
            ;;
        -b) # custom blessedSamples file
            BYO_SAMPLES=$2
            shift 2
            ;;
        -*) # invalid flag
            echo "ERROR: Unknown flag $1"
            usage 1
            ;;
        *) # there are no positional args
            echo "ERROR: Unknown argument $1"
            usage 1
            ;;
    esac
done

# Check input and requirements
: ${HZN_EXCHANGE_URL:?} ${HZN_EXCHANGE_USER_AUTH:?}

if [[ $EXCLUDE_IBM_PUBLISH == 'true' && -z $POLICY_ORG ]]; then
    echo "Error: if -X or EXCLUDE_IBM_PUBLISH is specified then -c or POLICY_ORG must also be specified, otherwise this script would not publish anything."
    exit 1
fi

if ! command -v hzn >/dev/null 2>&1; then
    echo "Error: the 'hzn' command must be installed before running this script"
    exit 2
fi

origDir=$PWD
# text file containing servies and patterns to publish
blessedSamples="$LOCAL_PATH_TO_EXAMPLES/tools/blessedSamples.txt"

# if using a custom blessedSamples file, it will need to be generated first
if [[ -n $BYO_SAMPLES ]]; then
    pathToBYOSamples="$origDir/$BYO_SAMPLES"
    if [[ ! -f "$pathToBYOSamples" ]]; then
        echo "Error: Custom samples file not found: $pathToBYOSamples"
        echo "Continuing with default samples..."
        unset BYO_SAMPLES
    fi
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
        if [[ $EXAMPLES_PREVIEW_MODE != 'true' ]]; then
            if [[ $PUBLISH_ALL_ORGS == 'true' ]]; then
                orgs=()
                orgs=( $(hzn exchange org list | jq .[]) )
                for org in ${orgs[@]}; do
                    if [[ $org != '"IBM"' && $org != '"root"' ]]; then
                        echo "Publishing deployment policy of $sample to $org org..."
                        HZN_ORG_ID=$org runCmdQuietly make publish-deployment-policy
                    fi
                done
            else
                echo "Publishing deployment policy of $sample to $POLICY_ORG org..."
                HZN_ORG_ID=$POLICY_ORG runCmdQuietly make publish-deployment-policy
            fi
        fi
    fi
}

publish_blessedSamples() {
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
                if [[ -n $BYO_SAMPLES ]]; then
                    HZN_ORG_ID=$POLICY_ORG runCmdQuietly make publish-only
                else
                    runCmdQuietly make publish-only
                fi
            fi
        fi

        # check if an org was specified to publish sample deployment policy 
        if [[ -n $POLICY_ORG || $PUBLISH_ALL_ORGS == 'true' ]]; then
            deployPolPublish "$sample"
        fi

        cd $origDir

    done < "$blessedSamples"
}

byo_blessedSamples() {
    echo "Publishing a custom set of example services to the Open Horizon Exchange..."
    > "$blessedSamples"  # Clear/create the file
    # generate a similarly formatted blessedSamples file
    while IFS= read -r line
    do
        # each $line contains the URL for a github repository
        repo=$(basename "$line")
        repo_name="${repo%.*}"
        git clone $line $LOCAL_PATH_TO_EXAMPLES/$repo_name
        chk $? "when cloning $line"

        echo "$repo_name" >> "$blessedSamples"

    done < "$pathToBYOSamples"
}

# Determine git tag to clone from
if [[ -z $EXAMPLES_REPO_TAG ]]; then
    EXAMPLES_REPO_TAG="v$(hzn version 2>/dev/null | grep 'Horizon CLI' | awk '{print $4}')"
    if [[ $EXAMPLES_REPO_TAG == 'v' ]]; then
        echo "Error: could not get CLI version from 'hzn version' "
        exit 3
    fi
    echo "Using examples repo tag $EXAMPLES_REPO_TAG derived from the hzn version"
fi

# Clone the repo at the specified tag point
if [[ -d "$LOCAL_PATH_TO_EXAMPLES" ]]; then
    echo "Directory $LOCAL_PATH_TO_EXAMPLES already exists, can not git clone into it. Will try to proceed assuming it was git cloned previously..."
else
    echo "Cloning $EXAMPLES_REPO to $LOCAL_PATH_TO_EXAMPLES ..."
    runCmdQuietly git clone $EXAMPLES_REPO $LOCAL_PATH_TO_EXAMPLES
fi
cd $LOCAL_PATH_TO_EXAMPLES
if [[ $EXAMPLES_REPO_TAG != 'master' ]]; then
    echo "Switching to tag $EXAMPLES_REPO_TAG ..."
    if git fetch origin +refs/tags/$EXAMPLES_REPO_TAG:refs/tags/$EXAMPLES_REPO_TAG 2>/dev/null; then
        git checkout tags/$EXAMPLES_REPO_TAG -b $EXAMPLES_REPO_TAG
        chk $? "switching to tag $EXAMPLES_REPO_TAG"
    else
        echo "Warning: examples tag '$EXAMPLES_REPO_TAG' does not exist, falling back to the master branch"
    fi
fi

if [[ $EXAMPLES_PREVIEW_MODE == 'true' ]]; then
    echo "Note: Running in preview mode, the samples will NOT actually be published..."
fi

# if using a custom blessedSamples file, it will need to be generated first
if [[ -n $BYO_SAMPLES ]]; then
    byo_blessedSamples
fi
publish_blessedSamples

# clean up
echo -e "Successfully published all examples to the exchange. Removing $LOCAL_PATH_TO_EXAMPLES directory."
if [[ $EXAMPLES_KEEP_LOCAL_REPO != 'true' ]]; then
    rm -f -r $LOCAL_PATH_TO_EXAMPLES
fi

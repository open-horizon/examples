#!/bin/bash

# This script gathers the necessary information and files to install Horizon and register an edge device

# Usage: ./script.sh <edge-device-type> [-t (to create a tar file containing gathered files)] <distribution>

# Parameters:
EDGE_DEVICE=$1 		# the type of edge device planned for install and registration < 32-bit-ARM , 64-bit-ARM , x86_64-LINUX , macOS >
PACKAGE_FILES=$2 	# if '-t' flag is set, create agentInstallFiles.tar.gz file containing gathered files 
DISTRO=$3			# for 64-bit-ARM and x86_64-LINUX 'xenial' instead of 'bionic' for the older version of Ubuntu

# Check if environment variables for `cloutctl login ...` are set 
function checkEnvVars () {
	echo "Checking environment variables..."
	echo ""

	if [ -z $EDGE_DEVICE ]; then
		echo "ERROR: Device type not specified."
		echo "Usage: ./script.sh <edge-device-type> [-t]"
		exit 1

	elif ! ([[ "$EDGE_DEVICE" == "32-bit-ARM" ]] || [[ "$EDGE_DEVICE" == "64-bit-ARM" ]] || [[ "$EDGE_DEVICE" == "x86_64-Linux" ]] || [[ "$EDGE_DEVICE" == "macOS" ]]); then
		echo "ERROR: Unknown device type."
		exit 1

	elif [ -z $ICP_URL ]; then
		echo "ERROR: ICP_URL environment variable is not set. Can not run 'cloudctl login ...'"
		exit 1 	

	elif [ -z $USER ]; then
		echo "ERROR: USER environment variable is not set. Can not run 'cloudctl login ...'"
		exit 1 

	elif [ -z $PW ]; then
		echo "ERROR: PW environment variable is not set. Can not run 'cloudctl login ...'"
		exit 1 
	fi
}

function cloudLogin () {
	echo "Connecting to cluster and configure kubectl..."
	echo ""

	cloudctl login -a $ICP_URL -u $USER -p $PW -n kube-public --skip-ssl-validation
	if [ $? -ne 0 ]; then
		echo "ERROR: 'cloudctl login' failed. Check if ICP_URL, USER, and PW environment variables are set correctly."
        exit 1
    fi
}

function getClusterName () {
	echo "Getting cluster name..."
	
	CLUSTER_NAME=$(kubectl get configmap -n kube-public ibmcloud-cluster-info -o jsonpath="{.data.cluster_name}")
	if [ $? -ne 0 ]; then
		echo "ERROR: Failed to get cluster name."
        exit 1
    fi

	echo "Cluster name: $CLUSTER_NAME"
	echo ""
}

function createAPIKey () {
	echo "Creating IBM Cloud Private platform API key..."

	cloudctl iam api-key-create "$EDGE_DEVICE API Key" -d "$EDGE_DEVICE API Key" -f key.txt
	if [ $? -ne 0 ]; then
		echo "ERROR: Failed to create API Key."
        exit 1
    fi

    API_KEY=$(cat key.txt | jq -r '.apikey')
    echo "API Key: $API_KEY"
    echo ""
}

function createAgentInstallConfig () {
	echo "Creating agent-install.cfg file..."

	cat << EndOfContent > agent-install.cfg
HZN_EXCHANGE_URL=$ICP_URL/ec-exchange/v1/
HZN_FSS_CSSURL=$ICP_URL/ec-css/
HZN_ORG_ID=$CLUSTER_NAME
HZN_EXCHANGE_USER_AUTH=iamapikey:$API_KEY
HZN_EXCHANGE_PATTERN=IBM/pattern-ibm.helloworld
EndOfContent
	if [ $? -ne 0 ]; then
		echo "ERROR: Failed to create agent-install.cfg file."
        exit 1
    fi

    echo "agent-install.cfg file created: "
	cat agent-install.cfg
	echo ""
}

function getICPCert () {
	echo "Getting the IBM Cloud Private self-signed certificate..."

	kubectl -n kube-public get secret ibmcloud-cluster-ca-cert -o jsonpath="{.data['ca\.crt']}" | base64 --decode > agent-install.crt
	if [ $? -ne 0 ]; then
		echo "ERROR: Failed to get the IBM Cloud Private self-signed certificate"
        exit 1
    fi
    echo ""
}

function gatherHorizonFiles () {
	echo "Locating the IBM Edge Computing for Devices installation content for $EDGE_DEVICE device..."

    # Determine edge device type
    if [[ "$EDGE_DEVICE" == "32-bit-ARM" ]]; then
		tar --strip-components 6 -zxvf ibm-edge-computing-x86_64-3.2.1.1.tar.gz ibm-edge-computing-x86_64-3.2.1.1/horizon-edge-packages/linux/raspbian/stretch/armhf
		if [ $? -ne 0 ]; then
			echo "ERROR: Failed to locate the IBM Edge Computing for Devices installation content"
        	exit 1
    	fi

	elif [[ "$EDGE_DEVICE" == "64-bit-ARM" ]]; then
		tar --strip-components 6 -zxvf ibm-edge-computing-x86_64-3.2.1.1.tar.gz ibm-edge-computing-x86_64-3.2.1.1/horizon-edge-packages/linux/ubuntu/bionic/arm64
		if [ $? -ne 0 ]; then
			echo "ERROR: Failed to locate the IBM Edge Computing for Devices installation content"
        	exit 1
    	fi

	elif [[ "$EDGE_DEVICE" == "x86_64-Linux" ]]; then
		tar --strip-components 6 -zxvf ibm-edge-computing-x86_64-3.2.1.1.tar.gz ibm-edge-computing-x86_64-3.2.1.1/horizon-edge-packages/linux/ubuntu/bionic/amd64
		if [ $? -ne 0 ]; then
			echo "ERROR: Failed to locate the IBM Edge Computing for Devices installation content"
        	exit 1
    	fi

	elif [[ "$EDGE_DEVICE" == "macOS" ]]; then
		tar --strip-components 3 -zxvf ibm-edge-computing-x86_64-3.2.1.1.tar.gz ibm-edge-computing-x86_64-3.2.1.1/horizon-edge-packages/macos
		if [ $? -ne 0 ]; then
			echo "ERROR: Failed to locate the IBM Edge Computing for Devices installation content"
        	exit 1
    	fi

	else
		echo "ERROR: Unknown device type."
		exit 1
	fi
	echo ""
}

function createTarFile () {
	echo "Creating agentInstallFiles.tar.gz file containing gathered files..."

	tar -czvf agentInstallFiles.tar.gz $(ls agent-install.cfg agent-install.crt *horizon*)
	if [ $? -ne 0 ]; then
		echo "ERROR: Failed to create agentInstallFiles.tar.gz file."
       	exit 1
    fi
	echo ""
}

main () {
	checkEnvVars

	cloudLogin

	getClusterName

	createAPIKey

	createAgentInstallConfig

	getICPCert

	gatherHorizonFiles

	if [[ "$PACKAGE_FILES" == "-t" ]]; then
		createTarFile
	fi
}
main



#!/bin/bash

# This script gathers the necessary information and files to install Horizon and register an edge device
#
# Usage: ./script.sh <edge-device-type> [-d <distribution>] [-t]
#
# Parameters:
#	required: 
# 		<edge-device-type>		the type of edge device planned for install and registration 
#								accepted values: < 32-bit-ARM , 64-bit-ARM , x86_64-LINUX , macOS >
#
#	optional: 	
#		-t 						create agentInstallFiles.tar.gz file containing gathered files
#		-d <distribution>		script defaults to 'bionic' build on linux
#								use this flag with < 64-bit-ARM or x86_64-LINUX > 
#								to specify `xenial` build 
#								Flag is ignored with < macOS >


# Parse script parameters
while (( "$#" )); do
  	case "$1" in
    	-d) # distribution specified 
      		if ! ([[ "$2" == "xenial" ]] || [[ "$2" == "bionic" ]]); then
      			echo "ERROR: Unknown linux distribution type."
      			exit 1
      		fi
      		DISTRO=$2
      		shift 2
      		;;
    	-t) # create tar file
      		PACKAGE_FILES=$1
      		shift
     		;;
    	*) # based on "Usage" this should be device type 
			if ! ([[ "$1" == "32-bit-ARM" ]] || [[ "$1" == "64-bit-ARM" ]] || [[ "$1" == "x86_64-Linux" ]] || [[ "$1" == "macOS" ]]); then
				echo "ERROR: Unknown device type."
				exit 1
			fi
      		EDGE_DEVICE=$1
      		shift
      		;;
  	esac
done
if [ -z $EDGE_DEVICE ]; then
	echo "ERROR: Device type not specified."
	echo "Usage: ./script.sh <edge-device-type> [-t]"
	exit 1
fi

 
function checkEnvVars () {
	echo "Checking system requirements..."
	cloudctl --help
	if [ $? -ne 0 ]; then
		echo "ERROR: cloudctl is not installed."
        exit 1
    fi

    kubectl --help
	if [ $? -ne 0 ]; then
		echo "ERROR: kubectl is not installed."
        exit 1
    fi

	echo "Checking environment variables..."

	if [ -z $ICP_URL ]; then
		echo "ERROR: ICP_URL environment variable is not set. Can not run 'cloudctl login ...'"
		exit 1 	

	elif [ -z $USER ]; then
		echo "ERROR: USER environment variable is not set. Can not run 'cloudctl login ...'"
		exit 1 

	elif [ -z $PW ]; then
		echo "ERROR: PW environment variable is not set. Can not run 'cloudctl login ...'"
		exit 1 
	fi
	echo ""
}

function cloudLogin () {
	echo "Connecting to cluster and configure kubectl..."

	cloudctl login -a $ICP_URL -u $USER -p $PW -n kube-public --skip-ssl-validation
	if [ $? -ne 0 ]; then
		echo "ERROR: 'cloudctl login' failed. Check if ICP_URL, USER, and PW environment variables are set correctly."
        exit 1
    fi
    echo ""
}

# Query the IBM Cloud Private cluster name
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

# Create a IBM Cloud Private platform API key
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

# With the information from the previous functions, create agent-install.cfg
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

# Get the IBM Cloud Private self-signed certificate
function getICPCert () {
	echo "Getting the IBM Cloud Private self-signed certificate agent-install.crt..."

	kubectl -n kube-public get secret ibmcloud-cluster-ca-cert -o jsonpath="{.data['ca\.crt']}" | base64 --decode > agent-install.crt
	if [ $? -ne 0 ]; then
		echo "ERROR: Failed to get the IBM Cloud Private self-signed certificate"
        exit 1
    fi
    echo ""
}

# Locate the IBM Edge Computing for Devices installation content
function gatherHorizonFiles () {
	echo "Locating the IBM Edge Computing for Devices installation content for $EDGE_DEVICE device..."

    # Determine edge device type, and distribution if applicable 
    if [[ "$EDGE_DEVICE" == "32-bit-ARM" ]]; then
		tar --strip-components 6 -zxvf ibm-edge-computing-x86_64-3.2.1.1.tar.gz ibm-edge-computing-x86_64-3.2.1.1/horizon-edge-packages/linux/raspbian/stretch/armhf
		if [ $? -ne 0 ]; then
			echo "ERROR: Failed to locate the IBM Edge Computing for Devices installation content"
        	exit 1
    	fi

	elif [[ "$EDGE_DEVICE" == "64-bit-ARM" ]]; then
		if [[ "$DISTRO" == "xenial" ]]; then
			tar --strip-components 6 -zxvf ibm-edge-computing-x86_64-3.2.1.1.tar.gz ibm-edge-computing-x86_64-3.2.1.1/horizon-edge-packages/linux/ubuntu/xenial/arm64
		else
			tar --strip-components 6 -zxvf ibm-edge-computing-x86_64-3.2.1.1.tar.gz ibm-edge-computing-x86_64-3.2.1.1/horizon-edge-packages/linux/ubuntu/bionic/arm64
		fi
		if [ $? -ne 0 ]; then
			echo "ERROR: Failed to locate the IBM Edge Computing for Devices installation content"
        	exit 1
    	fi

	elif [[ "$EDGE_DEVICE" == "x86_64-Linux" ]]; then
		if [[ "$DISTRO" == "xenial" ]]; then
			tar --strip-components 6 -zxvf ibm-edge-computing-x86_64-3.2.1.1.tar.gz ibm-edge-computing-x86_64-3.2.1.1/horizon-edge-packages/linux/ubuntu/xenial/amd64
		else	
			tar --strip-components 6 -zxvf ibm-edge-computing-x86_64-3.2.1.1.tar.gz ibm-edge-computing-x86_64-3.2.1.1/horizon-edge-packages/linux/ubuntu/bionic/amd64
		fi
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

# Download the latest version of the agent-install.sh script and make it executable
function pullAgentInstallScript () {
	echo "Pulling agent-install.sh script..."

	curl -O https://raw.githubusercontent.com/open-horizon/anax/v3.2.1/agent-install/agent-install.sh && \
		chmod +x ./agent-install.sh
	if [ $? -ne 0 ]; then
		echo "ERROR: Failed to pull agent-install.sh script from the anax repo."
       	exit 1
    fi
    echo ""
}

# Create a tar file of the gathered files for batch install 
function createTarFile () {
	echo "Creating agentInstallFiles.tar.gz file containing gathered files..."

	tar -czvf agentInstallFiles.tar.gz $(ls agent-install.sh agent-install.cfg agent-install.crt *horizon*)
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

	pullAgentInstallScript

	if [[ "$PACKAGE_FILES" == "-t" ]]; then
		createTarFile
	fi
}
main



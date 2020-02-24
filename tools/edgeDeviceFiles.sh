#!/bin/bash

# This script gathers the necessary information and files to install Horizon and register an edge device


function scriptUsage () {
	cat << EOF
ERROR: No arguments specified.

Usage: ./edgeDeviceFiles.sh <edge-device-type> [-t] [-k] [-d <distribution>] [-f <directory>]

Parameters:
  required:
    <edge-device-type>		the type of edge device planned for install and registration
				  accepted values: < 32-bit-ARM , 64-bit-ARM , x86_64-Linux , macOS >

  optional:
    -t 				create agentInstallFiles.tar.gz file containing gathered files
				  If this flag isn't set, the gathered files will be placed in the current directory
    -d 				<distribution>	script defaults to 'bionic' build on linux
				  use this flag with < 64-bit-ARM or x86_64-Linux >
				  to specify 'xenial' build
                                  use this with < 32-bit-ARM > to specify 'stretch' instead of defauly, 'buster'
				  Flag is ignored with < macOS >
    -k 				include this flag to create a new $USER-Edge-Device-API-Key. If this flag is not set,
				  the existing api keys will be checked for $USER-Edge-Device-API-Key and creation will
				  be skipped if it exists
    -f 				<directory> to move gathered files to. Default is current directory

Required Environment Variables:
    CLUSTER_URL			https://<cluster_CA_domain>:<port-number>
    USER 			your-cluster-admin-user
    PW				your-cluster-admin-password

EOF
	exit 1
}
if [[ "$#" = "0" ]]; then
	scriptUsage
fi

echo "Checking script parameters..."
while (( "$#" )); do
  	case "$1" in
    	-d) # distribution specified
      		if ! ([[ "$2" == "xenial" ]] \
					|| [[ "$2" == "bionic" ]] \
					|| [[ "$2" == "stretch" ]] \
					|| [[ "$2" == "buster" ]]); then
      			echo "ERROR: Unknown linux distribution type."
      			echo ""
      			exit 1
      		fi
      		DISTRO=$2
      		shift 2
      		;;
    	-t) # create tar file
      		PACKAGE_FILES=$1
      		shift
     		;;
     	-k) # create api key
      		CREATE_API_KEY=$1
      		shift
     		;;
     	-f) # directory to move gathered files to
			DIR=$2
      		shift 2
      		;;
    	*) # based on "Usage" this should be device type
			if ! ([[ "$1" == "32-bit-ARM" ]] || [[ "$1" == "64-bit-ARM" ]] || [[ "$1" == "x86_64-Linux" ]] || [[ "$1" == "macOS" ]]); then
				echo "ERROR: Unknown device type."
				echo ""
				exit 1
			fi
      		EDGE_DEVICE=$1
      		shift
      		;;
  	esac
done
if [ -z $EDGE_DEVICE ]; then
	scriptUsage
fi
echo " - valid parameters"
echo ""


function checkEnvVars () {
	echo "Checking system requirements..."
	cloudctl --help > /dev/null 2>&1
	if [ $? -ne 0 ]; then
		echo "ERROR: cloudctl is not installed."
        echo ""
        exit 1
    fi
    echo " - cloudctl installed"

    kubectl --help > /dev/null 2>&1
	if [ $? -ne 0 ]; then
		echo "ERROR: kubectl is not installed."
        echo ""
        exit 1
    fi
    echo " - kubectl installed"
    echo ""

	echo "Checking environment variables..."

	if [ -z $CLUSTER_URL ]; then
		echo "ERROR: CLUSTER_URL environment variable is not set. Can not run 'cloudctl login ...'"
		echo " - CLUSTER_URL=https://<cluster_CA_domain>:<port-number>"
		echo ""
		exit 1

	elif [ -z $USER ]; then
		echo "ERROR: USER environment variable is not set. Can not run 'cloudctl login ...'"
		echo " - USER=<your-cluster-admin-user>"
		echo ""
		exit 1

	elif [ -z $PW ]; then
		echo "ERROR: PW environment variable is not set. Can not run 'cloudctl login ...'"
		echo " - PW=<your-cluster-admin-password>"
		echo ""
		exit 1
	fi
	echo " - CLUSTER_URL set"
	echo " - USER set"
	echo " - PW set"
	echo ""
}

function cloudLogin () {
	echo "Connecting to cluster and configure kubectl..."
	echo "cloudctl login -a $CLUSTER_URL -u $USER -p $PW -n kube-public --skip-ssl-validation"

	cloudctl login -a $CLUSTER_URL -u $USER -p $PW -n kube-public --skip-ssl-validation
	if [ $? -ne 0 ]; then
		echo "ERROR: 'cloudctl login' failed. Check if CLUSTER_URL, USER, and PW environment variables are set correctly."
        echo ""
        exit 1
    fi
    echo ""
}

# Query the IBM Cloud Pak cluster name
function getClusterName () {
	echo "Getting cluster name..."
	echo "kubectl get configmap -n kube-public ibmcloud-cluster-info -o jsonpath=\"{.data.cluster_name}\""

	CLUSTER_NAME=$(kubectl get configmap -n kube-public ibmcloud-cluster-info -o jsonpath="{.data.cluster_name}")
	if [ $? -ne 0 ]; then
		echo "ERROR: Failed to get cluster name."
        echo ""
        exit 1
    fi

	echo " - Cluster name: $CLUSTER_NAME"
	echo ""
}

# Check if an IBM Cloud Pak platform API key exists
function checkAPIKey () {
	echo "Checking if \"$USER-Edge-Device-API-Key\" already exists..."
	echo "cloudctl iam api-keys | cut -d' ' -f4 | grep \"$USER-Edge-Device-API-Key\""

	KEY=$(cloudctl iam api-keys | cut -d' ' -f4 | grep "$USER-Edge-Device-API-Key")
	if [ -z $KEY ]; then
		echo "\"$USER-Edge-Device-API-Key\" does not exist. A new one will be created."
        CREATE_NEW_KEY=true
    else
    	echo "\"$USER-Edge-Device-API-Key\" already exists. Skipping key creation."
    	CREATE_NEW_KEY=false
    fi
    echo ""
}

# Create a IBM Cloud Pak platform API key
function createAPIKey () {
	echo "Creating IBM Cloud Pak platform API key..."
	echo "cloudctl iam api-key-create \"$USER-Edge-Device-API-Key\" -d \"$USER-Edge-Device-API-Key\" -f key.txt"

	cloudctl iam api-key-create "$USER-Edge-Device-API-Key" -d "$USER-Edge-Device-API-Key" -f key.txt
	if [ $? -ne 0 ]; then
		echo "ERROR: Failed to create API Key."
        echo ""
        exit 1
    fi

    API_KEY=$(cat key.txt | jq -r '.apikey')
    echo " - $USER-Edge-Device-API-Key: $API_KEY"
    echo ""
}

# With the information from the previous functions, create agent-install.cfg
function createAgentInstallConfig () {
	echo "Creating agent-install.cfg file..."

	cat << EndOfContent > agent-install.cfg
HZN_EXCHANGE_URL=$CLUSTER_URL/ec-exchange/v1/
HZN_FSS_CSSURL=$CLUSTER_URL/ec-css/
HZN_ORG_ID=$CLUSTER_NAME
EndOfContent
	if [ $? -ne 0 ]; then
		echo "ERROR: Failed to create agent-install.cfg file."
        echo ""
        exit 1
    fi

    echo "agent-install.cfg file created: "
	cat agent-install.cfg
	echo ""
}

# Get the IBM Cloud Pak self-signed certificate
function getClusterCert () {
	echo "Getting the IBM Cloud Pak self-signed certificate agent-install.crt..."
	echo "kubectl -n kube-public get secret ibmcloud-cluster-ca-cert -o jsonpath=\"{.data['ca\.crt']}\" | base64 --decode > agent-install.crt"

	kubectl --namespace kube-system get secret cluster-ca-cert -o jsonpath="{.data['tls\.crt']}" | base64 --decode > agent-install.crt
	if [ $? -ne 0 ]; then
		echo "ERROR: Failed to get the IBM Cloud Pak self-signed certificate"
        echo ""
        exit 1
    fi
    echo ""
}

# Locate the IBM Edge Computing for Devices installation content
function gatherHorizonFiles () {
	echo "Locating the IBM Edge Computing Manager for Devices installation content for $EDGE_DEVICE device..."
	echo "tar --strip-components n -zxvf ibm-ecm-4.0.0-x86_64.tar.gz ibm-ecm-4.0.0-x86_64/horizon-edge-packages/..."
	echo "Dist is $DISTRO"

    # Determine edge device type, and distribution if applicable
    if [[ "$EDGE_DEVICE" == "32-bit-ARM" ]]; then
		if [[ "$DISTRO" == "stretch" ]]; then
			tar --strip-components 6 -zxvf ibm-ecm-4.0.0-x86_64.tar.gz ibm-ecm-4.0.0-x86_64/horizon-edge-packages/linux/raspbian/stretch/armhf
		else
			tar --strip-components 6 -zxvf ibm-ecm-4.0.0-x86_64.tar.gz ibm-ecm-4.0.0-x86_64/horizon-edge-packages/linux/raspbian/buster/armhf
		fi
		if [ $? -ne 0 ]; then
			echo "ERROR: Failed to locate the IBM Edge Computing Manager for Devices installation content"
        	echo ""
        	exit 1
    	fi

	elif [[ "$EDGE_DEVICE" == "64-bit-ARM" ]]; then
		if [[ "$DISTRO" == "xenial" ]]; then
			tar --strip-components 6 -zxvf ibm-ecm-4.0.0-x86_64.tar.gz ibm-ecm-4.0.0-x86_64/horizon-edge-packages/linux/ubuntu/xenial/arm64
		else
			tar --strip-components 6 -zxvf ibm-ecm-4.0.0-x86_64.tar.gz ibm-ecm-4.0.0-x86_64/horizon-edge-packages/linux/ubuntu/bionic/arm64
		fi
		if [ $? -ne 0 ]; then
			echo "ERROR: Failed to locate the IBM Edge Computing Manager for Devices installation content"
        	echo ""
        	exit 1
    	fi

	elif [[ "$EDGE_DEVICE" == "x86_64-Linux" ]]; then
		if [[ "$DISTRO" == "xenial" ]]; then
			tar --strip-components 6 -zxvf ibm-ecm-4.0.0-x86_64.tar.gz ibm-ecm-4.0.0-x86_64/horizon-edge-packages/linux/ubuntu/xenial/amd64
		else
			tar --strip-components 6 -zxvf ibm-ecm-4.0.0-x86_64.tar.gz ibm-ecm-4.0.0-x86_64/horizon-edge-packages/linux/ubuntu/bionic/amd64
		fi
		if [ $? -ne 0 ]; then
			echo "ERROR: Failed to locate the IBM Edge Computing Manager for Devices installation content"
        	echo ""
        	exit 1
    	fi

	elif [[ "$EDGE_DEVICE" == "macOS" ]]; then
		tar --strip-components 3 -zxvf ibm-ecm-4.0.0-x86_64.tar.gz ibm-ecm-4.0.0-x86_64/horizon-edge-packages/macos
		if [ $? -ne 0 ]; then
			echo "ERROR: Failed to locate the IBM Edge Computing Manager for Devices installation content"
        	echo ""
        	exit 1
    	fi

	else
		echo "ERROR: Unknown device type."
		echo ""
		exit 1
	fi
	echo ""
}

# Download the latest version of the agent-install.sh script and make it executable
function pullAgentInstallScript () {
	echo "Pulling agent-install.sh script..."

	curl -O https://raw.githubusercontent.com/open-horizon/anax/v4.0/agent-install/agent-install.sh && \
		chmod +x ./agent-install.sh
	if [ $? -ne 0 ]; then
		echo "ERROR: Failed to pull agent-install.sh script from the anax repo."
       	echo ""
       	exit 1
    fi
    echo ""
}

# Create a tar file of the gathered files for batch install
function createTarFile () {
	echo "Creating agentInstallFiles-$EDGE_DEVICE.tar.gz file containing gathered files..."
	echo "tar -czvf agentInstallFiles-$EDGE_DEVICE.tar.gz \$(ls agent-install.sh agent-install.cfg agent-install.crt *horizon*)"

	tar -czvf agentInstallFiles-$EDGE_DEVICE.tar.gz $(ls agent-install.sh agent-install.cfg agent-install.crt *horizon*)
	if [ $? -ne 0 ]; then
		echo "ERROR: Failed to create agentInstallFiles-$EDGE_DEVICE.tar.gz file."
       	echo ""
       	exit 1
    fi
	echo ""
}

# Move gathered files to specified -f directory
function moveFiles () {
	echo "Moving files to $DIR..."
	if ! [[ -d "$DIR" ]]; then
    	echo "$DIR does not exist, creating it..."
    	mkdir $DIR
	fi

	mv $(ls agent-install.sh agent-install.cfg agent-install.crt *horizon*) $DIR
	if [ -f key.txt ]; then
    	mv key.txt $DIR
	fi

	if [ $? -ne 0 ]; then
		echo "ERROR: Failed to move files to $DIR."
       	echo ""
       	exit 1
    fi
    echo ""
}

# If an API Key was created, print it out
function printApiKey () {
	echo ""
	echo "************************** Your created API Key ******************************"
	echo ""
	echo "     $USER-Edge-Device-API-Key: $API_KEY"
	echo ""
	echo "********************* Save this value for future use *************************"
	echo ""
}

main () {
	checkEnvVars

	cloudLogin

	getClusterName

	checkAPIKey

	if [[ "$CREATE_API_KEY" == "-k" ]] || [[ "$CREATE_NEW_KEY" == "true" ]]; then
		createAPIKey
	fi

	createAgentInstallConfig

	getClusterCert

	gatherHorizonFiles

	pullAgentInstallScript

	if [[ "$PACKAGE_FILES" == "-t" ]]; then
		createTarFile
	fi

	if ! [ -z $DIR ]; then
		moveFiles
	fi

	if [[ "$CREATE_API_KEY" == "-k" ]] || [[ "$CREATE_NEW_KEY" == "true" ]]; then
		printApiKey
	fi

}
main



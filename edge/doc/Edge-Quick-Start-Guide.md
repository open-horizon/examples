# Edge Quick Start Guide

This guide provides a concise description of the process for setting up WIoTP/Horizon edge nodes and deploying existing services to them.

If you want to develop your own Horizon service, see the [Edge Developer Quick Start Guide](Edge-Developer-Quickstart-Guide.md).

Additional information is available, and questions may be asked, in our forum, at [https://discourse.bluehorizon.network/](https://discourse.bluehorizon.network/).

The Edge is based upon the open source Horizon project [https://github.com/open-horizon](https://github.com/open-horizon). There are therefore several references to Horizon, and the "hzn" Linux shell command in this document.

Let's get started...

## Setup Your Organization in the Watson IoT Platform
* Visit the IBM Watson IoT Platform pages and create an IBM cloud account and sign in
  * [https://console.bluemix.net/](https://console.bluemix.net/)
* Provision a Watson Internet of Things instance
  * click `Create resource`
  * on the left, under the `Platform` category, click `Internet of Things`
  * click `Internet of Things Platform`
  * **for now you must select `US South` as the region**
  * click Create
* Start your Internet of Things Platform instance:
  * click `Launch` to start the instance
  * Save your 6-character organization name. You can find it as the prefix of the web page's domain name, or at the top right in the ID field.
* Enable the beta Edge features in your IoT Platform instance:
  * on the left, hover over the gear and click `SETTINGS`
  * click `Experimental Features`
  * enable `Activate Experimental Features`
* Create a Gateway Type and instance:
  * on the left, hover over the device icon and click `DEVICES`
  * click the `Device Types` tab
  * click `Add Device Type`
  * click `Gateway`
  * enter a name for your new IoT Gateway Type
  * enable `Edge Services` for this gateway type
  * click `Next`, `Next`, and then click `Done`
  * click the `Browse` tab
  * click `Add Device`
  * select the gateway type you just created (it may be pre-selected)
  * enter a device ID for your edge node device
  * click `Next` and `Next` again
  * enter an authentication token and record it somewhere
  * click `Next` and click `Done`
* Create an API key so that the Edge node can access your Internet of Things Platform instance:
  * on the left, hover over the compass icon and click `APPS`
  * click `Generate API Key`
  * enter a Description
  * click `Next`
  * select Standard Application
  * click `Generate Key`
  * record the API Key and Authentication Token somewhere
* When finished, you should have the following information recorded:
  * An IBM cloud account name and password
  * Organization ID (6 characters).
  * A Gateway Type name
  * A Gateway instance ID, and its associated authentication token
  * An API Key name, and its associated authentication token

## Prepare Your Edge Node
* Install a recent Debian Linux variant on your Edge Node (e.g., ubuntu 16.04, which is used for the instructions below)
* Open a Linux shell with root privileges on your Edge Node (e.g., on the console, or over ssh)
```
sudo -s
```
* Install some utilities:
```
apt update && apt install -y curl wget gettext
```
* Ensure that you have the current docker version installed (since many distros are set up to run much older docker versions):
```
curl -fsSL get.docker.com | sh
```
* Configure the *apt* manager by adding the bluehorizon repo to /etc/apt/sources.list.d:
```
wget -qO - http://pkg.bluehorizon.network/bluehorizon.network-public.key | apt-key add -
aptrepo=updates
# aptrepo=testing    # or use this for the latest, development version
cat <<EOF > /etc/apt/sources.list.d/bluehorizon.list
deb [arch=$(dpkg --print-architecture)] http://pkg.bluehorizon.network/linux/ubuntu xenial-$aptrepo main
deb-src [arch=$(dpkg --print-architecture)] http://pkg.bluehorizon.network/linux/ubuntu xenial-$aptrepo main
EOF
```
* Install the horizon packages and MQTT client:
```
apt update && apt install -y horizon-wiotp mosquitto-clients
```
* Make sure the horizon package version shown at the bottom of the above step is "2.15.2" or later

For the rest of the guide you will not require root privileges, so you may optionally exit now from the root privileged shell you created above.

* The remaining commands shown in this document expect you to have the following environment variables set in your Linux shell environment.  Put these into a file, replacing the values that have "my" in them with your own values you recorded in the first section of the document.  Then source this file in your shell.

```
# These values contain the credentials you created earlier in the Watson IoT Platform web GUI
export HZN_ORG_ID=myorg
export WIOTP_DOMAIN=internetofthings.ibmcloud.com
export WIOTP_GW_TYPE=mygwtype
export WIOTP_GW_ID=mygwinstance
export WIOTP_GW_TOKEN='mygwinstancetoken'
export WIOTP_API_KEY='a-myapikeyrandomchars'
export WIOTP_API_TOKEN='myapikeytoken'

# This variable must be set appropriately for your specific Edge Node
export ARCH=amd64   # or arm for Raspberry Pi, or arm64 for TX2

# There is no need for you to edit these variables
export HZN_DEVICE_ID="g@${WIOTP_GW_TYPE}@$WIOTP_GW_ID"
export HZN_DEVICE_TOKEN="$WIOTP_GW_TOKEN"
export WIOTP_CLIENT_ID_APP="a:$HZN_ORG_ID:$WIOTP_GW_TYPE$WIOTP_GW_ID"
export WIOTP_CLIENT_ID_GW="g:$HZN_ORG_ID:$WIOTP_GW_TYPE:$WIOTP_GW_ID"
export HZN_EXCHANGE_USER_AUTH="$WIOTP_API_KEY:$WIOTP_API_TOKEN"
export HZN_EXCHANGE_API_AUTH="$WIOTP_API_KEY:$WIOTP_API_TOKEN"
```
## Verify Your Gateway Credentials and Access

List your gateway instance from the WIoTP cloud:
```
hzn wiotp device list $WIOTP_GW_TYPE $WIOTP_GW_ID | jq .
``` 

Use the mosquitto-clients package to verify your credentials by opening two Linux shells and subscribing to the IBM Watson IoT Platform MQTT message broker in one shell, and publishing a message to that broker in the other (which you should see in the subscribed shell).
* In the first shell, subscribe:
```
mosquitto_sub -v -h $HZN_ORG_ID.messaging.$WIOTP_DOMAIN -p 8883 -i "${WIOTP_CLIENT_ID_APP:0:38}" -u "$WIOTP_API_KEY" -P "$WIOTP_API_TOKEN" --capath /etc/ssl/certs -t iot-2/type/$WIOTP_GW_TYPE/id/$WIOTP_GW_ID/evt/status/fmt/json
```
* In the other shell, publish
```
mosquitto_pub -h $HZN_ORG_ID.messaging.$WIOTP_DOMAIN -p 8883 -i "$WIOTP_CLIENT_ID_GW" -u "use-token-auth" -P "$WIOTP_GW_TOKEN" --capath /etc/ssl/certs -t iot-2/type/$WIOTP_GW_TYPE/id/$WIOTP_GW_ID/evt/status/fmt/json -m '{"message": "Hello, world."}'
```
* You should see the "Hello, world." message appear in the output of the first shell

## Define an Additional Microservice and Workload in the Horizon Exchange

At this point, you could register your edge node with Horizon and have the default WIoTP core-iot service deployed to it. But we want to also show you how to have an additional workload deployed to your edge nodes. For this we will use a simple example microservice and workload created by the Horizon team that gathers CPU load statistics from your Edge node and publishes those statistics to your Watson IoT Platform.

* Clone the openhorizon examples project which contains files that you will need during the following steps:
```bash
cd ~
git clone https://github.com/open-horizon/examples.git
```
* Temporarily set the Horizon Exchange URL
```bash
export HZN_EXCHANGE_URL="https://$HZN_ORG_ID.$WIOTP_DOMAIN/api/v0002/edgenode/"
```
* Add the "cpu" microservice to your WIoTP organization and see that it was added:
```bash
hzn exchange microservice publish -f ~/examples/edge/services/cpu_percent/horizon.microservice/pre-signed/cpu-$ARCH.json
hzn exchange microservice list | jq .
```

* Configure the CPU load workload definition file using your environment variables, add it to your WIoTP organization, and see that it was added:
```bash
mkdir -p ~/hzn
hzn exchange workload publish -f ~/examples/edge/wiotp/cpu2wiotp/horizon/pre-signed/cpu2wiotp-$ARCH.json
hzn exchange workload list | jq .
```

## Augment the Edge Node Deployment Pattern

The Edge system deploys Patterns of code onto WIoTP Edge Node gateways. The deployment Pattern used for a particular gateway has the same name as its Gateway Type. By default the deployment pattern includes the WIoTP core-IoT service. Here you will update the deployment Pattern for your Gateway Type so that it also includes the prebuilt example CPU load Workload that we just added to your platform.

* Configure the CPU load pattern json file using your environment variables and add it to your pattern:
```bash
hzn exchange pattern insertworkload -f ~/examples/edge/wiotp/cpu2wiotp/horizon/pattern/insert-cpu2wiotp.json $WIOTP_GW_TYPE
```
* Verify that the CPU load Workload was inserted into the Pattern for your Gateway Type:
```bash
hzn exchange pattern list $WIOTP_GW_TYPE | jq .
```
* Unset the HZN_EXCHANGE_URL environment variable (because after registration in the next section `hzn` can get the value from the Horizon agent):
```bash
unset HZN_EXCHANGE_URL
```

## Register Your Edge Node
Now register your node with the Edge system and verify that it is connected properly before proceeding to develop and publish your code.

* **If you have run thru this document before** on this edge node, do this to clean up:
```
hzn unregister -f
```
* Register the node and start the Watson IoT Platform core-IoT service and the CPU workload:
```
wiotp_agent_setup --org $HZN_ORG_ID --deviceType $WIOTP_GW_TYPE --deviceId $WIOTP_GW_ID --deviceToken "$WIOTP_GW_TOKEN" -cn 'edge-connector'
```
After a short while, usually within just a minute or two (but rarely it could take up to 10 minutes) the Horizon Agreement Bots (AgBots) in the IBM Cloud will discover your Edge Node and establish an agreement with it to run all of the containers referenced in the deployment pattern.
* Verify that 2 agreements are made, one for `edge-core-iot-workload` and one for `cpu2wiotp`.  The output should indicate an `agreement_finalized_time`, and eventually an `agreement_execution_start_time` should also be populated.
```
hzn agreement list | jq . 
```

* Eventually the 6 docker containers should be running: 4 for the core-IoT service and 2 for the CPU example.  Verify this with:
```
docker ps
```

* If an agreement is not formed, or if the containers are not started, see [Why aren't the expected docker containers running on my edge node?](Troubleshooting.md#why-arent-the-expected-docker-containers-running-on-my-edge-node) to troubleshoot.

* Once the containers are all running, you should be able to verify that your workload is sending short messages like `{"cpu":1.49}` to the Watson IoT Platform in one of two ways:
  * Return to the Watson IoT Platform web pages: select Devices in the left panel, select your Gateway instance (it should have a blue dot next to it meaning it is "connected"), and click on the Recent Events. The default publish interval for the CPU example is 30 seconds, so you may have to wait that long before seeing the first message.
  * Or, use the mosquitto_sub command previously shown in this guide:
```
mosquitto_sub -v -h $HZN_ORG_ID.messaging.$WIOTP_DOMAIN -p 8883 -i "${WIOTP_CLIENT_ID_APP:0:38}" -u "$WIOTP_API_KEY" -P "$WIOTP_API_TOKEN" --capath /etc/ssl/certs -t iot-2/type/$WIOTP_GW_TYPE/id/$WIOTP_GW_ID/evt/status/fmt/json
```
* If you don't see MQTT messages coming from your edge node, look at the cpu2wiotp log messages for errors:
```
grep '_cpu2wiotp\[' /var/log/syslog
```

* You can also manually send a message from your edge node via the local IoT-Core edge connector (first start the mosquitto_sub command above in another shell). This simulates how a workload sends data to WIoTP via the edge connector:

```
mosquitto_pub -h localhost -p 8883 -i "a:myapp" --cafile /var/wiotp-edge/persist/dc/ca/ca.pem -t iot-2/evt/status/fmt/json -m '{"message": "Hello from the edge"}' -d
```

## Configuring other nodes to run this deployment Pattern

To use this deployment pattern on other edge nodes you don't have to repeat everything in this document. The summary of what to do on each additional edge node is:
1. Create the gateway in the WIoTP cloud UI (with the same gateway type)
1. Install a debian-based Linux operating system and docker
1. Install the Horizon packages
1. Register the edge node and to use this pattern by running `wiotp_agent_setup`

## What To Do If Things Go Wrong
Take a look at the [Edge Troubleshooting Guide](Troubleshooting.md).

You may also wish to explore the `hzn` command, a powerful tool for debugging this system:
* Online help is available within hzn:
```
hzn --help
```
* Help is also available for all of the sub commands in hzn, simply add "--help" after any command to get details:
```
hzn exchange pattern insertworkload --help
```

Docker commands are also useful for observing the operation of the system on edge nodes.  You can use the `docker ps` command to find the identifiers of the running containers then use `docker inspect ...` to get detailed information, such as the docker virtual private networks their virtual interfaces are connected to, their IP addresses on those networks, etc.

The Edge system code relies heavily on REST APIs and thus the `curl` command can also be a powerful tool for debugging.  In fact most of the functionality of the hzn command is implemented by invoking REST APIs from the Horizon Exchange in the cloud, or the local Horizon Agent on the Edge Node.  Often you can run an hzn command with the `-v` (verbose) argument to see the REST API methods being used under the covers for the command, and their results.  You can then use this information to directly interact with the APIs using curl:
* Example "verbose" output from hzn:
```
[verbose] GET http://localhost/status
[verbose] HTTP code: 200
[verbose] GET https://...ibmcloud.com/api/v0002/edge/orgs/.../patterns
[verbose] HTTP code: 403
```
You are encouraged to join the discussion on the Horizon forum at the link below where you may find answers to frequently asked questions, and ask questions of your own:
* [http://discourse.bluehorizon.network/](http://discourse.bluehorizon.network/)


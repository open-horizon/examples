## Pi3-Streamer Microservice
This defines the microservice for a LAN webcam using Raspberry Pi3 and a Pi Camera.  
Originally packaged in docker as [cogwerx-mjpg-streamer-pi3](https://github.com/open-horizon/cogwerx-mjpg-streamer-pi3)

## Setup Steps

### Manual Pre-Setup Steps: 
[Download](https://bluehorizon.network/documentation/disclaimer) a Raspbian image for your Pi 3 (we tested this using [Horizon](https://bluehorizon.network/)'s Raspbian image). Unzip and flash the image to your micro SD Card, (setup WiFi) and boot. Full setup instructions for that can be found [here](https://bluehorizon.network/documentation/adding-your-device).
Run raspi-config as root and set GPU memory and enable the Pi Cam:

    raspi-config

### Set the following options:
* Option 5 (Connections to peripherals): P1 (Camera) Enable the Pi Camera  
* Option 7 (Advanced Options): A3 (Memory Split): Set GPU memory to 256 MB  
Reboot  

&nbsp; &nbsp; &nbsp; <img src="https://user-images.githubusercontent.com/16260619/37161848-a253e6be-22a8-11e8-9e1b-73509ae8c4dd.png" width="480" />

You're done with pre-setup steps.

## Automatic Deployment on IBM Edge with Watson IoT Platform
Follow the steps in this [Edge Quick Start Guide](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md), up until "Define an Additional Microservice and Workload in the Horizon Exchange".  
At that point: stop, return here and continue with this guide.  
The Quick Start Guide will instruct you to define a device name and a device type. As an example, your information may look something like:  

    Device Type: arm32-PI3STRMR    (a general name for all devices of this type)  
    Device Name: PI3-Home          (a specific name for this device)  
    Device Token: jkdas9dusadkna   (some secure string, specific to this device)  
    API Key: 'generated-chars'  
    API Token: 'generated-chars'  

These values aren't visible outside of your IBM Cloud organization. The token is not retrievable after definition.  API keys may be used for all devices you define, or per device at your discretion.

### Start Using IBM Edge to Define and Deploy your Pi 3 LAN Streamer
At this point, you could register your edge node with Horizon and have the default WIoTP core-iot service deployed to it. Some additional definition is needed to deploy the Pi3 Streamer microservice and workload to your edge node.  

First, clone the openhorizon examples project which contains files that you will need during the following steps:

    cd ~
    git clone https://github.com/open-horizon/examples.git

### Signing Keys
We'll generate a signing key for this Pi to use in defining microservices that will be authorized to run on your devices.  This key will be used to sign the deployment definitions, and to verify the microservices when they begin to run on the Pi. This can take a few minutes to generate on a Pi 3.

 * Generate a signing key for horizon to use in publishing microservices and workloads. This can take a few minutes on the Pi. Once generated, import your key into horizon with `hzn key import`. Verify with `hzn key list`.
```bash
mkdir ~/keys && cd ~/keys
hzn key create <x509 org> <x509 cn>   # example: hzn key create ibm thomas@ibm.com
hzn key import --public-key-file=<key file name>
hzn key list   # You should see your key listed in the output
```

Your key should show in the output list, similar to the following:  
```
root@horizon-0000000079b68342:~# hzn key list
[
  {
    "id": "thomas-2e434f1456233d537a7122763884af1cba3e77-public.pem",
    "common_name": "thomas@us.ibm.com",
    "organization_name": "ibm",
    "serial_number": "2e:43:4f:14:56:23:3d:53:7a:71:22:76:38:84:af:1c:ba:3e:77",
    "not_valid_before": "2018-04-01 07:35:41 +0000 UTC",
    "not_valid_after": "2022-04-01 19:19:42 +0000 UTC"
  }
]
```

### Microservice and Workload Setup / Registration
At this point, you could register your edge node with IBM Edge with WIoTP and have the default WIoTP core-iot service deployed to it. Now we'll also define the Pi3-Streamer microservice and workload in your WIoTP org, such that registration of your device will cause your Edge to pull those containers, run them, and publish status to your Watson IoT Platform org.

* Clone the open-horizon examples project which contains files that you will need during the following steps:
```bash
cd ~
git clone https://github.com/open-horizon/examples.git 
```
* Temporarily set the Horizon Exchange URL
```bash
export HZN_EXCHANGE_URL="https://$HZN_ORG_ID.$WIOTP_DOMAIN/api/v0002/edgenode/"
```

#### Microservice  
You'll define the microservice in your WIoTP organization using the files already in the repo, plus your own specific credentials and config and then publish the definition to the Exchange. You'll reference your Docker hub account to do this.

First, set environment variables for your microservice. 
```bash
cd ~/examples/edge/services/pi3_streamer
cp horizon/envvars.sh.sample  horizon/envvars.sh
vim horizon/envvars.sh.sample  # or use your favorite text editor
```
Change the `HZN_ORG_ID` to your own WIoTP organization; provide your Docker Hub ID, and a name for your domain. (You can use a fictitious one if you like.)  Save the file and export the environment var's with the `source` command.
```bash
source horizon/envvars.sh
```

Next, list the microservices already in your org. Then take a look at the files in the directory.  You'll build your version of the microservice using `make`, you'll add the "pi3streamer" microservice to your WIoTP organization, and push the docker image up to your Docker Hub registry, and verify that the microservice was added to the exchange.  
```bash
hzn exchange microservice list | jq .   # Your microservice won't appear yet
make build                              # This will build your pi3streamer Docker container image
hzn dev microservice verify             # This will verify the definition in horizon/userinput.json and horizon/microservice.definition.json before publishing it to the exchange
hzn dev microservice publish -k $PRIVATE_KEY_FILE       # This will publish the ms definition to the exchange, and push your Docker image to your registry
hzn exchange microservice list | jq .   # Your microservice should now be listed in the exchange
```

Your microservice definition in the Exchange may look like the following:
```bash
root@horizon-0000000079b68342:~/examples/edge/services/pi3_streamer# hzn exchange microservice list
[
  "5fdjke/mydomain.net-microservices-pi3streamer_1.0.0_arm"
]
```

#### Workload  
The sole workload associated with the pi3streamer is in `examples/edge/wiotp/pi3streamer2wiotp`. Setting up the workload is similar to the previous microservice step. The Pi3 Streamer2WIoTP workload will run in its own Docker container and do the following:
* It will query the pi3streamer microservice's HTTP REST API via `curl` and inspect the output to determine the pi3streamer is up and running
* It will send a status message to WIoTP every 10 seconds (you can set the value specifically if you like)

First, set environment variables for your workload. You'll use the Device Type, Device ID, and Device Token credentials that you created in Watson IoT Platform for your Pi3.
```bash
cd ~/examples/edge/wiotp/pi3streamer2wiotp
cp horizon/envvars.sh.sample  horizon/envvars.sh
vim horizon/envvars.sh.sample  # or use your favorite text editor
```
Change the `HZN_ORG_ID` to your own WIoTP organization; provide your Docker Hub ID, and a name for your domain. (You can use a fictitious one if you like.)  Also provide your Device-specific "WIOTP_*" credentials that you created in Watson IoT Platform.
Save the file and export the environment var's with the `source` command.
```bash
source horizon/envvars.sh
```

Next, list the workloads already in your org. Then take a look at the files in the directory. You'll build your version of the workload using `make`, you'll add the "pi3streamer2wiotp" workload to your WIoTP organization, and push the docker image up to your Docker Hub registry, and verify that the workload was added to the exchange.  

```bash
hzn exchange workload list | jq .       # Your workload won't appear yet
make build                              # This will build your pi3streamer2wiotp Docker container image
hzn dev dependency fetch -p ~/examples/edge/services/pi3_streamer/horizon/  # This will define this workload as dependent on the pi3streamer microservice  (See our Developer Guide for details)
hzn dev dependency fetch -s https://internetofthings.ibmcloud.com/wiotp-edge/microservices/edge-core-iot-microservice --ver 2.4.0 -o IBM -a $ARCH -k /etc/horizon/trust/publicWIoTPEdgeComponentsKey.pem  # This will define an additional dependency on IBM's Edge Core IoT microservices, which provide MQTT messaging and container management
hzn dev workload verify                 # This will verify the definition in horizon/userinput.json and horizon/workload.definition.json before publishing it to the exchange
hzn dev workload publish -k $PRIVATE_KEY_FILE   # This will publish the wl definition to the exchange, and push your Docker image to your registry
hzn exchange workload list | jq .       # Your workload should now be listed in the exchange
```

#### Pattern
Next, add your workload definition to the Pattern. The Pattern defines the list of multiple containers that a WIoTP device type will run, for one or more architectures (x86, ARM, etc).  By default, every Pattern first requires IBM's Core IoT workload. You'll add your own workload here.

```bash
hzn exchange pattern insertworkload -k $PRIVATE_KEY_FILE -f pattern/insert-pi3streamer2wiotp.json $WIOTP_GW_TYPE  # soon you can use -K $PUBLIC_KEY_FILE and then will not have to import it
hzn exchange pattern list $WIOTP_GW_TYPE | jq .   # Your workload should be listed in the pattern
```

Now you're set to register your Pi 3 as a LAN Streamer

## Registration
To register your Pi 3 to run your workload, you'll provide some config details (Device Type, ID, and Token) and instruct the local agent to set up your Pi 3 as an IBM edge. 

Register the node and start the Watson IoT Platform core-IoT service and the CPU workload:
```bash
wiotp_agent_setup --org $HZN_ORG_ID --deviceType $WIOTP_GW_TYPE --deviceId $WIOTP_GW_ID --deviceToken "$WIOTP_GW_TOKEN"
```

At this point, your Pi 3 will contact the IBM Edge Exchange, establish your Device identity using your credentials, and make an Agreement to run your LAN streamer containers. The IBM Edge containers will come down to your Pi 3, and the copies you've defined locally will begin to run. 

Use the command `hzn agreement list` to view the agreements your Pi 3 has made with the IBM Exchange.  You should see something like this:

```
root@horizon-0000000079b68342:~/examples/edge/wiotp/pi3streamer2wiotp# hzn agreement list
[
  {
    "name": "Policy for edge-core-iot-microservice merged with ARM32-pi3streamer_internetofthings.ibmcloud.com-wiotp-edge-workloads-edge-core-iot-workload_IBM_arm",
    "current_agreement_id": "1a7cff85059b080e35662c7d728c62870fab9d95d58c6947b9c5e3be36b34ac8",
    "consumer_id": "IBM/wiotp-agbot-1",
    "agreement_creation_time": "2018-04-01 06:31:52 +0000 UTC",
    "agreement_accepted_time": "2018-04-01 06:32:08 +0000 UTC",
    "agreement_finalized_time": "2018-04-01 06:32:10 +0000 UTC",
    "agreement_execution_start_time": "",
    "agreement_data_received_time": "",
    "agreement_protocol": "Basic",
    "workload_to_run": {
      "url": "https://internetofthings.ibmcloud.com/wiotp-edge/workloads/edge-core-iot-workload",
      "org": "IBM",
      "version": "2.4.0",
      "arch": "arm"
    }
  },
  {
    "name": "Policy for pi3streamer merged with Policy for edge-core-iot-microservice merged with ARM32-pi3streamer_mydomain.net-workloads-pi3streamer2wiotp_5fdjke_arm",
    "current_agreement_id": "baf923bdc6ad17c1cea17b6a350c73c959e813c207842752f1679186df4870a6",
    "consumer_id": "IBM/wiotp-agbot-1",
    "agreement_creation_time": "2018-04-01 06:31:54 +0000 UTC",
    "agreement_accepted_time": "2018-04-01 06:32:12 +0000 UTC",
    "agreement_finalized_time": "2018-04-01 06:32:18 +0000 UTC",
    "agreement_execution_start_time": "",
    "agreement_data_received_time": "",
    "agreement_protocol": "Basic",
    "workload_to_run": {
      "url": "https://mydomain.net/workloads/pi3streamer2wiotp",
      "org": "5fdjke",
      "version": "1.0.0",
      "arch": "arm"
    }
  }
]

```

Use the commands `docker ps -a` and `docker stats` to view Docker's execution of your LAN streamer containers. 

After a few minutes, you'll see your docker containers running and listed like this:

```
root@horizon-0000000079b68342:~/examples/edge/wiotp/pi3streamer2wiotp# docker ps -a
CONTAINER ID        IMAGE                                                   COMMAND                  CREATED             STATUS              PORTS                                            NAMES
2fdf2f3b8fb0        openhorizon/arm_pi3streamer2wiotp                       "/bin/sh -c /workloa…"   6 minutes ago       Up 6 minutes                                                         baf923bdc6ad17c1cea17b6a350c73c959e813c207842752f1679186df4870a6-pi3streamer2wiotp
91cdf9fd1963        openhorizon/arm_pi3streamer                             "/usr/bin/entry.sh .…"   6 minutes ago       Up 6 minutes        0.0.0.0:8080->8080/tcp                           mydomain.net-microservices-pi3streamer_1.0.0_a8d093fe-4266-4ee3-aaf1-5017ad09d3a7-pi3streamer
2277a85a6d9c        wiotp-connect/edge/armhf/edge-core-iot-workload:1.0.3   "/start.sh"              6 minutes ago       Up 6 minutes                                                         1a7cff85059b080e35662c7d728c62870fab9d95d58c6947b9c5e3be36b34ac8-edge-core-iot-workload
d7da1e359a6f        wiotp-connect/edge/armhf/edge-mqttbroker:1.1.3          "/start.sh"              7 minutes ago       Up 7 minutes                                                         internetofthings.ibmcloud.com-wiotp-edge-microservices-edge-core-iot-microservice_2.4.0_08f35ede-ce88-4982-aab6-a8ac8b100333-edge-mqttbroker
61ac2c67ac24        wiotp-infomgmt/edge/armhf/edge-im:1.0.15                "/start.sh"              7 minutes ago       Up 7 minutes                                                         internetofthings.ibmcloud.com-wiotp-edge-microservices-edge-core-iot-microservice_2.4.0_08f35ede-ce88-4982-aab6-a8ac8b100333-edge-im
be9b92ff4433        wiotp-connect/edge/armhf/edge-connector:2.4.1           "/start.sh"              7 minutes ago       Up 7 minutes        0.0.0.0:1883->1883/tcp, 0.0.0.0:8883->8883/tcp   internetofthings.ibmcloud.com-wiotp-edge-microservices-edge-core-iot-microservice_2.4.0_08f35ede-ce88-4982-aab6-a8ac8b100333-edge-connector

```

## Your Pi 3 is a LAN Video Streamer

Using a web browser, visit your Pi3's IP address followed by 8080 (e.g. http://xxx.xxx.xxx.xxx:8080) on your LAN.
That's it! You should be able to see a simple web page with a static image from your Pi.  Connect to http://xxx.xxx.xxx.xxx:8080/?action=stream to see your video stream.

&nbsp; &nbsp; &nbsp; <img src="https://user-images.githubusercontent.com/16260619/37161339-3ccba3aa-22a7-11e8-8938-516ce59d5f2d.png" width="640" />

-------------------
### Steps:
1. Setup your hardware.  
&nbsp;&nbsp; See ["Initial Setup"](https://github.com/open-horizon/cogwerx-mjpg-streamer-pi3/blob/master/README.md)
2. Setup your IBM Cloud account and Watson IoT Platform org.  
&nbsp;&nbsp; See [Setup Your Organization in the Watson IoT Platform](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md#setup-your-organization-in-the-watson-iot-platform) in the Edge Quick Start Guide
3. Prepare your Pi3.  
&nbsp;&nbsp; See ["Prepare Your Edge Node"](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md#prepare-your-edge-node) in the Edge Quick Start Guide
4. Define and publish this microservice to your org using the `.json` template files in `./horizon`  
&nbsp;&nbsp; See the CPU example in [Define an Additional Microservice...](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md#define-an-additional-microservice-and-workload-in-the-horizon-exchange) in the Edge Quick Start Guide

Then move on to define the workload and pattern, using the code in `examples/edge/wiotp/pi3streamer2wiotp` and register your Pi3 to run the pattern. 

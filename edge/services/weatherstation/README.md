## Personal Weather Station Microservice (PWSMS)

### Repo contents: Docker container build scripts and supporting files  
- Makefile: executes container build, dev, run, publish steps
- Dockerfile.<ARCH>:  Docker container image build files for various architectures (amd64, RPi (armhf))
- horizon/: Files for definition of microservice
- weewx/: Files for weewx, an open source python-based personal weather station (PWS) utility

### Setup Steps  
#### Hardware Setup  
- Setup your hardware:  
  - Linux desktop (amd64) or RPi 2/3 (armhf)  
  - Supported PWS such as one of [these](https://bluehorizon.network/documentation/weather)  
  - Plug your PWS base station into your linux box via USB  
- Register your PWS at [wunderground.com](https://www.wunderground.com/personal-weather-station/mypws)  
  - Note your `Station ID` and `Station Key` (you'll record these for device registration)  
<img width="568" alt="screen shot 2018-03-15 at 1 22 36 pm" src="https://user-images.githubusercontent.com/16260619/37489250-3c743c48-2854-11e8-8925-79d94d7f7517.png">
  
- Install a linux distribution (ubuntu 16.04+ recommended)  
  - For RPi2/3: [Download](https://www.raspberrypi.org/downloads/raspbian/) a Raspbian image for your Pi 3 (we tested this on a Pi 3 using Raspbian Stretch). Unzip and flash the image to your micro SD Card, boot, [setup WiFi](https://www.raspberrypi.org/documentation/configuration/wireless/wireless-cli.md), enable [SSH](https://www.raspberrypi.org/documentation/remote-access/ssh/).
  - For Desktop: Download and install ubuntu 16.04 on your x86 linux deskop machine.

#### IBM Watson IoT Platform / Horizon Setup  
  - Follow the setup steps in the open-horizon [Edge Quick Start Guide](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md)
  - Follow the Watson IoT Platform Setup step in this [Edge Quick Start Guide](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md#setup-your-organization-in-the-watson-iot-platform). 
You will define a device name and a device type. As an example, your information may look something like:  

    Device Type: arm32-pws         (a general name for all devices of this type)  
    Device Name: PI3-PWS           (a specific name for this device)  
    Device Token: jkdas9dusadkna   (some secure string, specific to this device)  
    API Key: 'generated-chars'  
    API Token: 'generated-chars'  
  - [Prepare Your Edge Node](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md#prepare-your-edge-node)  (install horizon packages and prereqs)  
  - [Verify Gateway Credentials and Access](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md#verify-your-gateway-credentials-and-access)
 
### PWS Microservice and Workload Setup / Registration
At this point, you could register your edge node with Horizon and have the default WIoTP core-iot service deployed to it. Now we'll also define the PWS microservice and workload in your WIoTP org, such that registration of your device will cause your Edge to pull those containers, run them, and publish status to your Watson IoT Platform org.

### Signing Keys
We'll generate a signing key for this Pi to use in defining microservices that will be authorized to run on your devices.  This key will be used to sign the deployment definitions, and to verify the microservices when they begin to run. (This can take a few minutes to generate on a Pi 3.)

 * Generate a signing key for horizon to use in publishing microservices and workloads. Once generated, import your key into horizon with `hzn key import`. Verify with `hzn key list`.
```bash
mkdir -p ~/hzn/keys && cd ~/hzn/keys
hzn key create <x509 org> <x509 cn>   # example: hzn key create ibm thomas@ibm.com
export PRIVATE_KEY_FILE=$(ls ~/hzn/keys/*-private.key)
export PUBLIC_KEY_FILE=$(ls ~/hzn/keys/*-public.pem)
hzn key import --public-key-file=$PUBLIC_KEY_FILE
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


* Clone the openhorizon examples project which contains files that you will need during the following steps:
```bash
cd ~
git clone https://github.com/open-horizon/examples.git 
```
* Temporarily set the Horizon Exchange URL
```bash
export HZN_EXCHANGE_URL="https://$HZN_ORG_ID.$WIOTP_DOMAIN/api/v0002/edgenode/"
```  

First, set environment variables for your microservice. 
```bash
cd ~/examples/edge/services/weatherstation
cp horizon/envvars.sh.sample  horizon/envvars.sh
vim horizon/envvars.sh  # or use your favorite text editor
```

This is a view of that file:
```
# Set this to the organization you created in the Watson IoT Platform
export HZN_ORG_ID=myorg

export ARCH=arm   # arch of your edge node: amd64, or arm for Raspberry Pi, or arm64 for ODROIDC2 / Jetson TX2
export PWSMS_NAME=pwsms   # the name of the microservice, used in the docker image path and in the microservice url
export PWSMS_VERSION=1.1.0   # the microservice version, and also used as the tag for the docker image. Must be in OSGI version format.

export DOCKER_HUB_ID=mydockerhubid   # your docker hub username, sign up at https://hub.docker.com/sso/start/?next=https://hub.docker.com/
export MYDOMAIN=mydomain.com    # used in the microservice url

# There is normally no need for you to edit these variables
export WIOTP_DOMAIN=internetofthings.ibmcloud.com
export HZN_EXCHANGE_URL="https://$HZN_ORG_ID.$WIOTP_DOMAIN/api/v0002/edgenode/"
```

Change the `HZN_ORG_ID` to your own WIoTP organization; provide your `DOCKER_HUB_ID`, and a name for `MYDOMAIN`. (You can use a fictitious one if you like.)  Save the file and export the environment var's with the `source` command.
```bash
source horizon/envvars.sh
```
Next, list the microservices already in your org. Then take a look at the files in the directory.  You'll build your version of the microservice using `make`, you'll add the "pwsms" microservice to your WIoTP organization, and push the docker image up to your Docker Hub registry, and verify that the microservice was added to the exchange.  
```bash
hzn exchange microservice list | jq .   # Your microservice won't appear yet
make build                              # This will build your pi3streamer Docker container image
hzn dev microservice verify             # This will verify the definition in horizon/userinput.json and horizon/microservice.definition.json before publishing it to the exchange
docker login                            # Login to Docker Hub with your name/pwd prior to publishing your container image
hzn dev microservice publish -k $PRIVATE_KEY_FILE       # This will publish the ms definition to the exchange, and push your Docker image to your registry
hzn exchange microservice list | jq .   # Your microservice should now be listed in the exchange
```

Your microservice definition in the Exchange may look like the following:
```bash
root@horizon-0000000079b68342:~/examples/edge/services/weatherstation# hzn exchange microservice list
[
  "5fdjke/mydomain.net-microservices-pwsms_1.0.0_arm"
]
```

#### Workload  
The sole workload associated with the pwsms is in `examples/edge/wiotp/pws2wiotp`. Setting up the workload is similar to the previous microservice step. The PWS2WIoTP workload will run in its own Docker container and do the following:
* It will query the pwsms microservice's HTTP REST API via `curl` and inspect the output to determine the pwsms microservice is up and running
* It will send a status message with weather data to WIoTP every 10 seconds (you can set the value specifically if you like)

First, set environment variables for your workload. You'll use the Device Type, Device ID, and Device Token credentials that you created in Watson IoT Platform for your Edge node.
```bash
cd ~/examples/edge/wiotp/pws2wiotp
cp horizon/envvars.sh.sample  horizon/envvars.sh
vim horizon/envvars.sh  # or use your favorite text editor
```
Change the `HZN_ORG_ID` to your own WIoTP organization; If you haven't already, provide your Docker Hub ID, and a name for your domain. (You can use a fictitious one if you like.)  Also provide your Device-specific "WIOTP" credentials that you created in Watson IoT Platform, if you haven't already done so.
Save the file and export the environment var's with the `source` command.
```bash
source horizon/envvars.sh
```

Next, list the workloads already in your org. Then take a look at the files in the directory. You'll build your version of the workload using `make`, you'll add the "pi3streamer2wiotp" workload to your WIoTP organization, and push the docker image up to your Docker Hub registry, and verify that the workload was added to the exchange.  

```bash
hzn exchange workload list | jq .       # Your workload won't appear yet
make build                              # This will build your pws2wiotp Docker container image
hzn dev dependency fetch -p ~/examples/edge/services/weatherstation/horizon/  # This will define this workload as dependent on the pwsms microservice  (See our Developer Guide for details)
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

Now you're set to register your Edge node as a PWS data producer.

#### Set your PWS Device Settings
To register your Pi3 to run the pattern with your PWS config values (Station ID, Station Key, etc), you'll need to provide those values in horizon's template file, located in `/etc/wiotp-edge/`.  Backup the existing file and use the template file to fill in your values. 

```bash
mv /etc/wiotp-edge/hznEdgeCoreIoTInput.json.template /etc/wiotp-edge/hznEdgeCoreIoTInput.json.template.orig # Backup the original
envsubst < hznEdgeCoreIoTInput.json.template > /etc/wiotp-edge/hznEdgeCoreIoTInput.json.template
```  
Finally, edit your `/etc/wiotp-edge/hznEdgeCoreIoTInput.json.template` file and provide your chosen values:

```bash
...
      "variables": {
        "PWS_ST_TYPE": "FineOffsetUSB",  (<-- replace with your station driver++)
        "PWS_MODEL": "WS2080",           (<-- replace with your station type++)
        "PWS_WU_ID": "KCAENCIN70",       (<-- replace with your WU Station ID)
        "PWS_WU_KEY": "7HGR6HD3",        (<-- replace with your WU Station Key)
        "PWS_WU_RPDF": "False"           (True/False to use "rapidfire mode", sends data to WU more often)
      }
...
```  
++Weewx Hardware guide: http://www.weewx.com/docs/hardware.htm

#### Registration
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
* Verify that 2 agreements are made, one for `edge-core-iot-workload` and one for `pws2wiotp`.  The output should indicate an `agreement_finalized_time`, and eventually an `agreement_execution_start_time` should also be populated.
```
hzn agreement list | jq . 
```

* Eventually the 6 docker containers should be running: 4 for the core-IoT service and 2 for the PWS example.  Verify this with:
```
watch -n 1 docker ps
```

After a minute or so (depending on device architecture, internet connection), you should see data at your PWS page on Weather Underground:
<img width="1127" alt="screen shot 2018-03-15 at 2 48 49 pm" src="https://user-images.githubusercontent.com/16260619/37492856-2326335c-2860-11e8-9248-1a50dba0bca4.png">

and in WIoTP, under your device ID:  

<img width="747" alt="screen shot 2018-03-15 at 2 54 31 pm" src="https://user-images.githubusercontent.com/16260619/37493057-dbdbfbde-2860-11e8-8b94-7454e7bb7475.png">

### References  
* pywws: http://pywws.readthedocs.org/en/latest/index.html
* weewx: http://www.weewx.com/  (open source, many supported PWS's)
* Weather Undergdound PWS Info: https://www.wunderground.com/weatherstation/overview.asp

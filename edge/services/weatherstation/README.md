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
  - For RPi2/3: Download a raspbian image for your Pi (we tested this on a Pi3 using [Horizon]()'s raspbian image). Unzip and flash the image to your micro SD Card, (setup WiFi) and boot.  
  - For Desktop: Download and install ubuntu 16.04  

#### IBM Watson IoT Platform / Horizon Setup  
- Follow the setup steps in the open-horizon [Edge Quick Start Guide](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md)
  - [Setup Your Organization in the Watson IoT Platform](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md#setup-your-organization-in-the-watson-iot-platform)  
  - [Prepare Your Edge Node](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md#prepare-your-edge-node)  (install horizon packages and prereqs)  
  - [Verify Gateway Credentials and Access](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md#verify-your-gateway-credentials-and-access)
 
### PWS Microservice and Workload Setup / Registration
At this point, you could register your edge node with Horizon and have the default WIoTP core-iot service deployed to it. Now we'll also define the Pi3-Streamer microservice and workload in your WIoTP org, such that registration of your device will cause your Edge to pull those containers, run them, and publish status to your Watson IoT Platform org.

* Generate a signing key for horizon to use in publishing microservices and workloads. This can take a few minutes on the Pi. Once generated, import your key into horizon with `hzn key import`. Verify with `hzn key list`.
```bash
mkdir ~/keys && cd ~/keys
hzn key create <x509 org> <x509 cn>   # example: hzn key create ibm chris@ibm.com
hzn key import --public-key-file=<key file name>
hzn key list
```

* Clone the openhorizon examples project which contains files that you will need during the following steps:
```bash
cd ~
git clone -b pws2wiotp https://github.com/open-horizon/examples.git 
```
* Temporarily set the Horizon Exchange URL
```bash
export HZN_EXCHANGE_URL="https://$HZN_ORG_ID.$WIOTP_DOMAIN/api/v0002/edgenode/"
```  

* Export two additional environment variables for publishing your microservice and workload:
```bash
export PWSMS_VERSION=1.1.0
export PWS2WIOTP_VERSION=1.1.0
export DOCKER_HUB_ID=openhorizon
```

#### Microservice  
* List the microservices already in your org. Add the "pwsms" microservice to your WIoTP organization and see that it was added:  
```bash
hzn exchange microservice list | jq .
cd ~/examples/edge/services/weatherstation/horizon
envsubst < pwsms-template.json > ms_def.json
hzn exchange microservice publish -f ms_def.json -k <your_private_key_file>
hzn exchange microservice list | jq .
```

#### Workload  
* Configure the PWS2WioTP usage workload definition file using your environment variables, add it to your WIoTP organization, and see that it was added:  
```bash
hzn exchange workload list | jq .
cd ~/examples/edge/wiotp/pws2wiotp/horizon
envsubst < pws2wiotp-template.json > wl-def.json
hzn exchange workload publish -f wl_def.json -k <your_private_key_file>
hzn exchange workload list | jq .
```

#### Augment the Edge Node Deployment Pattern
The Edge system deploys Patterns of code onto WIoTP Edge Node gateways. The deployment Pattern used for a particular gateway has the same name as its Gateway Type. By default the deployment pattern includes the WIoTP core-IoT service. Here you will update the deployment Pattern for your Gateway Type so that it also includes the workload that we just added to your platform.

* Configure the PWS usage pattern json file using your environment variables and add it to your pattern:
```bash
hzn exchange pattern list $WIOTP_GW_TYPE | jq .
cd ~/examples/edge/wiotp/pws2wiotp/horizon/pattern/
envsubst < insert-pws2wiotp-template.json > ~/pi_def.json
hzn exchange pattern insertworkload -f pi_def.json $WIOTP_GW_TYPE -k <your_private_key_file>
```
* Verify that the pws2wiotp Workload was inserted into the Pattern for your Gateway Type:
```bash
hzn exchange pattern list $WIOTP_GW_TYPE | jq .
```
* [Optional] Unset the HZN_EXCHANGE_URL environment variable (because after registration in the next section `hzn` can get the value from the Horizon agent):
```bash
unset HZN_EXCHANGE_URL
```
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
        "PWS_ST_TYPE": "FineOffsetUSB   (<-- replace with your station driver++)",
        "PWS_MODEL": "WS2080            (<-- replace with your station type++)",
        "PWS_WU_ID": "KCAENCIN70        (<-- replace with your WU Station ID)",
        "PWS_WU_KEY": "7HGR6HD3         (<-- replace with your WU Station Key)",
        "PWS_WU_RPDF": "False           (True/False to use "rapidfire mode", sends data to WU more often)"
      }
...
```

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

* Eventually the 6 docker containers should be running: 4 for the core-IoT service and 2 for the Pi3streamer example.  Verify this with:
```
watch -n 1 docker ps
```

After a minute or so (depending on device architecture, internet connection), you should see data at your PWS page on Weather Underground:
<img width="1127" alt="screen shot 2018-03-15 at 2 48 49 pm" src="https://user-images.githubusercontent.com/16260619/37492856-2326335c-2860-11e8-9248-1a50dba0bca4.png">

and in WIoTP, under your device ID:
<img width="1114" alt="screen shot 2018-03-15 at 2 51 46 pm" src="https://user-images.githubusercontent.com/16260619/37492939-74b79f8a-2860-11e8-8f91-a4c6199cf157.png">

### References  
* pywws: http://pywws.readthedocs.org/en/latest/index.html
* weewx: http://www.weewx.com/  (open source, many supported PWS's)
  * ++Weewx Hardware guide: http://www.weewx.com/docs/hardware.htm
* Weather Undergdound PWS Info: https://www.wunderground.com/weatherstation/overview.asp

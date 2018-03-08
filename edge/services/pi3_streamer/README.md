## Pi3-Streamer Microservice
This defines the microservice for a LAN webcam using Raspberry Pi3 and a Pi Camera.  
Originally packaged in docker as [cogwerx-mjpg-streamer-pi3](https://github.com/open-horizon/cogwerx-mjpg-streamer-pi3)

### Steps:
1. Setup your hardware.  
&nbsp;&nbsp; See ["Initial Setup"](https://github.com/open-horizon/cogwerx-mjpg-streamer-pi3/blob/master/README.md)
2. Setup your IBM Cloud account and Watson IoT Platform org.  
&nbsp;&nbsp; See [Setup Your Organization in the Watson IoT Platform](https://github.com/open-horizon/examples/wiki/Edge-Quick-Start-Guide#setup-your-organization-in-the-watson-iot-platform) in the Edge Quick Start Guide
3. Prepare your Pi3.  
&nbsp;&nbsp; See ["Prepare Your Edge Node"](https://github.com/open-horizon/examples/wiki/Edge-Quick-Start-Guide#prepare-your-edge-node) in the Edge Quick Start Guide
4. Define and publish this microservice to your org using the `.json` template files in `./horizon`  
&nbsp;&nbsp; See the CPU example in [Define an Additional Microservice...](https://github.com/open-horizon/examples/wiki/Edge-Quick-Start-Guide#define-an-additional-microservice-and-workload-in-the-horizon-exchange) in the Edge Quick Start Guide

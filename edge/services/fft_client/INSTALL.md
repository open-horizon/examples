# FFT Example services

## What is FFT?

The Fast Fourier Transform algorithm is a set of equations that identify anomalies or patterns in cyclic patterns. 
It is commonly used for vibration analysis to determine when machinery is operating outside standard parameters. 

## What does this example do?

This example is composed of three micro-services that work together: `volantmq`, `fft_server`, and `fft_client`. 
The `volantmq` service is a [message queue broker](https://godoc.org/github.com/VolantMQ/volantmq) conforming to the MQTT 3.1 specification. 
The `fft_server` service analyzes sample audio snippets sorted into bins using the FFT algorithm and will post a message to a topic if it finds any signal that falls outside a pre-specified range. 
The optional `fft_client` uses a pre-configured audio input device to record audio snippets of specified duration and posts them to a topic for the `fft_server` service to analyze.
All three services are multi-platform, multi-architecture and written in Golang. 
This example demonstrates how to wrap a library (`portaudio`) into a lightweight micro-service, and the pattern of creating a set of services where the server is required and the client is optional. 
It is believed that you can substitute any MQTT 3.1 or 3.1.1-compliant service for volantmq; however, this has not been confirmed through testing. 

## Pre-requisites

### Infrastructure

This example has been tested on the arm64 (x86_64) and armhf (arm32) architectures, and provides build targets for both. 
You may run it on bare metal or in a VM with Ubuntu 18.04 (bionic).  
It will also run on a Raspberry Pi 3 or 4 with Raspbian buster.  
On the RPi, it has been confirmed to work with 1GB or 2GB RAM and minimal storage. 
On bare metal or in a VM, it functions fine with 1 vCPU, 4GG RAM, and 20GB storage although it will require and use far less.
You will also want to ensure that you have a functioning microphone (stereo or mono) attached to the device.

### Platform

This example assumes that you have Open Horizon Management Hub services installed and running, so you can publish service definitions and a pattern to that Exchange.  
It also requires the Agent to be installed, configured, and authorized to connect to that Exchange (confirmed with `hzn exchange user list`).
It will also work with the all-in-one pattern where both Management Hub and Agent are installed on the same edge node. 
It has not been tested with deployment to a cluster. 

## Installation

### Step 1: Clone the Repository

``` shell
git clone https://github.com/open-horizon/examples.git
cd examples
```

### Step 2: Editing the Services and Pattern Definition Files

TODO: Templatize the files and use env vars and Makefile to automatically generate

Change to the folder with the message broker service.  We'll start by configuring that one. 

``` shell
cd edge/services/volantmq
```

NOTE: The client and the server will be connecting to the message broker with separate authenticated credentials. 
A default password for each is currently hard-coded into the configuration files.  
If you wish to change that, you will need to use a different set of passwords and you will need to put the encoded version in various configuration files. 
Here's how to encode the passwords if you choose not to use the defaults.

Assuming you want to set the server password to "server-pass":

``` shell
echo -n "server-pass" | openssl dgst -sha256 | sed 's/^.* //'
```

Which will respond with `7b1bf1e4f9535de960093f1c303fe35f49167bdc103ba99ad7dc9d62e2807a1d`

Likewise, if you want to set the client password to "client-pass":

``` shell
echo -n "client-pass" | openssl dgst -sha256 | sed 's/^.* //'
```

Which will respond with `fbfc2da74af1af1945ba7bf403cde789091e39b13c420170080872323dd2d148`

Now we're ready to edit the two configuration files for the `volantmq` service.

Using `vi` or your favorite editor, edit `./horizon/hzn.json`.

On the line `"HZN_ORG_ID": "testorg",`, replace "testorg" with your ORG ID. If you have it specified as an environment variable, you can skip this step.
On the line `"DOCKER_IMAGE_BASE": "openhorizon/volantmq",`, replace "openhorizon" with your dockerhub login.
On the line `"SERVICE_NAME": "ibm.volantmq",`, remove "ibm.".
On the line `"SERVICE_VERSION": "1.0.1"`, replace "1.0.1" with the current version you wish to set it to.

Next, edit `./horizon/service.definition.json`.

On the line `"url": "ibm.volantmq"`, remove "ibm.".
On the line that begins `"defaultValue": "fft-server: 7b1b`, replace the "fft-server" encoded password with your new encoded password ONLY IF you aren't using the default. 
Likewise for the "fft-client" encoded password on the same line.

Now we'll edit the two configuration files for the `fft_server` service.

Change to the `fft_server` directory:

``` shell
cd ../fft_server
```

Using `vi` or your favorite editor, edit `./horizon/hzn.json`.

On the line `"HZN_ORG_ID": "testorg",`, replace "testorg" with your ORG ID.
On the line `"DOCKER_IMAGE_BASE": "openhorizon/fft-server",`, replace "openhorizon" with your dockerhub login.
On the line `"SERVICE_NAME": "ibm.fft-server",`, remove "ibm.".
On the line `"SERVICE_VERSION": "1.0.7"`, replace "1.0.7" with the current version you wish to set it to.

Next, edit `./horizon/service.definition.json`.

On the line `"url": "ibm.fft-server"`, remove "ibm.".
On the line `"url": "ibm.volantmq",`, remove "ibm.".
On the line `{ "name": "MQTT_SERVER_PASS", "label": "", "type": "string", "defaultValue": "server-pass" },`, 
replace "server-pass" with your new clear text password ONLY IF you aren't using the default. 

Now we'll edit five configuration files for the `fft_client` service.

Change to the `fft_client` directory:

``` shell
cd ../fft_client
```

Using `vi` or your favorite editor, edit `./horizon/hzn.json`.

On the line `"HZN_ORG_ID": "testorg",`, replace "testorg" with your ORG ID.
On the line `"DOCKER_IMAGE_BASE": "openhorizon/fft-client",`, replace "openhorizon" with your dockerhub login.
On the line `"SERVICE_NAME": "ibm.fft-client",`, remove "ibm.".
On the line `"SERVICE_VERSION": "1.0.7"`, replace "1.0.7" with the current version you wish to set it to.

Next, edit `./horizon/service.definition.json`.

On the line `"url": "ibm.fft-client"`, remove "ibm.".
On the line `"url": "ibm.volantmq",`, remove "ibm.".
On the line `{ "name": "MQTT_CLIENT_PASS", "label": "", "type": "string", "defaultValue": "client-pass" },`, 
replace "client-pass" with your new clear text password ONLY IF you aren't using the default. 

Then edit `./horizon/pattern-all-arches.json`.

On the line `"serviceUrl": "ibm.fft-server"`, remove "ibm.".
On the line `"version": "1.0.7"`, replace "1.0.7" with the current version you wish to set it to.
On the next line `"serviceUrl": "ibm.fft-server"`, remove "ibm.".
On the next line `"version": "1.0.7"`, replace "1.0.7" with the current version you wish to set it to.

And then edit `./horizon/pattern.json`.

On the line `"serviceUrl": "ibm.fft-server"`, remove "ibm.".
On the line `"version": "1.0.7"`, replace "1.0.7" with the current version you wish to set it to.

Last, edit `./horizon/userinput.json`.

On the line `"MQTT_CLIENT_PASS": "client-pass",`, replace "client-pass" with your new clear text password ONLY IF you aren't using the default.
On the line `"DEVICE_ID": -1`, replace -1 with the actual device ID.  See the section on detecting device ID to ensure you're using the correct value.
On the line `"url": "ibm.fft-server",`, remove "ibm.".
On the line `"MQTT_SERVER_PASS": "server-pass",`, replace "server-pass" with your new clear text password ONLY IF you aren't using the default.
On the line `"url": "ibm.volantmq",`, remove "ibm.".
On the line that begins `"VOLANTMQ_USERS": "fft-server: 7b1b`, replace the "fft-server" encoded password with your new encoded password ONLY IF you aren't using the default. 
Likewise for the "fft-client" encoded password on the same line.

Now you're finally done editing!

### Step 3: Creating the Docker Images

This is an optional step that will need to be completed if you are not planning to use existing Docker images.  
Please feel free to skip this step otherwise.  

If you are planning to build the Docker images for both amd64 and armhf architectures, you'll want to build on an x86_64 machine using qemu. 
Assuming that you are running Ubuntu, the following will install the required packages:

``` shell
apt-get -y update
apt-get -y install gcc make perl jq curl git ssh qemu alsa-utils alsa-tools
```

Then run the following command to enable qemu:

``` shell
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
```

Additionally, the Makefile expects the following environment variables:

``` shell
export SERVICE_NAME=fft-client
```

Finally, you'll want to login to dockerhub so that you can publish the built Docker images to your dockerhup account.

``` shell
docker login
```

If that was successful, please continue, otherwise do not continue until you can login.

Now change your working directory to the `volantmq` service:

``` shell
cd ../volantmq
```

If building for all architectures:

``` shell
make build
```

Or if building for amd only:

``` shell
make build-amd
```

And push the built Docker image layers to dockerhub:

``` shell
docker push <your docker login>/volantmq
```

Now change your working directory to the `fft_server` service:

``` shell
cd ../fft_server
```

If building for all architectures:

``` shell
make build
```

Or if building for amd only:

``` shell
make build-amd
```

And push the built Docker image layers to dockerhub:

``` shell
docker push <your docker login>/fft-server
```

Last, change your working directory to the `fft_client` service:

``` shell
cd ../fft_client
```

If building for all architectures:

``` shell
make build
```

Or if building for amd only:

``` shell
make build-amd
```

And push the built Docker image layers to dockerhub:

``` shell
docker push <your docker login>/fft-client
```

### Step 4: Configuring and Publishing the Service and Pattern Definition Files

Ensure your `hzn` agent is connecting to your exchange and is authorized:

``` shell
hzn exchange user list
```

If not, you'll want to ensure that the following three environment variables are exported:

``` shell
export HZN_EXCHANGE_URL=<exchange URL, ex. http://192.168.1.138:3090/v1/>
export HZN_ORG_ID=<ex. myorg>
export HZN_EXCHANGE_USER_AUTH=<ex. admin:yl123AbCDefG>
```

Additionally, the Makefile expects the following environment variables:

``` shell
export SERVICE_NAME=fft-client
```

If you have never published a Service Definition file from this Agent to the Exchange, you'll also need to create signinng keys:

``` shell
hzn key create -l 4096 <HZN_ORG_ID ex. myorg> <HZN_EXCHANGE_USER ex. admin>
```

Now change your working directory to the `volantmq` service:

``` shell
cd ../volantmq
```

Then publish the Service Definition to the Exchange:

``` shell
make publish
```

Now change your working directory to the `fft_server` service:

``` shell
cd ../fft_server
```

Then publish the Service Definition to the Exchange:

``` shell
make publish
```

Last, change your working directory to the `fft_client` service:

``` shell
cd ../fft_client
```

Then publish the Service Definition to the Exchange:

``` shell
make publish-service
```

Then publish the Pattern Definition to the Exchange:

``` shell
make publish
```

You should now be able to see the services listed in the Exchange:

``` shell
hzn exchange service list
```

And likewise, the pattern in the Exchange:

``` shell
hzn exchange pattern list
```

### Step 5: Registering for the Pattern

You should already be in the `fft_client` folder and able to connect to the Exchange.

Run the following to register this edge compute node with the Exchange for the pattern, substituting your `HZN_ORG_ID` for "<ex. myorg>":

``` shell
hzn register -p "<ex. myorg>/pattern-fft-client" -f ./horizon/userinput.json
```

Once an agreement has been formed (see `hzn agreement list`) and the containers are running (see `docker ps`), you should be ready to use the FFT example service.  
Congratulations!

### Step 6: Testing the Deployed Example

Using Mosquitto or any other compatible message broker client on the same edge node, 
connect to the message queue "results" topic to watch for detected sound patterns.  
Replace "client-pass" below with your new clear text password ONLY IF you aren't using the default.

``` shell
mosquitto_sub -h localhost -p 1883 -t results -u "fft-client" -P "client-pass" -q 2
```

If the microphone does not detect the target sound, you should see it reporting "false".

### Step 7: Unregistering

To unregister your edge compute node from the example pattern, type:

``` shell
hzn unregister -f
```

## Ephemera

### Detecting Device ID

Valid devices to use for audio capture require one or two input channels. 
Here's an example list of connected audio devices: 

``` text
PortAudio version: 0x00130600
Version text: 'PortAudio V19.6.0-devel, revision unknown'
Found 5 devices:
   #   In Out  Sample Rate  Name
  --  --- ---  -----------  -----------------------------------
   0:   0  8         44100  bcm2835 ALSA: IEC958/HDMI (hw:0,1)
   1:   0  8         44100  bcm2835 ALSA: IEC958/HDMI1 (hw:0,2)
   2:   0  2         44100  Plantronics C320-M: USB Audio (hw:1,0)
   3:   1  0         44100  USB PnP Sound Device: Audio (hw:2,0)
   4:   0  2         48000  dmix
```

Notice that only device 3 has an input channel, so that device is the only one that can be used with this FFT example service.

To get a list of valid devices from the `portaudio` library, you can compile a small application. 

Install the ALSA SDK:

``` shell
apt-get -y install libasound-dev gcc make
```

Download `portaudio`:

``` shell
wget http://www.portaudio.com/archives/pa_stable_v190600_20161030.tgz
tar -zxvf pa_stable_v190600_20161030.tgz
cd portaudio
./configure && make
```

Create a file named `fad.c` and paste the following contents into it:

``` c
/*
  fad  Find Audio Devices
  gcc fad.c libportaudio.a -lrt -lm -lasound -pthread -o fad
  Written by glendarling@us.ibm.com
  Copyright (c) 2020, IBM; all rights reserved.
*/
#include <stdio.h>
#include <stdlib.h>
#include "portaudio.h"
void main() {
  PaError err;
  const PaDeviceInfo *deviceInfo;
  int numDevices;
  freopen("/dev/null", "w", stderr);
  err = Pa_Initialize();
  freopen("/dev/tty", "w", stderr);
  if( err != paNoError ) {
    fprintf(stderr, "ERROR: Pa_Initialize returned 0x%x\n", err);
    exit(1);
  }
  printf("PortAudio version: 0x%08X\n", Pa_GetVersion());
  printf("Version text: '%s'\n", Pa_GetVersionInfo()->versionText);
  numDevices = Pa_GetDeviceCount();
  if(numDevices < 0) {
    fprintf(stderr, "ERROR: Pa_CountDevices returned 0x%x\n", numDevices);
    exit(1);
  }
  printf("Found %d devices:\n", numDevices);
  printf("   #   In Out  Sample Rate  Name\n");
  printf("  --  --- ---  -----------  -----------------------------------\n");
  for(int i=0; i<numDevices; i++) {
    deviceInfo = Pa_GetDeviceInfo( i );
    /*
      Fields in deviceInfo:
      Type            Name
      --------------  ------------------------
      int             structVersion
      const char *    name
      PaHostApiIndex  hostApi
      int             maxInputChannels
      int             maxOutputChannels
      PaTime          defaultLowInputLatency
      PaTime          defaultLowOutputLatency
      PaTime          defaultHighInputLatency
      PaTime          defaultHighOutputLatency
      double          defaultSampleRate
    */
    printf("  %2d:  %2d %2d  %12.0f  %s\n",
      i,
      deviceInfo->maxInputChannels,
      deviceInfo->maxOutputChannels,
      deviceInfo->defaultSampleRate,
      deviceInfo->name);
  }
  exit(0);
}
```

Then compile the code with `gcc` like the following:

``` shell
cp lib/.libs/libportaudio.a .
cp include/portaudio.h .
gcc fad.c lib/.libs/libportaudio.a -lrt -lm -lasound -pthread -o fad
```

Now you can run `fad` like this:

``` shell
./fad
```

On my VM it found one microphone, device 1:t†

``` text
PortAudio version: 0x00130600
Version text: 'PortAudio V19.6.0-devel, revision 396fe4b6699ae929d3a685b3ef8a7e97396139a4'
Found 10 devices:
   #   In Out  Sample Rate  Name
  --  --- ---  -----------  -----------------------------------
   0:   2  6         48000  Intel 82801AA-ICH: - (hw:0,0)
   1:   2  0         44100  Intel 82801AA-ICH: MIC ADC (hw:0,1)
   2:  128 128         48000  sysdefault
   3:   0  6         48000  front
   4:   0  4         48000  surround40
   5:   0 128         48000  surround41
   6:   0 128         48000  surround50
   7:   0 128         48000  surround51
   8:  128 128         48000  default
   9:   0  2         48000  dmix
```

Alternatively you can run `fft-server` docker manually and use `--list_devices` flag:

```shell
docker run -it --rm --entrypoint /bin/sh --device /dev/snd DOCKER_ID/fft-client:ARCH-VERSION


./client -b ${MQTT_BROKER} -u ${MQTT_CLIENT_USER} -p ${MQTT_CLIENT_PASS} -c ${MQTT_CLIENT_CLIENT} -r ${SAMPLE_RATE} -f ${RECORD_FRA
ME} -q ${MQTT_QOS} -l ${LOG_LEVEL} --device_id ${DEVICE_ID} --list_devices
```

Where DOCKER_ID and VERSION are variables from fft_server/hzn.json and ARCH is currene architecture you're running on: amd64 or arm.

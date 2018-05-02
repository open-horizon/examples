# Edge Developer Quickstart Guide

This Developer Quickstart Guide provides a simplified description of the process for developing, testing and deploying user-developed code in the Edge environment.
The [Edge Service Development Guidelines](Edge-Service-Development-Guidelines.md) provides guidance on how to best structure your service so it runs well in the WIoTP/Horizon Edge fabric.
The [Edge Developer Guide](https://github.com/open-horizon/examples/wiki/Edge-Developer-Guide) is a more detailed description of the Edge environment and the concerns that an Edge developer has to be aware of.

Note there is a concise [Quick Start Guide](Edge-Quick-Start-Guide.md) available, that shows how to get an existing workload up and running on your edge nodes very quickly without having to develop any code. **That [Quick Start Guide](Edge-Quick-Start-Guide.md) is also a prerequisite for this guide.**

Additional information is available, and questions may be asked, in our forum, at:
* [https://discourse.bluehorizon.network/](https://discourse.bluehorizon.network/)

The Edge is based upon the open source Horizon project [https://github.com/open-horizon](https://github.com/open-horizon). There are therefore several references to Horizon, and the `hzn` Linux command in this document.

Edge simplifies and secures the global deployment and maintenance of software on IoT edge nodes.
This document will guide you through the process of building, testing, and deploying IoT edge software, using the Watson IoT Platform to securely deploy and then fully manage the software on your IoT edge nodes all over the world.
IoT software maintenance in Edge with Watson IoT becomes fully automatic (zero-touch for your edge nodes), highly secure and easy to centrally control.

## Overview

This guide is intended for developers who want to experience the Edge software development process with Watson IoT Platform, using a simple example.
In this guide you will learn how to create Horizon microservices and workloads, how to test them and ultimately how to integrate them with the Watson IoT Platform.

As you progress through this guide, you will first build a simple microservice that extracts CPU usage information from the underlying edge node.
Then you will build a workload that samples CPU usage information from the microservice, computes an average and then publishes the average to Waton IoT Platform.
Along the way, you will be exposed to many concepts and capabilities that are documented in complete detail in the [Edge Developer Guide](https://github.com/open-horizon/examples/wiki/Edge-Developer-Guide).

## Before you begin

Currently this guide is intended to be used on an x86_64 machine. (This will be expanded in the future.)

To familiarize yourself with WIoTP Edge, we suggest you go through the entire [Quick Start Guide](Edge-Quick-Start-Guide.md). But even if you do not go through that entire guide, **you must at least do the first sections of it, up to and including [Verify Your Gateway Credentials and Access](Edge-Quick-Start-Guide.md#verify-your-gateway-credentials-and-access)**, on the same edge node that you are using for this guide. (**For now, use the commented out line `aptrepo=testing` in the apt repo section, so you get the latest Horizon debian packages. They are currently required for this guide. You should have at least version 2.17.2**) That guide will have you accomplish the following necessary steps:

- Create your WIoTP organization, gateway type and id, and API key.
- Install docker, horizon, and some utilities.
- Set environment variables needed in the rest of this guide.
- Verify your edge node's access to the WIoTP cloud services.

After completing those steps in the [Quick Start Guide](Edge-Quick-Start-Guide.md) , continue here. Install some basic development tools and clone the examples repo:
```bash
apt install -y git make
cd ~
git clone https://github.com/open-horizon/examples.git
```

## Create a microservice project

A typical edge service has 2 parts to it: a service that accesses data that is available on this edge node (a microservice), and logic that does analysis/processing of the data and optionally sends consolidated data to the cloud (workload). 

We will start the microservice project by creating a docker container that will respond to an HTTP request with CPU usage information.
The microservice implementation is a very simple bourne shell script.
On your development machine, create a project directory:
```bash
mkdir -p ~/hzn/ms/cpu
cd ~/hzn/ms/cpu
```

Next, create a docker container that exposes an HTTP API for obtaining CPU usage information.
Normally you would:
* Author a Dockerfile to hold the container definition, for example `~/hzn/ms/cpu/Dockerfile`.
* Author a shell script that will run when the container starts, for example `~/hzn/ms/cpu/start.sh`.
* Author code that runs when the TCP listener gets a message, for example `~/hzn/ms/cpu/service.sh`.
* Author a Makefile that will build and test this container on your development machine, for example `~/hzn/ms/cpu/Makefile`.
* Author Horizon metadata that enables Horizon to manage your microservice.

But the easiest method is to copy an existing microservice that you want to start from:
```bash
cp -a ~/examples/edge/services/cpu_percent/* .
```

The Makefile and several of the horizon files contain environment variables that will be replaced by their values when the files are used. Copy `horizon/envvars.sh.sample` to `horizon/envvars.sh`, edit it to set the environment variables to your own values (including getting a docker hub id, if you don't already have one), then source the file so its values will be available to the rest of the commands:
```bash
cp horizon/envvars.sh.sample horizon/envvars.sh
# put your values in horizon/envvars.sh (there is a .gitignore file for that)
source horizon/envvars.sh
```

The Makefile is setup to build the container, start it and run a simple test to ensure that it works.
1. Make the microservice container and run it locally:
```
make
```
1. Verify that the container was built and started successfully by looking for the output of the curl test at the end of the output:
```
    curl -s http://localhost:8347/v1/cpu | jq .
    {
      "cpu": 1.34
    }
```
1. Stop the running container:
```
make stop
```

### Microservice project metadata

Now that the microservice container implementation is working correctly, you can use the Horizon project metadata to enable Horizon to run it now, and ultimately deploy it to Edge nodes. The examples project you copied is "Horizon ready", so it already contains the Horizon metadata files. Note the files in the `horizon` sub-directory:
* `microservice.definition.json` - the Horizon metadata of this microservice. Note a few of the significant json fields:
    * `specRef`: along with the `version` and `arch` this is the unique identifier for this microservice, and ideally a URL to a web site that documents the microservice for potential users of it.
    * `deployment`: contains the docker image(s) that make up this microservice, and how the Horizon agent should run them on each edge node:
        * `image` - the full docker image name (including the registry, if not in docker hub)
        * `cpu` - this field name is also used as the docker defined DNS name that workloads can use to contact it.
* `userinput.json` - the runtime input values specified by the edge node owner for the microservice or the Horizon agent. Microservices should require as little input from the edge node owners as possible, ideally none, which is the case here. Note within the file:
    * `url`: this must match the `specRef` in `microservice.definition.json`

The Horizon CLI contains a set of sub-commands of `hzn dev` that are useful for running your microservices and workloads in a development environment. Verify that the project metadata has no errors in it.
```
cd ~/hzn/ms/cpu
hzn dev microservice verify
```

If the verify sub-command finds any inconsistencies or errors it will report the file and the location within the file where the error was detected.

### Test the microservice project
Test your microservice in the Horizon test environment.
The test environment closely simulates the Edge node environment in which your microservice will run when deployed to Edge nodes.
```
hzn dev microservice start
```
After a few seconds, the microservice will be started. The `docker ps` command will display the running container.
Notice that the microservice docker container name was formed from values in the microservice definition: `specRef`, `version`, and the service name (and a unique UUID).
The microservice is attached to a docker network. See it via `docker network ls`. The network has nearly the same name as the microservice container name (minus the service name).

Also note that the running container exposes no host network ports.
Microservices on a Horizon Edge node, run in a sandboxed environment where network access is restricted only to workloads that require the microservice.
This means that when running in a Horizon test environment, by default, the host cannot access the microservice's network ports.
This can be overridden, but in general is the desired and expected behavior.

Finally, look at the environment variables that have been passed into the container:
```
docker inspect `docker ps -q` | jq '.[0].Config.Env'
```
The variables prefixed with "HZN" are provided by the Horizon Edge node platform.
A microservice container can exploit these environment variables.
See the [Horizon Environment Variable Documentation](https://github.com/open-horizon/anax/blob/master/doc/managed_workloads.md) for a complete description of Horizon Edge node environment variables.

The microservice container test environment can be stopped (it will be started again later when the workload requires it):
```
hzn dev microservice stop
```

## Create a workload project

Workloads use microservices to gain access to data on the Edge node, and then processes that data in some way.
Let's develop a minimal workload and explore more complex usages of the `hzn dev` sub-commands.

On your development machine, create a project directory:
```bash
mkdir -p ~/hzn/workload/cpu2wiotp
cd ~/hzn/workload/cpu2wiotp
```

Create a docker container that computes averages for a set of CPU usage samples. As with a microservice, normally you would author code, a Dockerfile, a Makefile, and Horizon metadata. But the easiest method is to copy an existing workload that you want to start from:
```bash
cp -a ~/examples/edge/wiotp/cpu2wiotp/* .
```

The Makefile and several of the horizon files contain environment variables that will be replaced by their values when the files are used. Edit `horizon/envvars.sh` to set the environment variables to your own values and source it:
```bash
cp horizon/envvars.sh.sample horizon/envvars.sh
# put your values in horizon/envvars.sh (there is a .gitignore file for that)
source horizon/envvars.sh
```

The Makefile is setup to build the container, start it, and run a simple test to ensure that it works.
1. Make the workload container and run it locally:
```
make
```
1. Verify that the container was built and started successfully, and look for output like this at the bottom:
```
Starting infinite loop to read from microservice then publish...
 Interval 1 cpu: 53
 Interval 2 cpu: 54
```
1. Control-C to stop `docker logs` and then stop the running container:
```
make stop
```

### Workload project metadata

Now that the workload container implementation is working correctly, let's take a look at the Horizon project metadata in order to make this project testable in the Horizon test environment and to ultimately make it deployable to Edge nodes. The examples project you copied is "Horizon ready", so it already contains the Horizon metadata files. Note the files in the `horizon` sub-directory:

* `workload.definition.json` - the Horizon metadata of this workload. Note a few of the significant json fields:
    * `workloadUrl`: along with the `version` and `arch` this is the unique identifier for this workload, and ideally a URL to a web site that documents the workload for potential users of it.
    * `specRef`: a microservice that this workload uses/depends on. This must match the `specRef` value in the microservice.definition.json file that the workload depends on.
    * `userInput`: this section defines the input values that can be specified to this workload by the edge node owner. The `userinput.json` file must include a value for each variable that doesn't have a default value.
    * `deployment`: contains the docker image(s) that make up this workload, and how the Horizon agent should run them on each edge node:
        * `image` - the full docker image name (including the registry, if not in docker hub)
        * `cpu2wiotp` - this field name is also used as the docker defined DNS name that can be used to contact it.
        * `environment` - environment variables that should be passed to the containers (in addition to the variables Horizon automatically passes)
* `userinput.json` - the runtime input values specified by the edge node owner for the workload or the Horizon agent. Microservices should require as little input from the edge node owners as possible, ideally none, which is the case here. Note within the file:
    * `url`: this must match the `workloadUrl` in `workload.definition.json`
* `horizon/dependencies` - This directory will be populated later by `hzn dev dependency fetch` to contain metadata describing other projects (microservices) that this project depends on.

The default `userinput.json` configures cpu2wiotp to get data from the cpu microservice and publish data to WIoTP. We will get there eventually, but first we want to run cpu2wiotp in "stand-alone" mode in which it makes up its own cpu data and just prints it instead of sending it to WIoTP. To accomplish that, update the `workloads` section of `userinput.json` to look like this:
```
    "workloads": [
        {
            "org": "$HZN_ORG_ID",
            "url": "https://$MYDOMAIN/workloads/$CPU2WIOTP_NAME",
            "versionRange": "[0.0.0,INFINITY)",
            "variables": {
                "SAMPLE_SIZE": 5,
                "SAMPLE_INTERVAL": 2,
                "MOCK": true,
                "PUBLISH": false,
                "VERBOSE": "1"
            }
        }
    ]
```

Verify that the project has no errors in it.
```
cd ~/hzn/workload/cpu2wiotp
hzn dev workload verify
```

If the verify sub-command finds any inconsistencies or errors it will report the file and the location within the file were the error was detected.

### Test the workload project

Test your workload in the Horizon test environment.
The test environment closely simulates the Edge node environment in which your workload will run when deployed to Edge nodes.
```
hzn dev workload start
```
After a few seconds, the workload will be started. The `docker ps` command will display the running container.
Notice that the workload docker container name is derived from the hexadecimal agreement id and the service name.
An agreement id is a hexadecimal string that uniquely identifies the agreement.
When your workload is eventually deployed by a Horizon agbot, an agreement id will be assigned and it will appear similarly to how you see it now.

The workload is attached to a docker network. See it via `docker network ls`. The network has nearly the same name as the workload container name (minus the service name).

Finally, look at the environment variables that have been passed into the container:
```
docker inspect `docker ps -q` | jq '.[0].Config.Env'
```
Variables prefixed with "HZN" are provided by the Horizon Edge node platform.
See the [Environment Variable Documentation](https://github.com/open-horizon/anax/blob/master/doc/managed_workloads.md) for a complete description of Horizon Edge node variables.
The rest of the variables came either from `userinput.json` or from the `environment` section of `workload.definition.json`.

The log output from the workload container can be found in syslog:
```
tail -f /var/log/syslog | grep cpu2wiotp
```

Control-C to stop `tail`. Then the workload container test environment can be stopped:
```
hzn dev workload stop
```

## Add a workload dependency

So far, the workload is executing with mocked CPU usage data, not very interesting.
Update the workload project to use the CPU usage microservice project we created earlier:
```bash
hzn dev dependency fetch -p ~/hzn/ms/cpu/horizon
```

The addition of the new dependency has updated several pieces of metadata in your workload project. Notice that:
* `horizon/dependencies` now contains a copy of the dependency's microservice definition.
* The dependency was added to the `apiSpec` array in the `workload.definition.json` file. (In our case it was already there, but would have been added if it wasn't.)
* The `userinput.json` file was updated to include variable values for the new dependency. (In our case it was already there, but would have been added if it wasn't.)

If the Horizon metadata of a dependent project changes, the dependency must be recreated with the `hzn dev dependency fetch` command used above.

Remember that the workload is currently configured to mock the CPU usage microservice. Now we want to use the real data from the cpu microservice, so
remove the "MOCK" variable from the `userinput.json` file so that the default value of "false" will be used:
```
    "variables": {
        "PUBLISH": false,
        "SAMPLE_INTERVAL": 2,
        "SAMPLE_SIZE": 5,
        "VERBOSE": "1"
    }
```

Now start the workload again, and this time both the microservice and the workload will be started.
```
hzn dev workload start
```

After a few seconds, the log output from the workload container can be found in syslog:
```
tail -f /var/log/syslog | grep cpu2wiotp
```

Notice that the workload is now picking up real CPU usage samples from the microservice and computing an average.
This is accomplished without any open ports on the host because the workload container has joined the docker network of the microservice. (This can be seen using the commands `docker ps`, `docker network ls`, and `docker inspect <id>`.)
This is how the containers will be configured when running on Edge nodes.

Control-C to stop `tail`. Then the workload container and the microservice container test environment can be stopped:
```
hzn dev workload stop
```

Let's review what has been accomplished so far:
* A microservice project has been authored which produces a standalone container that can be independently tested.
* A workload project has been authored which produces a standalone container that can be independently tested.
* The workload project has also been extended to depend on the microservice project such that both projects can be run in the Horizon test environment, and tested together.

## Publish Data to Watson IoT Platform

This section shows you how to publish the CPU usage average to your Watson IoT Platform, using WIoTP's core-iot microservice.

Modify your project as follows:
* Update `userinput.json` under the `workload` section to enable the workload to publish the CPU average by setting PUBLISH to `true`:
```
    "variables": {
        "PUBLISH": true,
        "SAMPLE_INTERVAL": 2,
        "SAMPLE_SIZE": 5,
        "VERBOSE": "1"
    }
```
* Add the Wation IoT Platform core IoT microservice as a dependency of your workload project.
```
hzn dev dependency fetch -s https://internetofthings.ibmcloud.com/wiotp-edge/microservices/edge-core-iot-microservice --ver 2.4.0 -o IBM -a $ARCH -k /etc/horizon/trust/publicWIoTPEdgeComponentsKey.pem
```
* The command above would normally have added a section to the `userinput.json` file under `microservices` for the edge-core-iot-microservice, but it recognized that we already had that section.
* Note that the `deployment` section of `workload.definition.json` already contains `environment` and `binds` settings appropriate for using the core-iot microservice, because we inherited that when we copied from `~/examples/edge/wiotp/cpu2wiotp`.
* Set up some things the core-iot microservice needs:
```
cp /etc/wiotp-edge/edge.conf.template /etc/wiotp-edge/edge.conf
mkdir -p /var/wiotp-edge/persist
wiotp_create_certificate -p $WIOTP_GW_TOKEN
```
* Start the workload and look for the CPU messages in your Watson IoT Platform instance:
```bash
hzn dev workload start
```

* Run `docker ps` and notice the additional core-iot containers that are running. (The core-iot service provides many other functions in addition to the one we are using: sending MQTT messages to WIoTP cloud.)

* Subscribe to the WIoTP cloud MQTT topic that the workload is publishing to to see the cpu values:
```
mosquitto_sub -v -h $HZN_ORG_ID.messaging.$WIOTP_DOMAIN -p 8883 -i "$WIOTP_CLIENT_ID_APP" -u "$WIOTP_API_KEY" -P "$WIOTP_API_TOKEN" --capath /etc/ssl/certs -t iot-2/type/$WIOTP_GW_TYPE/id/$WIOTP_GW_ID/evt/status/fmt/json
```

* If you don't see messages coming to that, look again at the workload log for errors. (Note: it is normal to sometimes get 1 error at the beginning of the cpu2wiotp log, depending on how long it takes each container to initialize.)
```
tail -f /var/log/syslog | grep cpu2wiotp
```

* Control-C to stop `tail`. Then stop the workload:
```bash
hzn dev workload stop
```

At this point you have successfully completed a quick pass through the development process.
You have developed and tested a microservice and a workload as standalone containers and running within the Horizon Edge node test environment.
You have also integrated your projects with your Watson IoT Platform instance and are able to publish data to it.
In order to make these projects available for other Edge nodes to run them, you need to publish your projects.
The next section describes how to use `hzn dev` to do that.

## Deploying the projects to WIoTP/Horizon

When you are satisfied that the microservice and workload are working correctly, you can deploy them to the WIoTP so that any Edge node in your organization can run them.
The first step is to create a key pair that will be used to sign your docker images and the deployment configuration of your microservice and workload.
```bash
cd ~/hzn
hzn key create <x509-org> <x509-cn>
export PRIVATE_KEY_FILE=~/hzn/*-private.key
export PUBLIC_KEY_FILE=~/hzn/*-public.pem
```
where `x509-org` is a company name or organization name that is suitable to be used as an x509 certificate organization name, and `x509-cn` is an x509 certificate common name (preferably an email address issued by the `x509-org` organization).

The private key will be used to sign the microservice and workload, the public key is needed by any Edge node that wants to run the microservice and workload, so that it can verify their signatures. (Horizon can take care of distributing the public key to the edge nodes that need it.)

You will be storing your docker images in Docker Hub, so login now:
```
docker login -u $DOCKER_HUB_ID
```

Now publish the microservice:
```bash
cd ~/hzn/ms/cpu
hzn dev microservice publish -k $PRIVATE_KEY_FILE  # soon you can use -K $PUBLIC_KEY_FILE and then will not have to import it
```

You can verify that the microservice was published:
```bash
hzn exchange microservice list
```

Publish the workload:
```bash
cd ~/hzn/workload/cpu2wiotp
hzn dev workload publish -k $PRIVATE_KEY_FILE  # soon you can use -K $PUBLIC_KEY_FILE and then will not have to import it
```

You can verify that the workload was published:
```bash
hzn exchange workload list
```

**If you previously went through the entrie [Quick Start Guide](Edge-Quick-Start-Guide.md) using this same gateway type, remove the workload you added to that pattern:**
```
hzn exchange pattern removeworkload $WIOTP_GW_TYPE $HZN_ORG_ID https://internetofthings.ibmcloud.com/workloads/cpu2wiotp $ARCH
```

Add this workload to your gateway type deployment pattern and verify it is there:
```
hzn exchange pattern insertworkload -k $PRIVATE_KEY_FILE -f pattern/insert-cpu2wiotp.json $WIOTP_GW_TYPE  # soon you can use -K $PUBLIC_KEY_FILE and then will not have to import it
hzn exchange pattern list $WIOTP_GW_TYPE | jq .
```

## Using your workload on edge nodes of this type

You, or others in your organization, can now use this workload on many edge nodes. On each of those nodes:

* **If you have run thru this document before** on this edge node, do this to clean up:
```
hzn unregister -f
```

* Import your signing public key (soon you won't need to do this):
```
hzn key import -k $PUBLIC_KEY_FILE
```

* Register the node and start the Watson IoT Platform core-IoT service and the CPU workload:
```
wiotp_agent_setup --org $HZN_ORG_ID --deviceType $WIOTP_GW_TYPE --deviceId $WIOTP_GW_ID --deviceToken "$WIOTP_GW_TOKEN"
```

Note: in the command above, we did not specify an input file with user input for your workload, so all of the workload's `userInput` variables are set to their default values. This is the preferred way to run your workload, so it can easily be run on many edge nodes.

After a short while, usually within just a minute or two, agreements will be made to run the WIoTP core-iot service and your workload. See [Register Your Edge Node](Edge-Quick-Start-Guide.md#register-your-edge-node) as a reminder for how to check for agreements and your workload, and how to verify that CPU values are being sent to the cloud.

## Advanced topic - Expanding your project

Since your project Makefile and metadata files have environment variables in them for architecture, version, etc., you can easily test your services for different architectures, publish them for other gateway types, etc.

Here are a few things to consider when expanding your project:
* New versions of your microservice and workload:
    * When you want to roll out a new version, build the new docker images (with a new tag for the new version) and create new microservice/workload definitions with the new version number. Then replace the workload reference in your gateway type pattern.
    * The edge nodes that are already using that pattern will automatically see the new workload within a few minutes and re-negotiate the agreement to run the new version. You can manually force this by canceling the current agreement using `hzn agreement cancel <agreement-id>`
* Additional gateway types and architectures:
    * If you want multiple gateway types to run your workload, you do not need to create the microservice/workload definitions multiple times. You only need to add your workload reference to each gateway type pattern.
    * With WIoTP gateway types with Edge Services enabled, each type can only be used for a single architecture. So create a different gateway type for each architecture you want to support.

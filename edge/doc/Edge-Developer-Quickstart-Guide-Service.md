# Edge Developer Quickstart Guide

This Developer Quickstart Guide provides a simplified description of the process for developing, testing and deploying user-developed code in the Edge environment.
The [Edge Service Development Guidelines](Edge-Service-Development-Guidelines.md) provides guidance on how to best structure your service so it runs well in the WIoTP/Horizon Edge fabric.
The [Edge Developer Guide](https://github.com/open-horizon/examples/wiki/Edge-Developer-Guide) is a more detailed description of the Edge environment and the concerns that an Edge developer has to be aware of.

Note there is a concise [Quick Start Guide](Edge-Quick-Start-Guide.md) available, that shows how to get an existing service up and running on your edge nodes very quickly without having to develop any code. **That [Quick Start Guide](Edge-Quick-Start-Guide.md) is also a prerequisite for this guide.**

Additional information is available, and questions may be asked, in our forum, at:
* [https://discourse.bluehorizon.network/](https://discourse.bluehorizon.network/)

The Edge is based upon the open source Horizon project [https://github.com/open-horizon](https://github.com/open-horizon). There are therefore several references to Horizon, and the `hzn` Linux command in this document.

Edge simplifies and secures the global deployment and maintenance of software on IoT edge nodes.
This document will guide you through the process of building, testing, and deploying IoT edge software, using the Watson IoT Platform to securely deploy and then fully manage the software on your IoT edge nodes all over the world.
IoT software maintenance in Edge with Watson IoT becomes fully automatic (zero-touch for your edge nodes), highly secure and easy to centrally control.

## Overview

This guide is intended for developers who want to experience the Edge software development process with Watson IoT Platform, using a simple example.
In this guide you will learn how to create Horizon services, how to test them and ultimately how to integrate them with the Watson IoT Platform.

As you progress through this guide, you will first build a simple service that extracts CPU usage information from the underlying edge node.
Then you will build a another service that samples CPU usage information from the service, computes an average and then publishes the average to Waton IoT Platform.
Along the way, you will be exposed to many concepts and capabilities that are documented in complete detail in the [Edge Developer Guide](https://github.com/open-horizon/examples/wiki/Edge-Developer-Guide).

## Before you begin

Currently this guide is intended to be used on an x86_64 machine. (This will be expanded in the future.)

To familiarize yourself with WIoTP Edge, we suggest you go through the entire [Quick Start Guide](Edge-Quick-Start-Guide.md). But even if you do not go through that entire guide, **you must at least do the first sections of it, up to and including [Verify Your Gateway Credentials and Access](Edge-Quick-Start-Guide.md#verify-your-gateway-credentials-and-access)**, on the same edge node that you are using for this guide. (**For now, use the commented out line `aptrepo=testing` in the apt repo section, so you get the latest Horizon debian packages. They are currently required for this guide. You should have at least version 2.17.6. For now you must also use the bluehoriozn hybrid environment for this guide, because the agbot and exchange versions in the WIoTP production environment do not yet support services.**) That guide will have you accomplish the following necessary steps:

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

## Create a service project

A typical edge application has 2 parts to it: a service that accesses data that is available on this edge node, and another service that contains logic that does analysis/processing of the data and optionally sends consolidated data to the cloud. 

We will start the service project by creating a docker container that will respond to an HTTP request with CPU usage information.
The service implementation is a very simple bourne shell script.
On your development machine, create a project directory:
```bash
mkdir -p ~/hzn/service/cpu
cd ~/hzn/service/cpu
```

Next, create a docker container that exposes an HTTP API for obtaining CPU usage information.
Normally you would:
* Author a Dockerfile to hold the container definition, for example `~/hzn/service/cpu/Dockerfile`.
* Author a shell script that will run when the container starts, for example `~/hzn/service/cpu/start.sh`.
* Author code that runs when the TCP listener gets a message, for example `~/hzn/service/cpu/service.sh`.
* Author a Makefile that will build and test this container on your development machine, for example `~/hzn/service/cpu/Makefile`.
* Author Horizon metadata that enables Horizon to manage your service.

But the easiest method is to copy an existing service that you want to start from:
```bash
cp -a ~/examples/edge/services/cpu_percent/* .
mv horizon horizon.microservice; mv horizon.service horizon

```

The Makefile and several of the horizon files contain environment variables that will be replaced by their values when the files are used. Copy `horizon/envvars.sh.sample` to `horizon/envvars.sh`, edit it to set the environment variables to your own values (including getting a docker hub id, if you don't already have one), then source the file so its values will be available to the rest of the commands:
```bash
cp horizon/envvars.sh.sample horizon/envvars.sh
# put your values in horizon/envvars.sh (there is a .gitignore file for that)
source horizon/envvars.sh
```

The Makefile is setup to build the container, start it and run a simple test to ensure that it works.
1. Make the service container and run it locally:
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

### The CPU Service project metadata

Now that the service container implementation is working correctly, you can use the Horizon project metadata to enable Horizon to run it now, and ultimately deploy it to Edge nodes. The examples project you copied is "Horizon ready", so it already contains the Horizon metadata files. Note the files in the `horizon` sub-directory:
* `service.definition.json` - the Horizon metadata of this service. Note a few of the significant json fields:
    * `url`: along with the `version` and `arch` this is the unique identifier for this service, and ideally a URL to a web site that documents the service for potential users of it.
    * `deployment`: contains the docker image(s) that make up this service, and how the Horizon agent should run them on each edge node:
        * `image` - the full docker image name (including the registry, if not in docker hub)
        * `cpu` - this field name is also used as the docker defined DNS name that other services can use to contact it.
* `userinput.json` - the runtime input values specified by the edge node owner for the service or the Horizon agent. Services should require as little input from the edge node owners as possible, ideally none, which is the case here. Note within the file:
    * `url`: this must match the `url` in `service.definition.json`

The Horizon CLI contains a set of sub-commands of `hzn dev` that are useful for running your services in a development environment. Verify that the project metadata has no errors in it.
```
cd ~/hzn/service/cpu
hzn dev service verify
```

If the verify sub-command finds any inconsistencies or errors it will report the file and the location within the file where the error was detected.

### Test the CPU service project
Test your service in the Horizon test environment.
The test environment closely simulates the Edge node environment in which your service will run when deployed to Edge nodes.
```
hzn dev service start
```
After a few seconds, the service will be started. The `docker ps` command will display the running container.
Notice that the service docker container name was formed from values in the service definition: `url`, `version`, and the service name (and a unique UUID).
The service is attached to a docker network. See it via `docker network ls`. The network has nearly the same name as the service container name (minus the service name).

Also note that the running container exposes no host network ports.
Services on a Horizon Edge node, run in a sandboxed environment where network access is restricted only to services that require the service.
This means that when running in a Horizon test environment, by default, the host cannot access the service's network ports.
This can be overridden, but in general is the desired and expected behavior.

Finally, look at the environment variables that have been passed into the container:
```
docker inspect `docker ps -q` | jq '.[0].Config.Env'
```
The variables prefixed with "HZN" are provided by the Horizon Edge node platform.
A service container can exploit these environment variables.
See the [Horizon Environment Variable Documentation](https://github.com/open-horizon/anax/blob/master/doc/managed_workloads.md) for a complete description of Horizon Edge node environment variables.

The service container test environment can be stopped (it will be started again later when the cpu2wiotp service requires it):
```
hzn dev service stop
```

## Create a cpu2wiotp service project

Services use services to gain access to data on the Edge node, and then processes that data in some way.
Let's develop a minimal cpu2wiotp service and explore more complex usages of the `hzn dev` sub-commands.

On your development machine, create a project directory:
```bash
mkdir -p ~/hzn/service/cpu2wiotp
cd ~/hzn/service/cpu2wiotp
```

Create a docker container that computes averages for a set of CPU usage samples. As with the cpu service, normally you would author code, a Dockerfile, a Makefile, and Horizon metadata. But the easiest method is to copy an existing service that you want to start from:
```bash
cp -a ~/examples/edge/wiotp/cpu2wiotp/* .
mv horizon horizon.workload; mv horizon.service horizon
```

The Makefile and several of the horizon files contain environment variables that will be replaced by their values when the files are used. Edit `horizon/envvars.sh` to set the environment variables to your own values and source it:
```bash
cp horizon/envvars.sh.sample horizon/envvars.sh
# put your values in horizon/envvars.sh (there is a .gitignore file for that)
source horizon/envvars.sh
```

The Makefile is setup to build the container, start it, and run a simple test to ensure that it works.
1. Make the service container and run it locally:
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

### cpu2wiotp service project metadata

Now that the cpu2wiotp service container implementation is working correctly, let's take a look at the Horizon project metadata in order to make this project testable in the Horizon test environment and to ultimately make it deployable to Edge nodes. The examples project you copied is "Horizon ready", so it already contains the Horizon metadata files. Note the files in the `horizon` sub-directory:

* `service.definition.json` - the Horizon metadata of this service. Note a few of the significant json fields:
    * `url`: along with the `version` and `arch` this is the unique identifier for this service, and ideally a URL to a web site that documents the service for potential users of it.
    * `requiredServices.url`: a service that this service uses/depends on. This must match the `url` value in the cpu service.definition.json file that the cpu2wiotp service depends on.
    * `userInput`: this section defines the input values that can be specified to this service by the edge node owner. The `userinput.json` file must include a value for each variable that doesn't have a default value.
    * `deployment`: contains the docker image(s) that make up this service, and how the Horizon agent should run them on each edge node:
        * `image` - the full docker image name (including the registry, if not in docker hub)
        * `cpu2wiotp` - this field name is also used as the docker defined DNS name that can be used to contact it.
        * `environment` - environment variables that should be passed to the containers (in addition to the variables Horizon automatically passes)
* `userinput.json` - the runtime input values specified by the edge node owner for the service or the Horizon agent. Services should require as little input from the edge node owners as possible, ideally none, which is the case here. Note within the file:
    * `url`: this must match the `url` in `service.definition.json`
* `horizon/dependencies` - This directory will be populated later by `hzn dev dependency fetch` to contain metadata describing other projects (services) that this project depends on.

The default `userinput.json` configures cpu2wiotp to get data from the cpu service and publish data to WIoTP. We will get there eventually, but first we want to run cpu2wiotp in "stand-alone" mode in which it makes up its own cpu data and just prints it instead of sending it to WIoTP. To accomplish that, update the entry in the `services` section of `userinput.json` for the cpu2wiotp service:
```
        {
            "org": "$HZN_ORG_ID",
            "url": "https://$MYDOMAIN/services/$CPU2WIOTP_NAME",
            "versionRange": "[0.0.0,INFINITY)",
            "variables": {
                "SAMPLE_SIZE": 5,
                "SAMPLE_INTERVAL": 2,
                "MOCK": true,
                "PUBLISH": false,
                "VERBOSE": "1"
            }
        }
```

Likewise, `service.definition.json` has 2 dependencies listed in the `requiredServices` array. Remove both of them:
```
"requiredServices": [],
```

They will be added back automatically later on in this guide.

Verify that the project has no errors in it.
```
cd ~/hzn/service/cpu2wiotp
hzn dev service verify
```

If the verify sub-command finds any inconsistencies or errors it will report the file and the location within the file were the error was detected.

### Test the cpu2wiotp service project

Test your cpu2wiotp service in the Horizon test environment.
The test environment closely simulates the Edge node environment in which your service will run when deployed to Edge nodes.
```
hzn dev service start
```
After a few seconds, the cpu2wiotp service will be started. The `docker ps` command will display the running container.
Notice that the service docker container name is derived from the hexadecimal agreement id and the service name.
An agreement id is a hexadecimal string that uniquely identifies the agreement.
When your service is eventually deployed by a Horizon agbot, an agreement id will be assigned and it will appear similarly to how you see it now.

The service is attached to a docker network. See it via `docker network ls`. The network has nearly the same name as the service container name (minus the service name).

Finally, look at the environment variables that have been passed into the container:
```
docker inspect `docker ps -q` | jq '.[0].Config.Env'
```
Variables prefixed with "HZN" are provided by the Horizon Edge node platform.
See the [Environment Variable Documentation](https://github.com/open-horizon/anax/blob/master/doc/managed_workloads.md) for a complete description of Horizon Edge node variables.
The rest of the variables came either from `userinput.json` or from the `environment` section of `service.definition.json`.

The log output from the service container can be found in syslog:
```
tail -f /var/log/syslog | grep cpu2wiotp
```

Control-C to stop `tail`. Then the service container test environment can be stopped:
```
hzn dev service stop
```

## Add the CPU service dependency

So far, the cpu2wiotp service is executing with mocked CPU usage data, not very interesting.
Update the cpu2wiotp service project to use the CPU usage service project we created earlier:
```bash
hzn dev dependency fetch -p ~/hzn/service/cpu/horizon
```

The addition of the new dependency has updated several pieces of metadata in your cpu2wiotp service project. Notice that:
* `horizon/dependencies` now contains a copy of the dependency's service definition.
* The dependency was added to the `requiredServices` array in the `service.definition.json` file.
* The `userinput.json` file was updated to include variable values for the new dependency. (In our case it was already there, but would have been added if it wasn't.)

If the Horizon metadata of a dependent project changes, the dependency must be recreated with the `hzn dev dependency fetch` command used above.

Remember that the cpu2wiotp service is currently configured to mock the CPU usage service. Now we want to use the real data from the cpu service, so
remove the "MOCK" variable from the `userinput.json` file so that the default value of "false" will be used:
```
    "variables": {
        "PUBLISH": false,
        "SAMPLE_INTERVAL": 2,
        "SAMPLE_SIZE": 5,
        "VERBOSE": "1"
    }
```

Now start the cpu2wiotp service again, and this time both the services will be started.
```
hzn dev service start
```

After a few seconds, the log output from the cpu2wiotp service container can be found in syslog:
```
tail -f /var/log/syslog | grep cpu2wiotp
```

Notice that the cpu2wiotp service is now picking up real CPU usage samples from the cpu service and computing an average.
This is accomplished without any open ports on the host because the cpu2wiotp service container has joined the docker network of the cpu service. (This can be seen using the commands `docker ps`, `docker network ls`, and `docker inspect <id>`.)
This is how the containers will be configured when running on Edge nodes.

Control-C to stop `tail`. Then the cpu2wiotp service container and the cpu service container test environment can be stopped:
```
hzn dev service stop
```

Let's review what has been accomplished so far:
* A cpu service project has been authored which produces a standalone container that can be independently tested.
* A cpu2wiotp service project has been authored which produces a standalone container that can be independently tested.
* The cpu2wiotp service project has also been extended to depend on the cpu service project such that both projects can be run in the Horizon test environment, and tested together.

## Publish Data to Watson IoT Platform

This section shows you how to publish the CPU usage average to your Watson IoT Platform, using WIoTP's core-iot service.

Modify your project as follows:
* Update `userinput.json` under the `services` section to enable the cpu2wiotp service to publish the CPU average by setting PUBLISH to `true`:
```
    "variables": {
        "PUBLISH": true,
        "SAMPLE_INTERVAL": 2,
        "SAMPLE_SIZE": 5,
        "VERBOSE": "1"
    }
```
* Add the Wation IoT Platform core IoT service as a dependency of your cpu2wiotp service project.
```
hzn dev dependency fetch --url https://internetofthings.ibmcloud.com/wiotp-edge/services/core-iot --ver 2.4.0 -o IBM -a $ARCH -k /etc/horizon/trust/publicWIoTPEdgeComponentsKey.pem
```
* The command above would normally have added a section to the `userinput.json` file under `services` for the core-iot service, but it recognized that we already had that section.
* The dependency was added to the `requiredServices` array in the `service.definition.json` file.
* Note that the `deployment` section of `service.definition.json` already contains `environment` and `binds` settings appropriate for using the core-iot service, because we inherited that when we copied from `~/examples/edge/wiotp/cpu2wiotp`.
* Set up some things the core-iot service needs:
```
cp /etc/wiotp-edge/edge.conf.template /etc/wiotp-edge/edge.conf
mkdir -p /var/wiotp-edge/persist
wiotp_create_certificate -p $WIOTP_GW_TOKEN
```
* Start the cpu2wiotp service and look for the CPU messages in your Watson IoT Platform instance:
```bash
hzn dev service start
```

* Run `docker ps` and notice the additional core-iot containers that are running. (The core-iot service provides many other functions in addition to the one we are using: sending MQTT messages to WIoTP cloud.)

* Subscribe to the WIoTP cloud MQTT topic that the cpu2wiotp service is publishing to to see the cpu values:
```
mosquitto_sub -v -h $HZN_ORG_ID.messaging.$WIOTP_DOMAIN -p 8883 -i "$WIOTP_CLIENT_ID_APP" -u "$WIOTP_API_KEY" -P "$WIOTP_API_TOKEN" --capath /etc/ssl/certs -t iot-2/type/$WIOTP_GW_TYPE/id/$WIOTP_GW_ID/evt/status/fmt/json
```

* If you don't see messages coming to that, look again at the cpu2wiotp service log for errors. (Note: it is normal to sometimes get 1 error at the beginning of the cpu2wiotp log, depending on how long it takes each container to initialize.)
```
tail -f /var/log/syslog | grep cpu2wiotp
```

* Control-C to stop `tail`. Then stop the cpu2wiotp service:
```bash
hzn dev service stop
```

At this point you have successfully completed a quick pass through the development process.
You have developed and tested 2 services as standalone containers and running within the Horizon Edge node test environment.
You have also integrated your projects with your Watson IoT Platform instance and are able to publish data to it.
In order to make these projects available for other Edge nodes to run them, you need to publish your projects.
The next section describes how to use `hzn dev` to do that.

## Deploying the projects to WIoTP/Horizon

When you are satisfied that the services are working correctly, you can deploy them to the WIoTP so that any Edge node in your organization can run them.
The first step is to create a key pair that will be used to sign your docker images and the deployment configuration of your services.
```bash
cd ~/hzn
hzn key create <x509-org> <x509-cn>
export PRIVATE_KEY_FILE=~/hzn/*-private.key
export PUBLIC_KEY_FILE=~/hzn/*-public.pem
```
where `x509-org` is a company name or organization name that is suitable to be used as an x509 certificate organization name, and `x509-cn` is an x509 certificate common name (preferably an email address issued by the `x509-org` organization).

The private key will be used to sign the services, the public key is needed by any Edge node that wants to run the services, so that it can verify their signatures. (Horizon can take care of distributing the public key to the edge nodes that need it.)

You will be storing your docker images in Docker Hub, so login now:
```
docker login -u $DOCKER_HUB_ID
```

Now publish the cpu service:
```bash
cd ~/hzn/service/cpu
hzn exchange service publish -k $PRIVATE_KEY_FILE -K $PUBLIC_KEY_FILE -f horizon/service.definition.json
```

You can verify that the service was published:
```bash
hzn exchange service list
```

Publish the cpu2wiotp service:
```bash
cd ~/hzn/service/cpu2wiotp
hzn exchange service publish -k $PRIVATE_KEY_FILE -K $PUBLIC_KEY_FILE -f horizon/service.definition.json
```

You can verify that the cpu2wiotp service was published:
```bash
hzn exchange service list
```

Add the cpu2wiotp service to your gateway type deployment pattern and verify it is there:
```
hzn wiotp device remove $WIOTP_GW_TYPE $WIOTP_GW_ID; hzn wiotp type remove $WIOTP_GW_TYPE   # when edit works, we can do that instead
hzn wiotp type create $WIOTP_GW_TYPE $ARCH -s $MYDOMAIN-services-${CPU2WIOTP_NAME}_${CPU2WIOTP_VERSION}_$ARCH
hzn wiotp device create $WIOTP_GW_TYPE $WIOTP_GW_ID $WIOTP_GW_TOKEN
```

**If you are using the bluehorizon hybrid development environment, do not use the -s flag above, and also run these commands:**
```
hzn exchange node create -n "$HZN_DEVICE_ID:$WIOTP_GW_TOKEN"
hzn exchange pattern publish -p $WIOTP_GW_TYPE -f pattern/cpu2wiotp-service.json
hzn exchange agbot addpattern -o IBM stg-agbot-dal09-01.staging.bluehorizon.network $HZN_ORG_ID $WIOTP_GW_TYPE
hzn exchange agbot addpattern -o IBM stg-agbot-lon02-01.staging.bluehorizon.network $HZN_ORG_ID $WIOTP_GW_TYPE
hzn exchange agbot addpattern -o IBM stg-agbot-tok02-01.staging.bluehorizon.network $HZN_ORG_ID $WIOTP_GW_TYPE
```

## Using your cpu2wiotp service on edge nodes of this type

You, or others in your organization, can now use this cpu2wiotp service on many edge nodes. On each of those nodes:

* **If you have run thru this document before** on this edge node, do this to clean up:
```
hzn unregister -f
```

* Register the node and start the Watson IoT Platform core-IoT service and the cpu2wiotp service:
```
wiotp_agent_setup --org $HZN_ORG_ID --deviceType $WIOTP_GW_TYPE --deviceId $WIOTP_GW_ID --deviceToken "$WIOTP_GW_TOKEN"
```

**If you are using the bluehorizon hybrid development environment, register this way instead:**
```
hzn register -n "$HZN_DEVICE_ID:$WIOTP_GW_TOKEN" -f horizon/userinput.json $HZN_ORG_ID $WIOTP_GW_TYPE
```

Note: in the command above, we did not specify an input file with user input for your cpu2wiotp service, so all of the service's `userInput` variables are set to their default values. This is the preferred way to run your service, so it can easily be run on many edge nodes.

After a short while, usually within just a minute or two, agreements will be made to run the WIoTP core-iot service and your cpu2wiotp service. See [Register Your Edge Node](Edge-Quick-Start-Guide.md#register-your-edge-node) as a reminder for how to check for agreements and your service, and how to verify that CPU values are being sent to the cloud.

## Advanced topic - Expanding your project

Since your project Makefile and metadata files have environment variables in them for architecture, version, etc., you can easily test your services for different architectures, publish them for other gateway types, etc.

Here are a few things to consider when expanding your project:
* New versions of your services:
    * When you want to roll out a new version, build the new docker images (with a new tag for the new version) and create new service definitions with the new version number. Then replace the service reference in your gateway type pattern.
    * The edge nodes that are already using that pattern will automatically see the new service within a few minutes and re-negotiate the agreement to run the new version. You can manually force this by canceling the current agreement using `hzn agreement cancel <agreement-id>`
* Additional gateway types and architectures:
    * If you want multiple gateway types to run your service, you do not need to create the service definition multiple times. You only need to add your service reference to each gateway type pattern.
    * With WIoTP gateway types with Edge Services enabled, each type can only be used for a single architecture. So create a different gateway type for each architecture you want to support.

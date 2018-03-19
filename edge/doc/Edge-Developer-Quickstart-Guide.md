# Edge Developer Quickstart Guide

This Developer Quickstart Guide provides a simplified description of the process for developing, testing and deploying user-developed code in the Edge environment.
The [Edge Developer Guide](https://github.com/open-horizon/examples/wiki/Edge-Developer-Guide) is a more detailed description of the Edge environment and the concerns that an Edge developer has to be aware of.

Note there is a concise [Quick Start Guide](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md) available, that shows how to get an existing workload up and running on your edge nodes very quickly without having to develop any code. **That [Quick Start Guide](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md) is also a prerequisite for this guide.**

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

To familiarize yourself with WIoTP Edge, we suggest you go through the entire [Quick Start Guide](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md). But even if you do not go through that entire guide, **you must at least do the first sections of it, up to and including [Verify Your Gateway Credentials and Access](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md#verify-your-gateway-credentials-and-access)**, on the same edge node that you are using for this guide. (**For now, use the commented out line `aptrepo=testing` in the apt repo section, so you get the latest Horizon debian packages. They are currently required for this guide. You should have at least version 2.16.3**) That guide will have you accomplish the following necessary steps:

- Create your WIoTP organization, gateway type and id, and API key.
- Install docker, horizon, and some utilities.
- Set environment variables needed in the rest of this guide.
- Verify your edge node's access to the WIoTP cloud services.

After completing those steps in the [Quick Start Guide](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md) , continue here. Set this environment variable and get access to the examples repo:
```bash
export HZN_EXCHANGE_URL="https://$HZN_ORG_ID.internetofthings.ibmcloud.com/api/v0002/edgenode/"
cd ~
git clone https://github.com/open-horizon/examples.git
```

You will be storing your docker images in Docker Hub. If you don't already have an id, [sign up](https://hub.docker.com/sso/start/?next=https://hub.docker.com/) for one, set this environment variable and login:
```
export DOCKER_HUB_ID=<mydockerhubid>
docker login -u $DOCKER_HUB_ID
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

To accelerate development, use these prebuilt files:
```bash
cp ~/examples/edge/devguide/cpu/* .
cp ~/examples/edge/services/cpu_percent/*.sh .
```

The Makefile is setup to build the container, start it and run a simple test to ensure that it works.
1. Make the microservice container and run it locally:
```
make
```
1. Verify that the container was built and started successfully by looking for the output of the curl test at the end of the console output:
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

Now that the microservice container implementation is working correctly, we will introduce the Horizon project metadata in order to make this project testable in the Horizon test environment and to make it deployable to an Edge node.
The Horizon CLI contains a set of sub-commands within `hzn dev` that are useful for working with project metadata.
This guide will introduce them as the project is developed.

### Create microservice project metadata
This command creates the skeletal metadata files, in a `horizon` sub-directory, you will need to configure, test and eventually deploy the microservice.
```
hzn dev microservice new
```

These files were created in `horizon`:
* `microservice.definition.json` - the Horizon metadata of the microservice, the container it runs, the user inputs that have to be configured, etc.
* `userinput.json` - the runtime input values specified by the edge node owner for the microservice or the Horizon agent. You can also set these during development to test your microservice with various values.
* `horizon/dependencies` - This directory holds metadata describing other projects that this project depends on.
There are none for this project.


### Update microservice project metadata
Update the Horizon metadata files based on our example project. Modify the `microservice.definition.json` file:
1. Update the value of `specRef` with a URL that is unique to your organization, for example: `http://my.company.com/microservices/cpu`.
This URL will be used by other parts of the Edge system to refer to this microservice.
2. Update the value of `version` with a version number that complies with the [OSGI Version standard](https://www.osgi.org/wp-content/uploads/SemanticVersioning.pdf), e.g. 0.0.1 as your first version. It is simplest if you set this to the same version number used in the docker image tag.
3. Update `userInput` with any runtime arguments that the Edge node owner can (or must, if no default) set in order for the service to work correctly.
In this case, there are none so you can remove the skeletal array element so it is just `"userInput": [],`.
3. Update the `deployment` field under the `workloads` section with the docker container configuration.
 * In this case, change the empty service name inside the services map to a meaningful service name, e.g. "cpu".
This name will be used as the network domain name of the microservice container and will be used by workload containers when contacting it.
```
    "deployment": {
        "services": {
            "cpu": {
```
 * Also change the value of `image` to the docker image path (with tag) of the container that you built previously.
You can get the image name from the `docker images | grep cpu` command.
 * Set any environment variables in the `environment` array. These are variables and values that you want Horizon to pass into the container when it is started.
The Edge node owner cannot override these variables.
In this case, there are none so you can set the environment field to an empty array.

Your `microservice.definition.json` should look something like this:
```
    {
        "org": "<your_org>",
        "label": "",
        "description": "",
        "public": true,
        "specRef": "http://my.company.com/microservices/cpu",
        "version": "0.0.1",
        "arch": "amd64",
        "sharable": "multiple",
        "downloadUrl": "not used yet",
        "matchHardware": {},
        "userInput": [],
        "workloads": [
            {
                "deployment": {
                    "services": {
                        "cpu": {
                            "image": "cpu_microservice:0.0.1",
                            "privileged": false,
                            "environment": []
                        }
                    }
                },
                "deployment_signature": "",
                "torrent": ""
            }
        ]
    }
```

Update the `userinput.json` file for our example project:
1. Update the `global` section with Horizon agent attributes. See the [Attribute Documentation](https://github.com/open-horizon/anax/blob/master/doc/attributes.md) for a description of supported Horizon attributes.
In this case, there aren't any so you can remove the skeletal array element.
1. Update the `url` in the `microservices` section with the same value you set on `specRef` in the `microservice.definition.json` file.
The `microservices` section tells the test simulator which microservice to run.
In this case, we have only 1 microservice in our project.
1. Update the `variables` in the `microservices` section with values for user input variables that are defined in the `microservice.definition.json` file.
In this case, there are none, so you can remove the skeletal variable setting.
1. The `versionRange` indicates which version of the microservice to run.
By default, the value is set to all versions so it will always execute the microservice in this project.

Your `userinput.json` file should look something like this:
```
    {
        "global": [],
        "microservices": [
            {
                "org": "<your_org>",
                "url": "http://my.company.com/microservices/cpu",
                "versionRange": "[0.0.0,INFINITY)",
                "variables": {}
            }
        ]
    }
```

### Verify microservice project metadata
Now that you have defined and configured the Horizon microservice, verify that the project has no errors in it.
```
hzn dev microservice verify
```

If the verify sub-command finds any inconsistencies or errors it will report the file and the location within the file where the error was detected.

### Test the microservice project
Test your microservice in the Horizon test environment.
The test environment closely simulates the Edge node environment in which your microservice will run when deployed to an Edge node.
```
hzn dev microservice start
```
After a few seconds, the microservice will be started.
The `docker ps` command will display the running container.
Notice that the microservice docker container names are derived from the microservice definition's `specRef` URL, `version`, a unique UUID, and the service name.

The microservice is attached to a docker network. See it via `docker network ls`.
The network has nearly the same name as the microservice container name (minus the service name).
Also note that the running container exposes no host network ports.
Microservices on a Horizon Edge node, run in a sandboxed environment where network access is restricted only to workloads that require the microservice.
This means that when running in a Horizon test environment, by default, the host cannot access the microservice's network ports.
This is the desired and expected behavior.

Finally, use
```
docker inspect <container_name> | jq ".[0].Config.Env"
```
to inspect the environment variables that have been passed into the container.
The variables prefixed with "HZN" are provided by the Horizon Edge node platform.
A microservice container can exploit these environment variables.
See the [Horizon Environment Variable Documentation](https://github.com/open-horizon/anax/blob/master/doc/managed_workloads.md) for a complete description of Horizon Edge node environment variables.

The microservice container test environment can be stopped (it will be started again later when the workload requires it):
```
hzn dev microservice stop
```

## Create a workload project

Workloads use microservices to gain access to data on the Edge node.
Let's develop a minimal workload and explore more complex usages of the `hzn dev` sub-commands.

On your development machine, create a project directory:
```bash
mkdir -p ~/hzn/workload/cpu2wiotp/horizon
cd ~/hzn/workload/cpu2wiotp
```

Next, create a docker container that computes averages for a set of CPU usage samples.
You would:
* Author a Dockerfile to hold the container definition, for example `~/hzn/workload/cpu2wiotp/Dockerfile`.
* Author code to implement the workload logic, for example `~/hzn/workload/cpu2wiotp/workload.sh`.
* Author a makefile to build and test the workload as a standalone container, for example `~/hzn/workload/cpu2wiotp/Makefile`.

To accelerate development, use these prebuilt files:
```bash
cp ~/examples/edge/devguide/cpu2wiotp/* .
cp ~/examples/edge/wiotp/cpu2wiotp/workload.sh .
cp ~/examples/edge/wiotp/cpu2wiotp/pattern/insert-cpu2wiotp-template.json horizon
```

The Makefile is setup to build the container, start it, and run a simple test to ensure that it works.
1. Make the workload container and run it locally:
```
make
```
1. Verify that the container was built and started successfully by looking at the docker logs output at the bottom of the console.
It will look something like this:
```
Starting infinite loop to read from microservice then publish...
 Interval 1 cpu: 53
 Interval 2 cpu: 54
```
1. Control-C to stop `docker logs` and then stop the running container:
```
make stop
```

Now that the workload container implementation is working correctly, we will introduce the Horizon project metadata in order to make this project testable in the Horizon test environment and to make it deployable to an Edge node.

### Create workload project metadata

Create the workload project metadata.
This command will create the skeletal metadata files you will need to configure, test, and eventually deploy the workload.
```
hzn dev workload new
```

As with the microservice, a new directory called `horizon` was created in your project and it holds several files and a directory:
* `workload.definition.json` - This file holds the Horizon metadata of the workload, the container it runs, the user inputs that have to be configured, etc.
* `userinput.json` - The input values specified by the edge node owner for the workload or the Horizon agent.
* `horizon/dependencies` - This directory holds metadata describing other projects that this project depends on.
We will add a dependency later in this guide.

### Update workload project metadata

Update the Horizon metadata files based on our example project. Modify the `workload.definition.json` file:
1. Update the value of `workloadUrl` with a URL that is unique to your organization, for example: `http://my.company.com/workloads/cpu2wiotp`.
This URL will be used by other parts of the Edge system to refer to this workload.
1. Update the value of `version` with a version number that complies with the [OSGI Version standard](https://www.osgi.org/wp-content/uploads/SemanticVersioning.pdf), e.g. 0.0.1 as your first version. It is simplest if you set this to the same version number used in the docker image tag.
1. Update `userInput` with with any runtime arguments that the Edge node owner can (or must, if no default) set in order for the workload to work correctly.
In this example, there are 5 variables that condition the logic of the workload.
Define these variables within the `userInput` section of the `workload.definition.json` file:
 * `SAMPLE_SIZE` - the number of samples to include in the average. Set the type to "int" and the default to "6".
 * `SAMPLE_INTERVAL` - the delay between samples. Set the type to "int" and the default to "5".
 * `MOCK` - used to provide mock samples. Set the type to "boolean" and the default to "false".
 * `PUBLISH` - used to control whether or not the CPU average is publish to Watson IoT Platform. Set the type to "boolean" and set the defauilt to "true".
 * `VERBOSE` - used to provide detailed logging of the workload logic. Set the type to "string" and the default to "0".
1. Update the `deployment` field under the `workloads` section with the docker container configuration.
 * In this case, change the empty service name inside the services map to a meaningful service name, e.g. "cpu2wiotp".
 * Also change the value of `image` to the docker image path (with tag) of the container that you built previously.
 * Set any environment variables in the `environment` array.
These are variables and values that you want Horizon to pass into the container when it is started, but that the Edge node configuration cannot override.
In this case, there are no environment variables that need to be set so you can set the environment field to an empty array.

Your `workload.definition.json` should look something like this:
```
    {
        "org": "<your_org>",
        "label": "",
        "description": "",
        "public": true,
        "workloadUrl": "http://my.company.com/workloads/cpu2wiotp",
        "version": "0.0.1",
        "arch": "amd64",
        "downloadUrl": "not used yet",
        "apiSpec": [],
        "userInput": [
            {
                "name": "SAMPLE_SIZE",
                "label": "the number of samples before calculating the average",
                "type": "int",
                "defaultValue": "6"
            },
            {
                "name": "SAMPLE_INTERVAL",
                "label": "the number of seconds between samples",
                "type": "int",
                "defaultValue": "5"
            },
            {
                "name": "MOCK",
                "label": "mock the CPU sampling",
                "type": "boolean",
                "defaultValue": "false"
            },
            {
                "name": "PUBLISH",
                "label": "publish the CPU samples to WIoTP",
                "type": "boolean",
                "defaultValue": "true"
            },
            {
                "name": "VERBOSE",
                "label": "log everything that happens",
                "type": "string",
                "defaultValue": "0"
            }
        ],
        "workloads": [
            {
                "deployment": {
                    "services": {
                        "cpu2wiotp": {
                            "image": "cpu2wiotp_workload:0.0.1",
                            "privileged": false,
                            "environment": []
                        }
                    }
                },
                "deployment_signature": "",
                "torrent": ""
            }
        ]
    }
```

Modify the `userinput.json` file for our example project.
1. We again do not have any Horizon agent attributes to set, so remove the skeletal array element from the `global` section.
1. Update the `url` in the `workloads` section with the same value you set in `workloadUrl` in the `workload.definition.json` file.
In this case, we have only 1 workload in our project.
1. Update the `variables` in the `workloads` section with values for user input variables that are defined in the `workload.definition.json` file.
In this case, there are five.
Even though these variables have default values, we need to override them with values suitable for testing, because the defaults are not what we want at this stage of the project.

Your `userinput.json` file should look something like this:
```
    {
        "global": [],
        "workloads": [
            {
                "org": "<your_org>",
                "url": "http://my.company.com/workloads/cpu2wiotp",
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
    }
```

### Verify workload project metadata

Now that you have defined and configured the Horizon workload, verify that the project has no errors in it.
```
hzn dev workload verify
```

If the verify sub-command finds any inconsistencies or errors it will report the file and the location within the file were the error was detected.

### Test the workload project

Test your workload in the Horizon test environment.
The test environment closely simulates the Edge node environment in which your workload will run when deployed to an Edge node.
```
hzn dev workload start
```
After a few seconds, the workload will be started.
The `docker ps` command will display the running container.
Notice that the workload docker container name is derived from the hexadecimal agreement id and the service name.
An agreement id is a hexadecimal string that uniquely identifies the agreement.
When your workload is eventually deployed by an agbot, an agreement id will be assigned and it will appear similarly to how you see it now.

The workload is attached to a docker network. See it via `docker network ls`.
The network has nearly the same name as the workload container name (minus the service name).

Finally, use
```
docker inspect <container_name> | jq '.[0].Config.Env'
```
to inspect the environment variables that have been passed into the container.
Variables prefixed with "HZN" are provided by the Horizon Edge node platform.
A workload container can exploit the "HZN" environment variables.
See the [Environment Variable Documentation](https://github.com/open-horizon/anax/blob/master/doc/managed_workloads.md) for a complete description of Horizon Edge node variables.

The log output from the workload container can be found in syslog:
```
tail -f /var/log/syslog | grep cpu2wiotp
```

The workload container test environment can be stopped:
```
hzn dev workload stop
```

## Add a workload dependency

So far, the workload is executing with mocked CPU usage data, not very interesting.
Now let's update the workload project to use the CPU usage microservice as a dependency.
We are going to use the local microservice project we created earlier as the source of the dependency.
Notice that the -p flag in the following command points to the `horizon` directory where the microservice project's metadata resides.
```bash
hzn dev dependency fetch -p ~/hzn/ms/cpu/horizon
```

The addition of the new dependency has updated several pieces of metadata in your workload project.
Notice that:
* `horizon/dependencies` now contains a copy of the dependency's microservice definition.
* the dependency was added to the `apiSpec` array in the `workload.definition.json` file.
This array holds references to microservice dependencies.
* the `userinput.json` file was updated to include the new dependency for the test environment.

In general, if the Horizon metadata of a dependent project changes, the dependency must be recreated with the `hzn dev dependency` command used above.

Remember that the workload is currently configured to mock the CPU usage microservice.
Obviously we don't want to do that any longer since the real service is available to the workload project.
Remove the "MOCK" variable from the `userinput.json` file so that the default value of "false" will be used.
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
This is accomplished without any open ports on the host because the workload container has joined the docker network of the microservice.
This is how the containers will be configured when running on an Edge node.

The workload container and the microservice container test environment can be stopped:
```
hzn dev workload stop
```

Let's review what has been accomplished so far:
* A microservice project has been authored which produces a standalone container that can be independently tested.
* A workload project has been authored which produces a standalone container that can be independently tested.
* The workload project has also been extended to depend on the microservice project such that both projects can be run in the Horizon test environment, and tested together.

## Publish Data to Watson IoT Platform

This section shows you how to publish the CPU usage average to your Watson IoT Platform.

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
* Save off this version of `userinput.json`, because you'll need it later:
```
cp horizon/userinput.json horizon/userinput-without-core-iot.json
```
* Configure HTTPS authentication in the `global` section of the `userinput.json` file (so that the WIoTP containers can be fetched and loaded). Fill in your values for the environment variables:
```
    {
        "type": "HTTPSBasicAuthAttributes",
        "sensor_urls": [
            "https://us.internetofthings.ibmcloud.com/api/v0002/horizon-image/common"
        ],
        "publishable": false,
        "host_only": true,
        "variables": {
            "username": "$HZN_ORG_ID/$HZN_DEVICE_ID",
            "password": "$WIOTP_GW_TOKEN"
        }
    }
```
* Add the Wation IoT Platform core IoT microservice as a dependency of your workload project.
```
hzn dev dependency fetch -s https://internetofthings.ibmcloud.com/wiotp-edge/microservices/edge-core-iot-microservice --ver 2.4.0 -o IBM -a amd64 -k /etc/horizon/trust/publicWIoTPEdgeComponentsKey.pem
```
* The command above added a section to the `userinput.json` file under `microservices` for the edge-core-iot-microservice.
Fill in the variable values like this (substitute for the env vars):
```
    "variables": {
        "WIOTP_DEVICE_AUTH_TOKEN": "$WIOTP_GW_TOKEN",
        "WIOTP_DOMAIN": "$HZN_ORG_ID.messaging.internetofthings.ibmcloud.com"
    }
```
* Update `workload.definition.json` with new deployment environment variables:
```
    "environment": [
        "WIOTP_EDGE_MQTT_IP=edge-connector",
        "WIOTP_DOMAIN=internetofthings.ibmcloud.com",
        "WIOTP_PEM_FILE=/var/wiotp-edge/persist/dc/ca/ca.pem"
    ],
```
* Update `workload.definition.json` with a file binding to get access to the WIoTP MQTT Broker certificate.
Add a new field called `binds` under `workloads.deployment.services.cpu2wiotp`.
It is a peer to the `environment` field:
```
    "binds": [
        "/var/wiotp-edge:/var/wiotp-edge"
    ],
```
* Start the workload and look for the CPU messages in your Watson IoT Platform instance:
```bash
hzn dev workload start
```

* Subscribe to the WIoTP cloud MQTT topic that the workload is publishing to to see the cpu values:
```
mosquitto_sub -v -h $HZN_ORG_ID.messaging.$WIOTP_DOMAIN -p 8883 -i "$WIOTP_CLIENT_ID_APP" -u "$WIOTP_API_KEY" -P "$WIOTP_API_TOKEN" --capath /etc/ssl/certs -t iot-2/type/$WIOTP_GW_TYPE/id/$WIOTP_GW_ID/evt/status/fmt/json
```

* If you don't see messages coming to that, look again at the workload log for errors. (Note: it is normal to sometimes get 1 error at the beginning of the cpu2wiotp log, depending on how long it takes each container to initialize.)
```
tail -f /var/log/syslog | grep cpu2wiotp
```

* Stop the workload:
```bash
hzn dev workload stop
```

At this point you have successfully completed a quick pass through the development process.
You have developed and tested a microservice and a workload as standalone containers and running within the Horizon Edge node test environment.
You have also integrated your projects with your Watson IoT Platform instance and are able to publish data to it.
In order to make these projects available for other Edge nodes to run them, you need to publish your projects.
The next section describes how to use `hzn dev` to do that.

## Deploying the projects to WIoTP/Horizon
**Note: this section is still under development. It should work as-is, but we are still testing and refining it.**

When you are satisfied that the microservice and workload are working correctly, you can deploy them to the Edge so that any Edge node in your organization can run them.
The first step is to create a key pair that will be used to sign the deployment configuration of your microservice and workload.
```bash
cd ~/hzn
hzn key create <x509-org> <x509-cn>
export PRIVATE_KEY_FILE="~/hzn/*-private.key"
export PUBLIC_KEY_FILE="~/hzn/*-public.pem"
```
where `x509-org` is a company name or organization name that is suitable to be used as an x509 certificate organization name, and `x509-cn` is an x509 certificate common name (preferably an email address issued by the `x509-org` organization).

This command will generate a private key and a public key.
The private key will be used to sign the microservice and workload, the public is needed by any Edge node that wants to run the microservice and workload, so that it can verify their signatures.

Now publish the microservice:
```bash
cd ~/hzn/ms/cpu
hzn dev microservice publish -k $PRIVATE_KEY_FILE
```

You can verify that the microservice was published:
```bash
hzn exchange microservice list
```

When the microservice is successfully published, upload your microservice image to the docker registry that you are using. The microservice project has this built into the makefile:
```bash
make publish
```

Publish the workload:
```bash
cd ~/hzn/workload/cpu2wiotp
hzn dev workload publish -k $PRIVATE_KEY_FILE
```

You can verify that the workload was published:
```bash
hzn exchange workload list
```

When the workload is successfully published, upload your workload container to the docker registry that you are using. The workload project has this built into the makefile:
```bash
make publish
```

Add this workload to your gateway type deployment pattern:

* Replace environment variables in `horizon/insert-cpu2wiotp-template.json` with your values:
```
envsubst < horizon/insert-cpu2wiotp-template.json > horizon/insert-cpu2wiotp.json
```

* Replace these 2 lines in `horizon/insert-cpu2wiotp.json` with your values:
```
      "workloadUrl": "https://internetofthings.ibmcloud.com/workloads/cpu2wiotp",
          "version": "",
```

* Add your workload to your gateway type deployment pattern and verify it is there:
```
hzn exchange pattern insertworkload -k $PRIVATE_KEY_FILE -f horizon/insert-cpu2wiotp.json $WIOTP_GW_TYPE
hzn exchange pattern list $WIOTP_GW_TYPE | jq .
```

## Using your workload on edge nodes of this type

You, or others in your organization, can now use this workload on many edge nodes. On each of those nodes:

* Copy your public key, and `~/hzn/workload/cpu2wiotp/horizon/userinput.json` to it.

* Import your to the Horizon agent:
```bash
hzn key import -k $PUBLIC_KEY_FILE
```

* **If you have run thru this document before** on this edge node, do this to clean up:
```
hzn unregister -f
```

* Register the node and start the Watson IoT Platform core-IoT service and the CPU workload:
```
wiotp_agent_setup --org $HZN_ORG_ID --deviceType $WIOTP_GW_TYPE --deviceId $WIOTP_GW_ID --deviceToken "$WIOTP_GW_TOKEN" -f ~/hzn/workload/cpu2wiotp/horizon/userinput-without-core-iot.json
```

After a short while, usually within just a minute or two, agreements will be made to run the WIoTP core-iot service and your workload. See [Register Your Edge Node](https://github.com/open-horizon/examples/blob/master/edge/doc/Edge-Quick-Start-Guide.md#register-your-edge-node) as a reminder for how to check for agreements and your workload, and how to verify that CPU values are being sent to the cloud.

## Advanced topic - Expanding your project

Your project metadata currently has hardcoded values for architecture, gateway type, and workload version. You can easily make your project more flexible by parameterizing your project metadata files by using environment varibles in your files and processing your files with a tool like `envsubst`.
If you then add a few new recipes to your `Makefile` you can automate the process of creating the metadata files needed by horizon.
See `~/examples/edge/wiotp/cpu2wiotp/Makefile` for an example of how to do this (specifically the targets `hznbuild` and `hznstart`).

Here are a few things to consider when expanding your project:
* New versions of your microservice and workload:
    * When you want to roll out a new version, build the new docker images (with a new tag for the new version) and create new microservice/workload definitions with the new version number. Then replace the workload reference in your gateway type pattern.
    * The edge nodes that are already using that pattern will automatically see the new workload within a few minutes and re-negotiate the agreement to run the new version. You can manually force this by canceling the current agreement using `hzn agreement cancel <agreement-id>`
* Additional gateway types and architectures:
    * If you want multiple gateway types to run your workload, you do not need to create the microservice/workload definitions multiple times. You only need to add your workload reference to each gateway type pattern.
    * With WIoTP gateway types with Edge Services enabled, each type can only be used for a single architecture. So create a different gateway type for each architecture you want to support.

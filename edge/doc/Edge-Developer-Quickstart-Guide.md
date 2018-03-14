# Edge Developer Quickstart Guide

This Developer Quickstart Guide provides a simplified description of the process for developing, testing and deploying user-developed code in the Edge environment.
The [Edge Developer Guide](https://github.com/open-horizon/examples/wiki/Edge-Developer-Guide) is a more detailed description of the Edge environment and the concerns that an Edge developer has to be aware of.

Note there is a concise [Quick Start Guide](https://github.com/open-horizon/examples/wiki/Edge-Quick-Start-Guide) available, that shows how to get an existing workload up and running on your edge nodes very quickly without having to develop any code. **That [Quick Start Guide](https://github.com/open-horizon/examples/wiki/Edge-Quick-Start-Guide) is also a prerequisite for this guide.**

Additional information is available, and questions may be asked, in our forum, at:
* [https://discourse.bluehorizon.network/](https://discourse.bluehorizon.network/)

The Edge is based upon the open source Horizon project [https://github.com/open-horizon](https://github.com/open-horizon). There are therefore several references to Horizon, and the "hzn" Linux shell command in this document.

Edge simplifies and secures the global deployment and maintenance of software on IoT edge nodes.
This document will guide you through the process of building, testing, and deploying IoT edge software, using Watson IoT to securely deploy and then fully manage the software on your IoT edge nodes all over the world.
IoT software maintenance in Edge with Watson IoT becomes fully automatic (zero-touch for your edge nodes), highly secure and easy to centrally control.

## Overview

This guide is intended for developers who want to experience the Edge software development process with Watson IoT Platform, using a simple example.
In this guide you will learn how to create Horizon microservices and workloads, how to test them and ultimately how to integrate them with the Watson IoT Platform.

As you progress through this guide, you will first build a simple microservice that extracts CPU usage information from the underlying machine.
Then you will build a workload that samples CPU usage information from the microservice, computes an average and then publishes the average to Waton IoT Platform.
Along the way, you will be exposed to many concepts and capabilities that are documented in complete detail in the [Edge Developer Guide](https://github.com/open-horizon/examples/wiki/Edge-Developer-Guide).

## Before you begin

To familiarize yourself with WIoTP Edge, we suggest you go through the entire [Quick Start Guide](https://github.com/open-horizon/examples/wiki/Edge-Quick-Start-Guide). But even if you do not go through that entire guide, **you must at least do the first sections of it, up to and including [Verify Your Gateway Credentials and Access](https://github.com/open-horizon/examples/wiki/Edge-Quick-Start-Guide#verify-your-gateway-credentials-and-access)**, on the same edge node that you are using for this guide. (For now use the commented out line `aptrepo=testing` in the apt repo section.) That will have you accomplish the following necessary steps:

- Create your WIoTP organization, gateway type and id, and API key.
- Install docker, horizon, and some utilities.
- Set environment variables needed in the rest of this guide.
- Verify your edge node's access to the WIoTP cloud services.

Now set this environment variable and get access to the examples repo:
```bash
export HZN_EXCHANGE_URL="https://$HZN_ORG_ID.internetofthings.ibmcloud.com/api/v0002/edgenode/"
cd ~
https://github.com/open-horizon/examples.git
```

## Create a microservice project

We will start the microservice project by creating a docker container and implementation that will respond to an HTTP request with CPU usage information.
The microservice implementation is a very simple bourne shell script.
On your development machine, create a project directory:
```bash
mkdir -p ~/hzn/ms/cpu
cd ~/hzn/ms/cpu
```

Next, create a docker container that exposes an HTTP API for obtaining CPU usage information.
You would:
* Author a Dockerfile to hold the container definition, for example ~/hzn/ms/cpu/Dockerfile.
* Author a shell script that will run when the container starts, for example ~/hzn/ms/cpu/start.sh.
* Author the shell script that runs when the TCP listener gets a message, for example ~/hzn/ms/cpu/service.sh.
* Author a Makefile that will build and test this container on your development machine, for example ~/hzn/ms/cpu/Makefile.
The makefile is kept as simple as possible for the purposes of illustration.
Later in this guide, you will see how to create something more complex.

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
The Horizon CLI contains a set of sub-commands within `hzn dev` that are useful for working with project metadata data.
This guide will introduce them as the project is developed.

### Create microservice project metadata
Create the microservice project metadata.
This command will create the skeletal metadata files you will need to configure, test and eventually deploy the microservice.
```
hzn dev microservice new
```

The `hzn dev microservice verify` command can be used to learn which parts of your Horizon metadata need to be adjusted so that the project becomes testable in the Horizon test environment.
This guide will walk you through the updates you need to make to get an example running.
You can use this command to verify your changes as you proceed.

You will notice that a new directory called `horizon` has been created in your project and it holds several files and one directory.
The `horizon` directory is the default directory for project metadata.
The default can be overridden by using the -d flag.
* `microservice.definition.json` - This file holds the Horizon definition of the microservice, the container it runs, the user inputs that have to be configured, etc.
* `userinput.json` - This file holds configuration data that describes how the microservice's Horizon attributes should be configured and values for user inputs.
Think of this file as the inputs needed to configure the microservice for a test.
* `horizon/dependencies` - This directory holds metadata describing other projects that this project depends on.
There are none for this project.


### Update microservice project metadata
Update the Horizon metadata files based on our project.
Use an editor of your choice to modify the `microservice.definition.json` file:
1. Update the value of `specRef` with a URL that is unique to your organization, for example: `http://my.company.com/microservices/cpu`.
This URL will be used by other parts of the Edge system to refer to this microservice.
2. Update the value of `version` with a version number that complies with the [OSGI Version standard](https://www.osgi.org/wp-content/uploads/SemanticVersioning.pdf), e.g. 0.0.1 as your first version.
3. Update `userInput` with any configuration variables that the microservice will need.
These are variables that the Edge node configuration can (or must) set in order for the service to work correctly.
In this case, there are none so you can remove the skeletal array element.
3. Update `workloads` with the service and docker metadata that comprises the microservice implementation.
 * In this case, change the empty service name inside the services map to a meaningful service name, e.g. "cpu".
This name is the network domain name used by workload containers to access remoteable APIs offered by this microservice, so you will
need to use this same name when the workload implementation invokes the microservice.
```
    "deployment": {
        "services": {
            "cpu": {
```
 * Also change the value of `image` to the docker image name (with tag) of the container that we built previously.
You can get the image name from the `docker images` command.
 * Set any environment variables in the `environment` array. These are variables and values that you want Horizon to pass into the container when it is started.
The Edge node configuration cannot override these variables.
In this case, there are none so you can set the environment field to an empty array.

Your `microservice.definition.json` should look something like the following.
Please note that `<your_org>` will have the value of the HZN_ORG_ID environment variable:
```
    {
        "org": "<your_org>",
        "label": "",
        "description": "",
        "public": false,
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

Use an editor of your choice to modify the `userinput.json` file.
As mentioned earlier, the purpose of this file is to configure the microservice as it might be configured on an Edge node, enabling you to test your microservice in a simulated Horizon test environment.
In order to do this, update the `userinput.json` file with the Edge node configuration that you would like to use:
1. Update the `global` section with Horizon attributes. See the [Attribute Documentation](https://github.com/open-horizon/anax/blob/master/doc/attributes.md) for a description of supported Horizon attributes.
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

If the verify sub-command finds any inconsistencies or errors it will report the file and the location within the file were the error was detected.

### Test the microservice project
Test your microservice in the Horizon test environment.
The test environment closely simulates the Edge node environment in which your microservice will run when deployed to an Edge node.
```
hzn dev microservice start
```
After a few seconds, the microservice will be started.
The `docker ps` command will display the running container.
Notice that the microservice docker instance names are derived from the microservice definition's `specRef` URL, `version`, a unique UUID, and the service name.

The microservice is attached to a docker network: `docker network ls`.
The network has nearly the same name as the microservice container instance (minus the service name).
Also note that the running container exposes no host network ports.
Microservices on a Horizon Edge node, run in a sandboxed environment where network access is restricted only to workloads that require the microservice.
This means that when running in a Horizon test environment, by default, the host cannot access the microservice container's network endpoints.
That is, `make check` fails.
This is the desired and expected behavior.

Finally, use
```
docker inspect <instance_name> | jq ".[0].Config.Env"
```
to inspect the environment variables that have been passed into the container implementation.
The variables prefixed with "HZN" are provided by the Horizon Edge node platform.
A microservice container implementation can exploit these environment variables.
See the [Environment Variable Documentation](https://github.com/open-horizon/anax/blob/master/doc/managed_workloads.md) for a complete description of Horizon Edge node environment variables.

The microservice container test environment can be stopped:
```
hzn dev microservice stop
```

## Create a workload project

Workloads use microservices to gain access to data on the Edge node.
Let's develop a minimal workload and explore more complex usages of the `hzn dev` sub-commands.

On your development machine, create a project directory:
```bash
mkdir -p ~/hzn/workload/cpu2wiotp
cd ~/hzn/workload/cpu2wiotp
```

Next, create a docker container that computes averages for a set of CPU usage samples.
You would:
* Author a Dockerfile to hold the container definition, for example ~/hzn/workload/cpu2wiotp/Dockerfile.
* Author a shell script to implement the workload logic, for example ~/hzn/workload/cpu2wiotp/workload.sh.
* Author a makefile to build and test the workload as a standalone container, for example ~/hzn/workload/cpu2wiotp/Makefile.
The makefile is kept as simple as possible for the purposes of illustration.
Later in this guide, you will see how to create something more complex.

To accelerate development, use these prebuilt files:
```bash
cp ~/examples/edge/devguide/cpu2wiotp/* .
cp ~/examples/edge/wiotp/cpu2wiotp/workload.sh .
```

The Makefile is setup to build the container, start it and run a simple test to ensure that it works.
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
1. Stop the running container:
```
make stop
```

Now that the workload container implementation is working correctly, we will introduce the Horizon project metadata in order to make this project testable in the Horizon test environment and to make it deployable to an Edge node.

### Create workload project metadata

Create the workload project metadata.
This command will create the skeletal metadata files you will need to configure, test and eventually deploy the workload.
```
hzn dev workload new
```

The `hzn dev workload verify` command can be used to learn which parts of your Horizon metadata need to be adjusted so that the project becomes testable in the Horizon test environment.
This guide will walk you through the updates you need to make to get the example running.
You can use this command to verify your changes as you proceed.

You will notice that a new directory called `horizon` has been created in your project and it holds several files and a directory.
* `workload.definition.json` - This file holds the Horizon definition of the workload, the container it runs, the user inputs that have to be configured, etc.
* `userinput.json` - This file holds configuration data that describes how the workload's Horizon attributes should be configured and values for user inputs.
Think of this file as the inputs needed to configure the workload for a test.
* `horizon/dependencies` - This directory holds metadata describing other projects that this project depends on.
We will add a dependency later in this guide.

### Update workload project metadata

Update the Horizon metadata files based on our project.
Use an editor of your choice to modify the `workload.definition.json` file:
1. Update the value of `workloadUrl` with a URL that is unique to your organization, for example: `http://my.company.com/workloads/cpu2wiotp`.
This URL will be used by other parts of the Edge system to refer to this workload.
1. Update the value of `version` with a version number that complies with the [OSGI Version standard](https://www.osgi.org/wp-content/uploads/SemanticVersioning.pdf), e.g. 0.0.1 as your first version.
1. Update `userInput` with any configuration variables that the workload will need.
These are variables that the Edge node configuration can set (or must set if no default is provided) in order for the workload logic to work correctly.
In this case, there are 5 variables that condition the logic of the workload.
Define these variables within the `userInput` section of the `workload.definition.json` file:
 * `SAMPLE_SIZE` - the number of samples to include in the average. Set the type to "int" and the default to "6".
 * `SAMPLE_INTERVAL` - the delay between samples. Set the type to "int" and the default to "5".
 * `MOCK` - used to provide mock samples. Set the type to "boolean" and the default to "false".
 * `PUBLISH` - used to control whether or not the CPU average is publish to Watson IoT Platform. Set the type to "boolean" and set the defauilt to "true".
 * `VERBOSE` - used to provide detailed logging of the workload logic. Set the type to "string" and the default to "0".
1. Update `workloads` with the service and docker metadata that comprises the workload implementation.
 * In this case, change the empty service name inside the services map to a meaningful service name, e.g. "cpu2wiotp".
 * Also change the value of `image` to the docker image name (with tag) of the container that we built previously.
 * Set any environment variables in the `environment` array.
These are variables and values that you want Horizon to pass into the container when it is started, but that the Edge node configuration cannot override.
In this case, there are no environment variables that need to be set so you can set the environment field to an empty array.

Your `workload.definition.json` should look something like this.
Please note that `<your_org>` will have the value of the HZN_ORG_ID environment variable:
```
    {
        "org": "<your_org>",
        "label": "",
        "description": "",
        "public": false,
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
                "defaultValue": "10"
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
                            "environment": [],
                            "image": "cpu2wiotp_workload:0.0.1",
                            "privileged": false
                        }
                    }
                },
                "deployment_signature": "",
                "torrent": ""
            }
        ]
    }
```

Use an editor of your choice to modify the `userinput.json` file.
As mentioned earlier, the purpose of this file is to configure the workload as it might be configured on an Edge node, enabling you to test your workload in a simulated Horizon test environment.
In order to do this, update the `userinput.json` file with the Edge node configuration that you would like to use:
1. Update the `global` section with Horizon attributes. See the [Attribute Documentation](https://github.com/open-horizon/anax/blob/master/doc/attributes.md) for a description of supported Horizon attributes.
In this case, there aren't any so you can remove the skeletal array element.
1. Update the `url` in the `workloads` section with the same value you set on `workloadUrl` in the `workload.definition.json` file.
The `workloads` section tells the test simulator which workload to run.
In this case, we have only 1 workload in our project.
1. Update the `variables` in the `workloads` section with values for user input variables that are defined in the `workload.definition.json` file.
In this case, there are five.
Even though these variables have default values, we need to override them with values suitable for testing, because the defaults are not what we want at this stage of the project.
1. The `versionRange` indicates which version of the workload to run.
By default, the value is set to all versions so it will always execute the workload in this project.

Your `userinput.json` file should look something like this:
```
    {
        "global": [],
        "microservices": [],
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
Notice that the workload docker instance names are derived from the hexadecimal agreement id and the service name.
An agreement id is simply a hexadecimal encoded very large number that uniquely identifies the agreement.
When your workload runs as a result of an agreement between it and an agbot, an agreement id will be assigned and it will appear similarly to how you see it now.

The workload is attached to a docker network: `docker network ls`.
The network has nearly the same name as the workload container instance (minus the service name).

Finally, use
```
docker inspect <instance_name> | jq '.[0].Config.Env'
```
to inspect the environment variables that have been passed into the container implementation.
The variables defined by the workload are shown.
Variables prefixed with "HZN" are provided by the Horizon Edge node platform.
A workload container implementation can exploit the "HZN" environment variables.
See the [Environment Variable Documentation](https://github.com/open-horizon/anax/blob/master/doc/managed_workloads.md) for a complete description of Horizon Edge node variables.

The workload container test environment can be stopped:
```
hzn dev workload stop
```

## Add a workload dependency

So far, the workload is executing with mocked CPU usage data, not very interesting.
Now let's update the workload project, introduing the CPU usage microservice as a dependency.
We are going to use the local microservice project we created earlier as the source of the dependency.
Notice that the -p flag in the following command points to the default `horizon` directory where the microservice project's metadata resides.
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

After a few seconds, the log output from the workload container can be found in syslog.
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

The last part of this guide shows you how to publish the CPU usage average to your Watson IoT Platform.

Modify your project as follows:
* Configure HTTPS authentication in the global section of the `userinput.json` file (so that the WIoTP containers can be fetched and loaded):
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
hzn dev dependency fetch -s https://internetofthings.ibmcloud.com/wiotp-edge/microservices/edge-core-iot-microservice --ver 2.3.1 -o IBM -a amd64 -k /etc/horizon/trust/publicWIoTPEdgeComponentsKey.pem
```
* The command above added a section to the `userinput.json` file under `microservices` for the edge-core-iot-microservice.
Fill in the variable values like this (substitute for the env vars):
```
    "variables": {
        "WIOTP_DEVICE_AUTH_TOKEN": "$WIOTP_GW_TOKEN",
        "WIOTP_DOMAIN": "$HZN_ORG_ID.messaging.internetofthings.ibmcloud.com"
    }
```
* Update `userinput.json` under the `workload` section to enable the workload to publish the CPU average (turn on PUBLISH and set the auth token, doing env var substitution):
```
    "variables": {
        "PUBLISH": true,
        "SAMPLE_INTERVAL": 2,
        "SAMPLE_SIZE": 5,
        "VERBOSE": "1"
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

I've noticed that sometimes the first publish fails, ignore this message and allow the workload to continue running:
```
Connection Refused: not authorised.
```

* Stop the workload:
```bash
hzn dev workload stop
```

At this point you have successfully completed a quick pass through the development process.
You have developed and tested a microservice and a workload as standalone containers and running within the Horizon Edge node test environment.
You have also integrated your projects with your Watson IoT Platform instance and are able to publish data to it.
In order to make these projects available for other Edge nodes to run them, you need to publish your projects.
The next section describes how to use `hzn dev` to publish you projects.

## Deploying the projects
When you are satisfied that the microservice and workload are working correctly, you can deploy them to the Edge so that any Edge node in your organization can run them.
The first step is to create a key pair that will be used to sign the deployment configuration of your microservice and workload.
```bash
hzn key create <x509-org> <x509-cn>
```
where `x509-org` is a company name or organization's name that is suitable to be used as an x509 certificate organization name, and `x509-cn` is an x509 certificate common name (preferably an email address issued by the `x509-org` organization).

This command will generate a public key and a private key.
The private key will be used to sign the microservice and workload, the public is needed by any Edge node that wants to run the microservice and workload.

The second step is to publish the microservice.
```bash
cd ~/hzn/ms/cpu
hzn dev microservice publish -k <private-key>
```
where `private-key` is the private key you generated in the previous step.

You can verify that the microservice was published:
```bash
hzn exchange microservice list
```

When the microservice is successfully published, upload your microservice container to the docker registry that you are using. The microservice project has this built into the makefile:
```bash
make publish
```

The third step is to publish the workload, and the process is remarkably similar to publishing a microservice.
```bash
cd ~/hzn/workload/cpu2wiotp
hzn dev workload publish -k <private-key>
```
where `private-key` is the private key you generated in the first step.

You can verify that the workload was published:
```bash
hzn exchange workload list
```

When the workload is successfully published, upload your workload container to the docker registry that you are using. The workload project has this built into the makefile:
```bash
make publish
```

The final step occurs on the node where you want to run the microservice and workload.
Since that's probably not your development machine, you should save the key pair someplace secure until you need them again.
The public signing key generated in the first step must be imported into the Horizon agent on the Edge node.
You must do this before invoking `hzn register` on the node.
```bash
hzn key import -k <public-key>
```
where `public-key` is the public key generated in the first step.

## Expand the project to multiple hardware architectures

The project metadata is hard wired to your WIoTP configuration (organization) and to the hardware architecture of your development machine.
This information is burned into the Horizon metadata files.
The final step in this guide is to enchance the project to make the metadata more reusable across different users and different hardware architectures.
Accomplishing this is simple, if you simply capture this variability in environment variables and parameterize the metadata files with those environment variables.
Using a tool like `envsubst` is ideal for creating concrete horizon metadata based on parameterized metadata files.
Then add a few new recipes to your `Makefile` and you can automate the entire process of creating the metadata files needed by horizon.
See the [Horizon examples project](https://github.com/open-horizon/examples) for an example of how to do this.

Both the cpu_percent and the cpu2wiotp sub projects are setup in this way.
Notice how the Makefile in each sub project creates the `horizon_build` directory and puts the conditioned metadata in there.
This enables the use of `hzn dev microservice start` and `hzn dev worklaod start` to work correctly for these projects.
This is the recommended approach when the project needs to be used by developers in more than 1 organization, or working with containers on more than 1 hardware architecture, etc.

Earlier in this guide you cloned the examples project.
You can see this approach in action by:
```
cd ~/examples/edge/wiotp/cpu2wiotp
make hznstart
```

The workload and all the dependent microservice containers have been started and are sending data to your WIoTP.

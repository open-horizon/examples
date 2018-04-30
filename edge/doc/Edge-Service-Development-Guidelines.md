# Guidelines for Developing Horizon Edge Services

## Introduction

Automating running services at the edge is substantially different from running services in the cloud:
* The number of edge nodes can be much greater.
* The networks to edge nodes can be unreliable and much slower. And edge nodes are often behind firewalls so connections usually can not be initiated from the cloud to the edge nodes.
* Edge nodes are usually not set up by ops staff, and they may be in remote locations. So you have to assume that once initially set up, there won't be a person around to do anything on the edge node.
* Edge nodes are usually less trusted environments than cloud servers.

These differences mean that different techniques must be used to deploy and manage software at the edge. Horizon is designed for this, but the services that are written by other developers must follow the following guidelines to fit into this structure well.

## Guidelines For Services

1. **Size matters** - service containers must be as small as possible so they deploy well over slow networks and to small edge devices.
    * Preferred programming languages that help achieve this are:
        * **Best**: go, rust, c, sh
        * **Ok**: c++, python, bash
        * **Not recommended**: nodejs, java and JVM-based languages (scala, etc.)
	* Some techniques to make docker images smaller:
        * Use alpine as the base linux image.
        * When installing packages in an alpine-based image, use: `apk --no-cache --update add …` (this avoids storing the pkg cache, which is not needed for runtime).
        * Recognize that deleting files in a subsequent layer (Dockerfile statement) doesn't actually remove the space from the image (because it is still in the previous layer), unless you squash the image at the end (docker save it to a tar file and docker load it). If you don't want to squash it, another common technique is to do many steps on the same `RUN` statement, separated by `&&`, so you can delete temporary files in the same layer as they were created.
        * Do not include any build tools in the runtime image. Best practice is for the build script/makefile to use a separate docker container to build the runtime artifacts (compile code, etc.), and then just copy the artifacts (executables, etc) into the runtime image.
1. **Self-Contained** - because the service has to be shipped over the network to a wide variety of edge nodes, the service containers should everything the service need bundled into the container (certificates, etc.). Don't rely on users doing something specific to the edge node in order for the service to run successfully.
1. **Well-Designed Configuration** - edge nodes need to be as close to zero-ops as possible. This includes the services that run on them. Horizon automates the deploy and management of the services, but the services must be structured properly to enable Horizon to do this without human intervention:
    * **UserInput variables** specified in the service definition are values that are different for every edge node and must be specified by the edge node owner, for example a device token. For each service, this list of variables should be as short as possible, preferably none, or at least have default values for most of them so they don't always have to be specified by the edge node owner.
    * **Environment configuration variables** are specified the deployment field in the service definition, and can be overridden in a pattern that uses that service. These are values that can be the same for a entire organization or at least a device type. Since they are specified in the service or pattern definitions in the exchange, they can be centrally managed, and Horizon can push out new values automatically to all edge nodes using the service. It also means the edge node owner does not have to do any manual work for these to be set.
    * **Locally overriding variables** - if it is necessary for your service to support this, it should be done in a way such that the service has appropriate defaults, so the node owner only has to provide values if the defaults are not satisfactory. The userInput variables with default values can be used for this. If a service has a lot of configurable settings, it may choose to support a config file, but a default version of the config file should be shipped in the service container so the edge node owner doesn't always have to provide the config file.


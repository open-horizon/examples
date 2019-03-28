# Horizon speedtest2wiotp Service

This Service used the speedtest REST service to get data about WAN connectivity, and then publishes this data to the IBM Watson IoT Platform ("WIoTP") MQTT service.

Interested parties with the appropriate credentials may then observe this data from anywhere on the planet by subscribing to this WIoTP MQTT service.

## Preconditions

The standard Linux `make` tool is used to operate on this code.  Please see the local `Makefile` for additional details.

In addition, you must get an IBM Cloud account (free or paid) and provision a WIoTP "Organization" and within that organization create a Device Type, and an instance of that type (a Device ID and Dvice Token). To subscribe, you will also need to provision an API key name and access token. See "Using the WIoTP Web Pages" below for instructions.

Finally, you must enter these details into the `wiotp.config` file in this directory.

Edge Nodes that register for the corresponding pattern will also need to fill in the `horizon/userinput.json` file with appropriate credentials as well. Each Edge Node will require its own Device ID to publish.

## Building

To build and tag the `speedtest2wiotp` Service docker container for the local architecture, within this directory run make with no target:
```
    $ make
```

## Testing

To test the `speedtest2wiotp` Service container, first build and run it, then run subscribe to the WIoTP MQTT service, e.g.:

```
    $ make
    $ make test
```

## Pushing To DockerHub

When you are ready, `docker login` to your DockerHub account. Once that succeeds then you can push an appropriately-tagged image to account `openhorizon` in DockerHub with this command:

```
    $ make push
```

## Publishing the Service to the Exchange

Once you have managed to push the image to DockerHub, then you can publish it to the Horizon Exchange as a "public" service in the "IBM" organization. Begin by setting up your IBM org credentials in your shell environment and then run this command:

```
    $ make service-publish
```

## Publishing a Pattern to the Exchange

Once you have managed to publish the Service to the Exchange, you can create a corresponding "public" pattern to the Horizon Exchange in the "IBM" organization that references this Service. Begin by setting up your IBM org credentials in your shell environment and then run this command:

```
    $ make pattern-publish
```

## Development Environment

To facilitate development of the `speedtest2wiotp` Service, you may wish to use the `dev` target:

```
    $ make dev
```

This will build the `speedtest2wiotp` Service container, then mount this working directory as `/outside` within the container and run `/bin/sh` in the container.  In that shell, `cd /outside` and then you can work on the original files in persistent storage outside the container, and also run them within the context of the container.

## Using the WIoTP Web Pages

More detail to come soon, here!  :-)



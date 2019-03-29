# Horizon speedtest2wiotp Service

This Service used the speedtest REST service to get data about WAN connectivity, and then publishes this data to the IBM Watson IoT Platform ("WIoTP") MQTT service.

Interested parties with the appropriate credentials may then observe this data from anywhere on the planet by subscribing to this WIoTP MQTT service.

## Preconditions

The standard Linux `make` tool is used to operate on this code, so it must be installed. Please see the local `Makefile` for additional usage details.

In addition, you require an IBM Cloud account (free or paid) to provision a WIoTP `Organization` and within that organization create an IoT `Device Type`, and an instance of that type (specifying a `Device ID` and `Device Token`). See "Using the WIoTP Web Pages" below for instructions.

When you register an Edge Node with a pattern containing `speedtest2wiotp` you will need to provide a `userinput.json` file defining these values, e.g.:

```
{
    "services": [
        {
            "org": "IBM",
            "url": "github.com.open-horizon.examples.speedtest2wiotp",
            "versionRange": "[0.0.0,INFINITY)",
            "variables": {
                "WIOTP_ORG": "theorg",
                "WIOTP_DEVICE_TYPE":  "MyDeviceType",
                "WIOTP_DEVICE_ID":    "MyDeviceId0",
                "WIOTP_DEVICE_TOKEN": "MyDeviceId0Token"

            }
        }
    ]
}
```

Please note that each Edge Node will require its own Device ID to publish.

## Local Testing Preconditions

If you wish to do (optional) local testing of this service's publiation to the IBMN Cloud Watson IoT Platform you will need to install an appropriate subscription tool, and setup some credentials.

The `mosquitto_sub` command (found in the popular `mosquitto-clients` package) and the `jq` utility are used for testing, so they must be installed if you wish to do this.

You must also provision a Watson IoT Platform API key name and access token to enable you to authenticate and view the data streams from your Edge Node devices. See "Using the WIoTP Web Pages" below for instructions.

To use the `make test` target for testing, you must configure a few environment variables with your credentials, e.g.:

```
export WIOTP_ORG="theorg"
export WIOTP_API_KEYNAME="a-theorg-thekeyname"
export WIOTP_API_TOKEN="token-for-that-key"
export WIOTP_DEVICE_TYPE="MyDeviceType"
export WIOTP_DEVICE_ID="MyDeviceId0"
```

## Building

To build and tag the `speedtest2wiotp` Service docker container for the local architecture, within this directory run `make build`:

```
    $ make build
```

## Testing

To run this container with its dependency, the `speedtest` Service, you need to tell Horiozn abouot the dependency:

```
   $ hzn dev dependency fetch --arch=$ARCH --org IBM --url github.com.open-horizon.examples.speedtest
```

(where $ARCH is the local hardware architecture, using the Horizon name for that, which is the Go language name. Typically this will be `amd64`, `arm`, or `arm64`).

For the next step, in addition to setting up ARCH, you need to setup a semver VERSION variable, e.g.:

```
    $ export VERSION='1.0.0'
```

Once that is setup, you can start `speedtest2wiotp` together with `speedtest`, by using this command:

```
   $ hzn dev service start -S -f horizon/userinput.json
```

Note that you must have previously setup the `userinput.json` file appropriately as described above.

To verify that the `speedtest2wiotp` Service container is publishing to Watson IoT Platform, first build and run it, then you can use `make test` to subscribe to the IBM Cloud WIoTP MQTT data feed from your device, e.g.:

```
    $ make test
```

Note that you must have previously appropriately configured environment vairables for your API key credentials, etc. as described above.


## Pushing To DockerHub

When you are ready to share an update, `docker login` to your DockerHub account. Once that succeeds then you can push an appropriately-tagged image to account `openhorizon` in DockerHub with this command:

```
    $ make push
```

## Publishing the Service to the Exchange

Once you have managed to push the image to DockerHub, then you can publish it to the Horizon Exchange as a "public" service in the "IBM" organization. Begin by setting up your IBM org credentials in your shell environment and then run this command:

```
    $ make service-publish
```

## Publishing a Pattern to the Exchange

Once you have managed to publish the changed Service (and any changed dependencies) to the Exchange, you can create a corresponding "public" pattern to the Horizon Exchange in the "IBM" organization that references this Service. Begin by setting up your IBM org credentials in your shell environment and then run this command:

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



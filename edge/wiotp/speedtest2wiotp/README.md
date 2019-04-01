# Horizon speedtest2wiotp Service

This Service used the speedtest REST service to get data about WAN connectivity, and then publishes this data to the IBM Watson IoT Platform ("WIoTP") MQTT service.

Interested parties with the appropriate credentials may then observe this data from anywhere on the planet by subscribing to this WIoTP MQTT service.

## Preconditions

The standard Linux `make` tool is used to operate on this code, so it must be installed. Please see the local `Makefile` for additional usage details.

In addition, you require an IBM Cloud account (free or paid) to provision a WIoTP `Organization` and within that organization create an IoT `Device Type`, and an instance of that type (specifying a `Device ID` and `Device Token`). See "Preparing WIoTP" below for detailed instructions.

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

You must also provision a Watson IoT Platform API key name and access token to enable you to authenticate and view the data streams from your Edge Node devices. See "Preparing WIoTP" below for detailed instructions.

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

## Preparing WIoTP

To begin, if you do not already have one, you need to create an IBM cloud account (free, or paid), and sign in to that account at this URL: https://cloud.ibm.com/login

Next, provision a Watson Internet of Things instance. Begin by selecting the "Create Resource" button at the top right to bring up the resource catalog. Then select "Internet of Things" in the left menu panel of the catalog. Then select the card for "Internet of Things Platform" when it comes up on the right. On the Internet of Things Platform page, give your service a name, select a region, organization, space, and then scroll down to select a pricing plan. Then select the "Create" button to provision your new service.

When your Internet of Things Platform instance comes up, select the "Launch" button to start the instance. When it comes up, save your 6-character organization name. You can find it as the prefix of the web page's domain name, or at the top right in the ID field. E.g., if the page URL is "https://abc123.internetofthings.ibmcloud.com/dashboard/devices/browse" then "abc123" is your IoT Platform "organization" name, and you will need to know this later.

Next you need to define a "Device Type" and then create a "Device ID" instance of that type for each machine you register. You can do this programmatically using the platform's REST APIs, but to do just one, it is convenient to use the web UI.

- In the menu panel on the left, hover over the device icon and select DEVICES
- Select the Device Types tab
- Select Add Device Type
- Select Device (not Gateway)
- Enter a name for your new IoT Device Type
- Select Next, Next, and then Done.
- Select the Browse tab
- Select Add Device
- Select the Device Type you just created (it may be pre-selected)
- Enter a Device ID for your edge node device
- Select Next and Next again
- Enter an authentication token *and record it somewhere* (once you leave this page you can never get back that token, so if you lost it you would need to delete the Device ID and create a new one).
Select Next and click Done

At thia point you have a Watson IoT Platform Device created, and you can setup your Edge Node as this Device, and it can send data to Watson IoT Platform.

Now you need to create an API key so you can subscribe to the data stream from your Internet of Things Platform Device:

- In the menu panel on the left, hover over the compass icon and select APPS
- Select Generate API Key
- Give it a name
- Optionally, enter a Description
- Select Next
- Select Standard Application
- Select Generate Key
- record the API Key and Authentication Token somewhere (again, it is important to save this since it cannot be recovered after you leave this page)

When you have finished all of the above, you should have the following information recorded:

- Your IBM cloud account name and password
- Your IoT Platform organization ID (6 characters).
- Your Device Type name
- Your Device instance ID, and its associated authentication token
- Your API Key name, and its associated authentication token

You will need all of these things to follow through the steps above.

## Configuring "Cards" to Visualize Your Data on the WIoTP Web Pages

The IBM Watson IoT Platform provides some free visualization tools for your data streams. This section will show you how to configure visualization "cards" for the Dowload Speed, Upload Speed, and Latency/Ping Delay data in the Speedtest data stream, similar to those shown below:

![cards](https://github.com/open-horizon/examples/tree/master/edge/wiotp/speedtest/cards.png "Cards in a WIoTP Board")

Begin by hovering over the top icon in the left menu and select BOARDS. Then select the "+ Create New Board" button. Give the board a name (e.g., "SpeedTest"), and optionally, a description. Select the options lou wish (e.g., to add this board to your favorites and navbar). The select the Next button.

Choose with whom you wish to share the board, then select the Submit button. This board will be created and will show up with your other boards on the boards page. Now select your new board.

On the page for your new board, select the "+ Add New Card" button. Select a "Devices" type card, e.g., "Line Chart". Then select which of your Device IDs you wish to visualize (select individually or use a serach). Then select the Next button.

On the "Create ... Card" page select "(+)" to connect a new data set, select the event you wish to visualize (the one called "status" is for this example code). select the property you wish to visualize, e.g., "upload". This will graph the "upload" field values from each "status" message received. Select the type "Number" and enter a large "Max" of "10000000000". Then select the "Next" button.

Now configure the type and size of chart you want to use, then select the Next button. Then finally make a color choice and select the Submit button.

You should see this card appear on your board, and start reflecting your actual data values as status message are received.

Repeat this porocess to add other cards to visualize other fields in the status event, if desired.


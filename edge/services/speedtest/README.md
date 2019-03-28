# Horizon Speedtest REST Service

The shared "speedtest" Service REST API provides WAN connectivity data.

## Preconditions

The standard Linux `make` tool is used to operate on this code.  Please see the local `Makefile` for additional details.

## Building

To build and tag the `speedtest` Service docker container for the local architecture, within this directory run make with no target:
```
    $ make
```

## Testing

To test the `speedtest` Service container, build and run it, then run a simple curl probe if its REST API, e.g.:

```
    $ make
    $ make test
```

If the first speed test is still running when you run `make test`, you should expect to see output similar to the following:

```
 $ make test
curl -sS localhost:5659/v1/speedtest | jq
{
  "error": "Patience, please. No speed test data received yet."
}
```

Once the `speedtest` Service has "warmed up" and test data is available, you will receive a more informative response, similar to that shown below:

```
 $ make test
curl -sS localhost:5659/v1/speedtest | jq
{
  "download": 1576985993.253844,
  "upload": 133087928.26745234,
  "ping": 10.986,
  "server": {
    "url": "http://speedtest.rd.ks.cox.net/speedtest/upload.php",
    "lat": "37.6922",
    "lon": "-97.3372",
    "name": "Wichita, KS",
    "country": "United States",
    "cc": "US",
    "sponsor": "Cox - Wichita",
    "id": "16623",
    "host": "speedtest.rd.ks.cox.net:8080",
    "d": 43.13860244182284,
    "latency": 10.986
  },
  "timestamp": "2019-03-26T15:36:26.791333Z",
  "bytes_sent": 151519232,
  "bytes_received": 409373932,
  "share": "http://www.speedtest.net/result/8141400046.png",
  "client": {
    "ip": "169.63.203.237",
    "lat": "37.751",
    "lon": "-97.822",
    "isp": "SoftLayer Technologies",
    "isprating": "3.7",
    "rating": "0",
    "ispdlavg": "0",
    "ispulavg": "0",
    "loggedin": "0",
    "country": "US"
  }
}
```

## Pushing To DockerHub

When you are ready, `docker login` to your DockerHub account. Once that succeeds then you can push an appropriately-tagged image to account `openhorizon` in DockerHub with this command:

```
    $ make push
```

## Publishing to the Exchange

Once you have managed to push the image to DockerHub, then you can publish it to the Horizon Exchange as a "public" service in the "IBM" organization. Begin by setting up your IBM org credentials in your shell environment and then run this command:

```
    $ make service-publish
```

## Development Environment

To facilitate development of the `speedtest` Service, you may wish to use the `dev` target:

```
    $ make dev
```

This will build the `speedtest` Service container, then mount this working directory as `/outside` within the container and run `/bin/sh` in the container.  In that shell, `cd /outside` and then you can work on the original files in persistent storage outside the container, and also run them within the context of the container.


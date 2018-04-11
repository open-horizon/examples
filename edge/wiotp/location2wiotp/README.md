# Horizon Location Workload

The Location workload polls the shared "gps" microservice REST APIs to
get location data (latitude, longitude, elevation) which has either been
statically provided by the owner, or derived from the device's IP address,
or may have been provided by GPS hardware when available.  It also polls the
gps microservice REST API to get satellite data,  The location and satellite
data is then sent to the data ingest service.

## Preconditions

The standard Linux `make` tool is used to operate on this code.  Please see the local `Makefile` for additional details.  The standard Blue Horizon image does not include `make`, so you will want to install make as follows:
```
    $ apt-get install make
```

## Building

To build and tag the location workload docker container for the local architecture, go to this directory and run make with no target:
```
    $ make
```

To build for a different architecture, for example to build arm on x86, set the ARCH environment variable first:
```
    $ export ARCH=arm
    $ make
```

## Testing

To test the location workload container, build it, run it as a daemon for testing purposes, then finally run the test program.  Execute these commands from this directory to accomplish that:
```
    $ make
    $ make daemon
    $ make test
```

## Publishing to the Docker Registry

To publish the location workload docker container (for this local architecture) to the Docker Hub registry, and create the defintion for it in the Horizon Exchange:
```
    $ make exchange-publish    # this will also build the docker image
```

## Development and Test Development

To facilitate development of the location workload, use the `develop` target:
```
    $ make develop
```
This will build the location workload container, then mount this working directory and run `/bin/sh` in that container.  In that shell, `cd /outside` and then you can work on the original files here outside the container, and run them in the context of the container.

In that same shell you can develop the test code. A make target is also provided to facilitate monitoring publications to the central MQTT service:
```
    $ make listener
```

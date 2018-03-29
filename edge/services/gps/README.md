# Horizon GPS REST Microservice

The shared "gps" microservice REST API provides location data
(i.e., latitude, longitude, elevation) which has either been
statically provided by the owner, or derived from the device's IP address,
or may have been provided by GPS hardware when available.  The source of
the location data is configured by means of well-known variables that are
expected to be set in the process environment of the `gps` container when
it is run by the Horizon infrastructure.

## Preconditions

The standard Linux `make` tool is used to operate on this code.  Please see the local `Makefile` for additional details.  The standard Horizon image does not include `make`, so you will want to install make as follows:
```
    $ apt-get install make
```

## Building

To build and tag the `gps` microservice docker container for the local architecture, go to this directory and run make with no target:
```
    $ make
```

## Testing

To test the `gps` microservice container, build it, run it as a daemon for testing purposes, then finally run the test program.  Execute these commands from this directory to accomplish that:
```
    $ make
    $ make daemon
    $ make test
```

## Publishing to the Horizon Docker Registry

To publish the `gps` microservice docker container (for this local architecture) to the Docker Hub registry, and create the defintion for it in the Horizon Exchange:
```
    $ make
    $ make exchange-publish
```

## Development and Test Development

To facilitate development of the `gps` microservice, use the `develop` target:
```
    $ make develop
```
This will build the `gps` microservice container, then mount this working directory and run `/bin/sh` in that container.  In that shell, `cd /outside` and then you can work on the original files here outside the container, and run them in the context of the container.

Similarly, to facilitate development of the test container for the `gps` microservice, use the `develop-test` target:
```
    $ make develop-test
```
This will build the test container, then mount this working directory and run `/bin/bash` in that container.  In that shell, `cd /outside` and then you can work on the original test code files here outside the container, and run them in the context of the container.

## Debugging

To debug the gps microservice, you can connect directly to the servcie from any shell on the host as follows:
```
    export gps_ip=`docker inspect gps | grep IPAddress | tail -1 | sed 's/.*: "//;s/",//'`
    curl -s http://$gps_ip:31779/v1/gps | jq
```

If the gpsd daemon (not locally developed code) is suspect, you can debug it by using docker exec to make a shell in the gps micorservice container, and then connecting to the gpsd socket:
```
    docker exec -it gps /bin/sh
    telnet localhost 2947
    ?WATCH={"enable":true,"json":true}
    e
```
`e` is the exit command.  The WATCH command will stream the data until you stop it.  You can Google for more info on the *gpsd wire protocol*.

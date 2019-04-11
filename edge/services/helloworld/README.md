# Horizon Hello World Example Edge Service

## Using the Hello World Example Edge Service

- First, go through the "Try It" page "Installing Horizon Software On Your Edge Machine" to set up your edge node.
- Get an IBM cloud account (and for now have your org created in the exchange)
- Set your exchange org:
```
export HZN_ORG_ID="<yourorg>"
```
- Set your exchange user credentials in the Horizon-supported environment variable and verify it:
```
export HZN_EXCHANGE_USER_AUTH="iamapikey:<myapikey>"
hzn exchange user list
```
- Choose a id and token for your edge node, create it, and verify it:
```
export HZN_EXCHANGE_NODE_AUTH="<mynodeid>:<mynodetoken>"
hzn exchange node create -n $HZN_EXCHANGE_NODE_AUTH
hzn exchange node list ${HZN_EXCHANGE_NODE_AUTH%%:*}
```
- Register your edge node with Horizon to use the helloworld pattern:
```
hzn register -n "$HZN_EXCHANGE_NODE_AUTH" $HZN_ORG_ID IBM/pattern-helloworld
```
- Look at the Horizon agreement until it is finalized and then see the running container:
```
hzn agreement list
docker ps
```
- See the helloworld service output:
```
# soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep helloworld
# For now on Mac:
docker logs -f $(docker ps -q --filter name=helloworld)
``` 
- Unregister your edge node, stopping the helloworld service:
```
hzn unregister -f
```

## First-Time Edge Service Developer - Building and Publishing Your Own Version of the Hello World Example Edge Service

If you want to create your own Horizon edge service, follow the next 2 sections to copy the hello world example and start modifying it.

### Preconditions for Developing Your Own Service

- First, go through the steps in the section above to run the IBM helloworld service on an edge node.
- Get a docker hub id at https://hub.docker.com/ , if you don't already have one. (This example is set up to store the docker image in docker hub, but by modifying DOCKER_IMAGE_BASE you can store it in another registry.) Login to the docker registry using your id:
```
echo 'mydockerpw' | docker login -u mydockehubid --password-stdin
```
- If you are reading this from github.com (not locally), clone this repo and cd to this directory:
```
git clone git@github.com:open-horizon/examples.git
cd edge/services/helloworld
# copy it where you want to work on it, so you can commit it to your own git repo
# Soon you will be able to instead use: hzn dev service new -s <service-name> -v <version> -i <image>
```
- If you have the HZN_ORG_ID environment variable set from previous work, unset it (this value will now come from `horizon/hzn.cfg`):
```
unset HZN_ORG_ID
```
- Set the variable values in `horizon/hzn.cfg` to your own values.
- Soon these steps will not be needed, but for now do them:
    - Enable `hzn` to read `horizon/hzn.cfg`: `alias hzn='source horizon/hzn.cfg && hzn'`
    - Set the exchange URL: `export HZN_EXCHANGE_URL=https://alpha.edge-fabric.com/v1`
    - Set the architecture:
```
export ARCH=$(uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')
```
- As part of the above section "Using the Hello World Example Edge Service", you created your Exchange user credentials and edge node credentials. Ensure they are set and verify them:
```
export HZN_EXCHANGE_USER_AUTH="iamapikey:<myapikey>"
hzn exchange user list
export HZN_EXCHANGE_NODE_AUTH="<mynodeid>:<mynodetoken>"
hzn exchange node list ${HZN_EXCHANGE_NODE_AUTH%%:*}
```

### Building and Publishing Your Own Version of the Hello World Example Edge Service

- Edit `service.sh`, for example changing "Hello" to "Hey there"
    - Note: this service is a shell script simply for brevity, but you can write your service in any language.
- Build the hello world docker image:
```
make
```
- Test the service by having Horizon start it locally:
```
hzn dev service start -S
```
- See the docker container running and look at the output:
```
docker ps
# soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep helloworld
# For now on Mac:
docker logs -f $(docker ps -q --filter name=helloworld)
```
- See the environment variables Horizon passes into your service container:
```
docker inspect $(docker ps -q --filter name=helloworld) | jq '.[0].Config.Env'
```
- Stop the service:
```
hzn dev service stop
```
- Create a service signing key pair (if you haven't already done so):
```
mkdir -p horizon/keys   # soon this will not be needed
hzn key create -d horizon/keys IBM my@email.com
# soon these 2 commands will not be needed:
source horizon/hzn.cfg && mv horizon/keys/*-private.key $HZN_PRIVATE_KEY_FILE
source horizon/hzn.cfg && mv horizon/keys/*-public.pem $HZN_PUBLIC_KEY_FILE
```
- Have Horizon push your docker image to your registry and publish your service in the Horizon Exchange and see it there:
```
# soon the -k and -K flags will not be necessary
hzn exchange service publish -f horizon/service.definition.json -k $HZN_PRIVATE_KEY_FILE -K $HZN_PUBLIC_KEY_FILE
hzn exchange service list
```
- Publish your edge node deployment pattern in the Horizon Exchange and see it there:
```
# soon the -p flag will not be needed
hzn exchange pattern publish -f horizon/pattern.json -p pattern-helloworld-$ARCH
hzn exchange pattern list
```
- Register your edge node with Horizon to use your deployment pattern:
```
hzn register -n "$HZN_EXCHANGE_NODE_AUTH" $HZN_ORG_ID pattern-helloworld-$ARCH
```
- Look at the Horizon agreement until it is finalized and then see the running container:
```
hzn agreement list
docker ps
```
- See the helloworld service output:
```
# soon you will use 'hzn service log ...' for all platforms
# For now on Linux:
tail -f /var/log/syslog | grep helloworld
# For now on Mac:
docker logs -f $(docker ps -q --filter name=helloworld)
``` 
- Unregister your edge node, stopping the helloworld service:
```
hzn unregister -f
```

## Further Learning

To see more Horizon features demonstrated, continue on to the [cpu2msghub example](../../msghub/cpu2msghub).

## Process for the Horizon Development Team to Make Updates to the Helloworld Service

- Do the steps in the Preconditions section above, except:
    - export `HZN_EXCHANGE_URL` to the staging instance
    - export `HZN_EXCHANGE_USER_AUTH` to your credentials in the IBM org
- Make whatever code changes are necessary
- Increment `SERVICE_VERSION` in `horizon/hzn/cfg`
- Make the files that `HZN_PRIVATE_KEY_FILE` and `HZN_PUBLIC_KEY_FILE` point to actually be symbolic links to the common keys we use to sign all of our examples.
- Build, test, and publish for all architectures:
```
make publish-all-arches
```
Note: building all architectures works on mac os x, and can be made to work on ubuntu via: http://wiki.micromint.com/index.php/Debian_ARM_Cross-compile , https://wiki.debian.org/QemuUserEmulation

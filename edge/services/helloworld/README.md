# Horizon Hello World Example Edge Service

## Preconditions for Developing Your Own Service

- First, go through the quick start guide for setting up your edge node and running the helloworld service.
- Get a docker hub id at https://hub.docker.com/ , if you don't already have one. (This example is set up to store the docker image in docker hub, but by modifying DOCKER_IMAGE_BASE you can store it in another registry.) Login to the docker registry using your id:
```
echo 'mydockerpw' | docker login -u mydockehubid --password-stdin
```
- If you are reading this from github.com (not locally), clone this repo and cd to this directory:
```
git clone git@github.com:open-horizon/examples.git
cd edge/services/helloworld
# copy it where you want to work on it, so you can commit it to your own git repo
# Soon you will be able to instead use: hzn dev service new ...
```
- Set the variable values in `horizon/hzn.cfg` to your own values.
- Soon these steps will not be needed, but for now do them:
  - Enable `hzn` to read `horizon/hzn.cfg`: `alias hzn='source horizon/hzn.cfg && hzn'`
  - Set the architecture:
```
export ARCH=$(uname -m | sed -e 's/aarch64.*/arm64/' -e 's/x86_64.*/amd64/' -e 's/armv.*/arm/')
```
  - Set the exchange URL: `export HZN_EXCHANGE_URL=https://alpha.edge-fabric.com/v1`
- As part of the quick start mentioned in the first step above, you created your Exchange user credentials and edge node credentials. Set those here in the Horizon-supported environment variables and verify them:
```
export HZN_EXCHANGE_USER_AUTH="iamapikey:<myapikey>"
hzn exchange user list
export HZN_EXCHANGE_NODE_AUTH="<mynodeid>:<mynodetoken>"
hzn exchange node list ${HZN_EXCHANGE_NODE_AUTH%%:*}
```

## Building and Publishing Your Own Version of the Hello World Example Edge Service

- Build the hello world docker image:
```
make
```
- Have Horizon start the service locally:
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
- Push your docker image to your registry and publish your service in the Horizon Exchange and see it there:
```
# soon the -k and -K flags will not be necessary
hzn exchange service publish -f horizon/service.definition.json -k $HZN_PRIVATE_KEY_FILE -K $HZN_PUBLIC_KEY_FILE
hzn exchange service list
```
- Publish your edge node deployment pattern in the Horizon Exchange and see it there:
```
# soon the -p flag will not be needed
hzn exchange pattern publish -f horizon/pattern/pattern-helloworld.json -p pattern-helloworld-$ARCH
hzn exchange pattern list
```
- Register your edge node with Horizon to use your deployment pattern (choose a node id and token):
```
hzn register $HZN_ORG_ID pattern-helloworld-$ARCH
```
- Look at the Horizon agreement until 1 it is finalized and then see the running container:
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

# Further Learning

To see more Horizon features demonstrated, continue on to the cpu2msghub example.

# Creating Your Own Watson Speech to Text to IBM Event Streams Edge Service for Raspberry Pi

Follow the steps in this page to create your Watson Speech to Text to IBM Event Streams Edge Service.

## Preconditions for Developing Your Own Service

1. If you have not already done so, complete the steps in these sections:

  - [Preconditions for Using the Watson Speech to Text to IBM Event Streams Example Edge Service](README.md#preconditions)
  - [Using the Watson Speech to Text to IBM Event Streams Edge Service with Deployment Pattern](README.md#using-watsons2text-pattern)

2. If you are using macOS as your development host, configure Docker to store credentials in `~/.docker`:

  - Open the Docker **Preferences** dialog
  - Uncheck **Securely store Docker logins in macOS keychain**

3. If you do not already have a docker ID, obtain one at https://hub.docker.com/ . Log in to Docker Hub using your Docker Hub ID:

  ```bash
  export DOCKER_HUB_ID="<dockerhubid>"
  echo "<dockerhubpassword>" | docker login -u $DOCKER_HUB_ID --password-stdin
  ```

  Output example:

  ```bash
  WARNING! Your password will be stored unencrypted in /home/pi/.docker/config.json.
  Configure a credential helper to remove this warning. See
  https://docs.docker.com/engine/reference/commandline/login/#credentials-store

  Login Succeeded
  ```

4. Create a cryptographic signing key pair. This enables you to sign services when publishing them to the exchange. **This step only needs to be done once.**

  ```bash
  hzn key create "<x509-org>" "<x509-cn>"
  ```

  where `<x509-org>` is your company name, and `<x509-cn>` is typically set to your email address.

5. Install `git` and `jq`:

  ```bash
  sudo apt install -y git jq
  ```


## <a id=build-publish-your-wst> Building and Publishing Your Own Version of the Watson Speech to Text to IBM Event Streams Service

1. Clone this git repo:

```bash
cd ~   # or wherever you want
git clone git@github.com:open-horizon/examples.git
cd examples/
```

2. Check your Horizon CLI version:

```bash
hzn version
```

3. Starting with Horizon version `v2.29.0-595` you can `checkout` to a version of the example services that directly corresponds to your Horizon CLI version with these commands: 

```bash
export EXAMPLES_REPO_TAG="v$(hzn version 2>/dev/null | grep 'Horizon CLI' | awk '{print $4}')"
git checkout tags/$EXAMPLES_REPO_TAG -b $EXAMPLES_REPO_TAG
```
**Note:** if you are using an older version of the `hzn` CLI you can checkout to the branch that corresponds to the major version you are using. For example Horizon CLI version: `2.27.0-173` can run `git checkout v2.27`

4. Copy the `watson_speech2text` dir to where you will start development of your new service:

```
cp -a edge/evtstreams/watson_speech2text ~/myservice     # or wherever
cd ~/myservice
```

5. Set the values in `horizon/hzn.json` to your own values.

6. Edit `watsonspeech2text.py` however you want.
    - Note: this service is a shell script simply for brevity, but you can write your service in any language.

7. Build the watsons2text docker image:

```bash
make
```

8. Test the service by having Horizon start it locally:
```bash
hzn dev service start -S
```

9. Check that the container is running:
```bash
sudo docker ps 
```

9. See the watsons2text service output:

```bash
hzn dev service log -f $SERVICE_NAME
```

10. See the environment variables Horizon passes into your service container:
```bash
docker inspect $(docker ps -q --filter name=${SERVICE_NAME}) | jq '.[0].Config.Env'
```

11. Stop the service:
```bash
hzn dev service stop
```

12. Have Horizon push your docker image to your registry and publish your service in the Horizon Exchange and see it there:
```bash
hzn exchange service publish -f horizon/service.definition.json
hzn exchange service list
```

13. Publish your edge node deployment pattern in the Horizon Exchange and see it there:
```bash
hzn exchange pattern publish -f horizon/pattern.json
hzn exchange pattern list
```

14. Register your edge node with Horizon to use your deployment pattern:
```bash
hzn register -p pattern-${SERVICE_NAME}-${ARCH} -f horizon/userinput.json
```

15. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:
```bash
hzn agreement list
```

16. Once the agreement is made, list the docker container edge service that has been started as a result:
```bash
sudo docker ps
```


17. On any machine, install [kafkacat](https://github.com/edenhill/kafkacat#install), then subscribe to the Event Streams topic to see the json data that watsons2text is sending:

```bash
kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$EVTSTREAMS_API_KEY -X ssl.ca.location=$EVTSTREAMS_CERT_FILE -t $EVTSTREAMS_TOPIC
```


18. See the watsons2text service output:
```bash
hzn dev service log -f $SERVICE_NAME
``` 

19. Unregister your edge node, stopping the watsons2text service:
```bash
hzn unregister -f
```


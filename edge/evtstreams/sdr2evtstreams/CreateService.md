# Creating Your Own SDR To IBM Event Streams Edge Service

Follow the steps in this page to create your SDR To IBM Event Streams Edge Service.

## Preconditions for Developing Your Own Service

1. If you have not already done so, complete the steps in these sections:

  - [Preconditions for Using the SDR To IBM Event Streams Example Edge Service](README.md#preconditions)
  - [Using the SDR To IBM Event Streams Edge Service with Deployment Pattern](README.md#using-sdr2evtstreams-pattern)

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

  On **Linux**:

  ```bash
  sudo apt install -y git jq
  ```

  On **macOS**:

  ```bash
  brew install git jq
  ```


## <a id=building-your-own-sdr2evtstreams-pattern></a> Building and Publishing Your Own Version of the SDR To IBM Event Streams Edge Service

1. Clone this git repo:
```bash
cd ~   # or wherever you want
git clone git@github.com:open-horizon/examples.git
```

2. Checkout the branch that corresponds to your horizon CLI version. To get the branch name, remove the last bullet and any numbers after it, then prepend a `v` at the beginning:
```bash
$ hzn version
Horizon CLI version: 2.27.0-173 # Branch name in this example is v2.27
Horizon Agent version: 2.27.0-173
$ git checkout v2.27
```

3. Copy the `sdr2evtstreams` dir to where you will start development of your new service:
```bash
cp -a examples/edge/evtstreams/sdr2evtstreams ~/myservice     # or wherever
cd ~/myservice
```

4. Set the values in `horizon/hzn.json` to your own values.

5. Edit `main.go` however you want.
    - Note: this service is written in go, but you can write your service in any language.

6. Build the sdr2evtstreams docker image:
```bash
make
```

7. Test the service by having Horizon start it locally:
```bash
hzn dev service start -S
```

8. Check that the containers are running:
```bash
sudo docker ps
```

9. See the sdr2evtstreams service output:
```bash
hzn service log -f ibm.sdr2evtstreams
```

10. See the environment variables Horizon passes into your service container:
```bash
docker inspect $(docker ps -q --filter name=sdr2evtstreams) | jq '.[0].Config.Env'
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

14. Register your edge node with Horizon to use your deployment pattern (substitute `<service-name>` for the `SERVICE_NAME` you specified in `horizon/hzn.json`):
```bash
hzn register -p pattern-<service-name>-$(hzn architecture) -f horizon/userinput.json
```

15. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:
```bash
hzn agreement list
```

16. Once the agreement is made, list the docker container edge service that has been started as a result:
```bash
sudo docker ps
```

17. On any machine, subscribe to the Event Streams topic to see the json data that sdr2evtstreams is sending:
```bash
kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$EVTSTREAMS_API_KEY -t $EVTSTREAMS_TOPIC
```

18. See the sdr2evtstreams service output:
```bash
hzn service log -f ibm.sdr2evtstreams
```

19. Unregister your edge node, stopping the sdr2evtstreams service:
```bash
hzn unregister -f
```

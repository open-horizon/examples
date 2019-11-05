# Creating Your Own SDR To IBM Event Streams Edge Service

Follow the steps in this page to create your first simple Horizon edge service.

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
```
cd ~   # or wherever you want
git clone git@github.com:open-horizon/examples.git
```
2. Copy the `sdr2evtstreams` dir to where you will start development of your new service:
```
cp -a examples/edge/evtstreams/sdr2evtstreams ~/myservice     # or wherever
cd ~/myservice
```
3. Set the values in `horizon/hzn.json` to your own values.

### Building and Publishing Your Own Version of the SDR To IBM Event Streams Edge Service

1. Edit `main.go` however you want.
    - Note: this service is written in go, but you can write your service in any language.
2. Build the sdr2evtstreams docker image:
```
make
```
3. Test the service by having Horizon start it locally:
```
hzn dev service start -S
```
4. Check that the containers are running:
```
sudo docker ps
```

5. See the sdr2evtstreams service output:

	on **Linux**:
	```
	tail -f /var/log/syslog | grep sdr2evtstreams[[]
	```

	on **Mac**:
	```
	docker logs -f $(docker ps -q --filter name=sdr2evtstreams)
	```

6. See the environment variables Horizon passes into your service container:
```
docker inspect $(docker ps -q --filter name=sdr2evtstreams) | jq '.[0].Config.Env'
```
7. Stop the service:
```
hzn dev service stop
```
8. Create a service signing key pair in `~/.hzn/keys/` (if you haven't already done so):
```
hzn key create <my-company> <my-email>
```
9. Have Horizon push your docker image to your registry and publish your service in the Horizon Exchange and see it there:
```
hzn exchange service publish -f horizon/service.definition.json
hzn exchange service list
```
10. Publish your edge node deployment pattern in the Horizon Exchange and see it there:
```
hzn exchange pattern publish -f horizon/pattern.json
hzn exchange pattern list
```
11. Register your edge node with Horizon to use your deployment pattern (substitute for `SERVICE_NAME` the value you specified above for `hzn dev service new -s`):
```
hzn register -p pattern-SERVICE_NAME-$(hzn architecture) -f horizon/userinput.json
```
12. Look at the Horizon agreement until it is finalized and then see the running container:
```
hzn agreement list
docker ps
```
13. On any machine, subscribe to the Event Streams topic to see the json data that sdr2evtstreams is sending:
```
kafkacat -C -q -o end -f "%t/%p/%o/%k: %s\n" -b $EVTSTREAMS_BROKER_URL -X api.version.request=true -X security.protocol=sasl_ssl -X sasl.mechanisms=PLAIN -X sasl.username=token -X sasl.password=$EVTSTREAMS_API_KEY -X ssl.ca.location=$EVTSTREAMS_CERT_FILE -t $EVTSTREAMS_TOPIC
```

14. See the sdr2evtstreams service output:

	on **Linux**:
	```
	tail -f /var/log/syslog | grep sdr2evtstreams[[]
	```

	on **Mac**:
	```
	docker logs -f $(docker ps -q --filter name=sdr2evtstreams)
	``` 

15. Unregister your edge node, stopping the sdr2evtstreams service:
```
hzn unregister -f
```

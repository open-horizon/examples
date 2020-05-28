# Creating Your Own Operator Edge Service

Follow the steps on this page to create your first operator that deploys an edge service to a cluster, and learn how you can pass horizon environment variables (in addition to any other environment variables needed) to your deployed service pods. 

## Preconditions for Developing Your Own Service

1. If you have not already done so, complete the steps in these sections:

  - [Preconditions for Using the Operator Example Edge Service](README.md#preconditions)
  - [Using the Operator Example Edge Service with Deployment Pattern](README.md#using-helloworld-pattern)

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

## <a id=build-publish-your-op> Building and Publishing Your Own Operator Example Edge Service

In order to deploy a containerized edge service to an edge cluster, a software developer first has to build a Kubernetes Operator that deploys the containerized edge service in a Kubernetes cluster. There are several options when writing a Kubernetes operator. To start, the Kubernetes open source documentation has an [operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) overview article which is a good resource to learn about operators. Visit the [operator-sdk getting started](https://github.com/operator-framework/getting-started#overview) page to find the steps followed to create the operator used in this example. 

When following the steps in the operator-sdk documentation, the [controller file](https://github.com/operator-framework/getting-started#add-a-new-controller) is where you can specify the containerized edge service you want to deploy to the edge cluster. If you look at [line 196](https://github.com/t-fine/examples/blob/2851c6abb17ccf2fbd760f5e5f494f3c4a668328/edge/services/operator/simple-operator/pkg/controller/ibmserviceoperator/ibmserviceoperator_controller.go#L196) of the controller file you can see the `open-horizon/ibm.helloworld_amd64:1.0.0` docker image is specified. 

If you have gone through the `ibm.helloworld` example before this then you know it uses the `HZN_DEVICE_ID` in the service log output. If you want to propagate the horizon environment variables (or any other environment variables necessary for your edge service to function properly) to the service pods it will require a few minor additions to your controller file. The code needed to make environment variables to your service is mostly confined to the reconcile loop and can be seen on [line 178](https://github.com/t-fine/examples/blob/37d88da91c13052e26461d4c1f6bb164ae9abec2/edge/services/operator/simple-operator/pkg/controller/ibmserviceoperator/ibmserviceoperator_controller.go#L178), and [line 198](https://github.com/t-fine/examples/blob/37d88da91c13052e26461d4c1f6bb164ae9abec2/edge/services/operator/simple-operator/pkg/controller/ibmserviceoperator/ibmserviceoperator_controller.go#L198) with the addition of importing `os` on [line 6](https://github.com/t-fine/examples/blob/37d88da91c13052e26461d4c1f6bb164ae9abec2/edge/services/operator/simple-operator/pkg/controller/ibmserviceoperator/ibmserviceoperator_controller.go#L6). 

1. Clone this git repo:
```bash
cd ~   # or wherever you want
git clone git@github.com:open-horizon/examples.git
```

2. Copy the `operator` dir to where you will start development of your new service:
```bash
cp -a examples/edge/services/operator ~/myoperator     # or wherever
cd ~/myoperator/simple-operator/deploy
```

3. Set the values in `horizon/hzn.json` to your own values and update the path in `horizon/service.definition.json` to point to the `ibm.operator.tar.gz` file
  - Node: the `ibm.operator.tar.gz` file contains everything in the `deploy/` directory.

Testing operators is different from testing Horizon services for edge devices. If you have created your own operator file for this example, you can test it by following the steps on the `operator-sdk` getting started page under the [Run as a Deployment inside the cluster](https://github.com/operator-framework/getting-started#1-run-as-a-deployment-inside-the-cluster) section.

5. With the operator tested, instruct Horizon to push your docker image to your registry and publish your service in the Horizon Exchange:

  ```bash
  hzn exchange service publish -f horizon/service.definition.json
  hzn exchange service list
  ```

6. **Optional:** Modify the `service.policy.json` file located in the `horizon/` directory to contain a constraint of your choosing and add it to your published operator service:
  ```bash
  hzn exchange service addpolicy -f horizon/service.definition.json <your-operator-service>
  ```
  
7. Modify the `deployment.policy.json` file located in the `horizon/` directory by changing the `"constraints": ["openhorizon.service == ibm.operator"]` to something else uniquely identifying the edge cluster you want to run your operator service.

8. Publish this Deployment Policy to the Exchange giving it a memorable name of your choosing:
   ```bash
   hzn exchange business addpolicy -f deployment.policy.json <choose-any-policy-name>
   ```
   
9. Modify the `node.policy.json` file located in the `horizon/` directory by changing the `"properties":`  value to that of the constraint you spedified in the `deployment.policy.json` so they match and will form an agreement.

10. Get your modified `node.policy.json` file onto your edge cluster then tegister your cluster with your new node policy:
  ```bash
  hzn register --policy horizon/node.policy.json
  ```
  
11. Veryfy that the `simple-operator` deployment is up and runing:
```bash 
kubectl get pod -n openhorizon-agent
```

- If everything deployed correctly you will see the operator pod in addition to three `example-ibmserviceoperator` pods running similar to following output
```bash 
NAME                                          READY   STATUS    RESTARTS   AGE
agent-6d8b8895f-bpwm9                         1/1     Running   0          2d21h
example-ibmserviceoperator-7d6849c487-5pmcb   1/1     Running   0          88s
example-ibmserviceoperator-7d6849c487-926xt   1/1     Running   0          88s
example-ibmserviceoperator-7d6849c487-j9z58   1/1     Running   0          88s
simple-operator-5cd47878fc-gjcl6              1/1     Running   0          96s
```

12. Verify that the operator is running successfully by checking its logs:
```bash
kubectl logs simple-operator-<op-id> -n openhorizon-agent
```

- If the operator is opperating correctly the logs should look similar to the following output:
```bash
root@gormand1:~# kubectl logs simple-operator-5cd47878fc-gjcl6 -n openhorizon-agent
{"level":"info","ts":1590090974.906244,"logger":"cmd","msg":"Operator Version: 0.0.1"}
{"level":"info","ts":1590090974.9062827,"logger":"cmd","msg":"Go Version: go1.14.3"}
{"level":"info","ts":1590090974.9063098,"logger":"cmd","msg":"Go OS/Arch: linux/amd64"}
{"level":"info","ts":1590090974.9063208,"logger":"cmd","msg":"Version of operator-sdk: v0.17.1"}
{"level":"info","ts":1590090974.9066129,"logger":"leader","msg":"Trying to become the leader."}
{"level":"info","ts":1590090975.1842163,"logger":"leader","msg":"No pre-existing lock was found."}
{"level":"info","ts":1590090975.1895027,"logger":"leader","msg":"Became the leader."}
{"level":"info","ts":1590090975.443696,"logger":"controller-runtime.metrics","msg":"metrics server is starting to listen","addr":"0.0.0.0:8383"}
{"level":"info","ts":1590090975.4556718,"logger":"cmd","msg":"Registering Components."}
{"level":"info","ts":1590090976.0551462,"logger":"metrics","msg":"Metrics Service object created","Service.Name":"simple-operator-metrics","Service.Namespace":"openhorizon-agent"}
{"level":"info","ts":1590090976.3069465,"logger":"cmd","msg":"Could not create ServiceMonitor object","error":"no ServiceMonitor registered with the API"}
{"level":"info","ts":1590090976.3070312,"logger":"cmd","msg":"Install prometheus-operator in your cluster to create ServiceMonitor objects","error":"no ServiceMonitor registered with the API"}
{"level":"info","ts":1590090976.3070836,"logger":"cmd","msg":"Starting the Cmd."}
{"level":"info","ts":1590090976.308044,"logger":"controller-runtime.manager","msg":"starting metrics server","path":"/metrics"}
{"level":"info","ts":1590090976.3089783,"logger":"controller-runtime.controller","msg":"Starting EventSource","controller":"ibmserviceoperator-controller","source":"kind source: /, Kind="}
{"level":"info","ts":1590090976.40968,"logger":"controller-runtime.controller","msg":"Starting EventSource","controller":"ibmserviceoperator-controller","source":"kind source: /, Kind="}
{"level":"info","ts":1590090976.5111873,"logger":"controller-runtime.controller","msg":"Starting Controller","controller":"ibmserviceoperator-controller"}
{"level":"info","ts":1590090976.5112767,"logger":"controller-runtime.controller","msg":"Starting workers","controller":"ibmserviceoperator-controller","worker count":1}
{"level":"info","ts":1590090976.511785,"logger":"controller_ibmserviceoperator","msg":"Reconciling IBMserviceOperator","Request.Namespace":"openhorizon-agent","Request.Name":"example-ibmserviceoperator"}
{"level":"info","ts":1590090976.612423,"logger":"controller_ibmserviceoperator","msg":"Creating a new Deployment","Request.Namespace":"openhorizon-agent","Request.Name":"example-ibmserviceoperator","Deployment.Namespace":"openhorizon-agent","Deployment.Name":"example-ibmserviceoperator"}
```

13. View the `ibm.helloworld` output to verify that the operator sucessfully deployed the service and the environment variables were passed into the pod:
```bash
kubectl logs example-ibmserviceoperator-<ex-op-id> -n openhorizon-agent
```

- if the environment variables were received by the worker pods the ouput should look similar to the following:
```bash
kubectl logs example-ibmserviceoperator-<ex-op-id> -n openhorizon-agent
tfine-cluster-apollo1 says: Hello from the cluster!!!
tfine-cluster-apollo1 says: Hello from the cluster!!!
tfine-cluster-apollo1 says: Hello from the cluster!!!
```

14. Unregister your edge node (which will also stop the myhelloworld service):
```bash
hzn unregister -f
```

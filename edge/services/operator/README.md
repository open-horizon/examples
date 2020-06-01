# Horizon Operator Example Edge Service

This is a simple example of using and creating an operator as edge service.

- [Preconditions for Using the Operator Example Edge Service](#preconditions)
- [Using the Operator Example Edge Service with Deployment Policy](#using-operator-policy)
- [Using the Operator Example Edge Service with Deployment Pattern](PatternRegister.md)
- [Creating Your Own Operator Edge Service](CreateService.md)
- Further Learning - to see more Horizon features demonstrated, continue on to the [cpu2evtstreams example](../../evtstreams/cpu2evtstreams).

## <a id=preconditions></a> Preconditions for Using the Operator Example Edge Service

If you haven't done so already, you must do these steps before proceeding with the operator example:

1. Install the Horizon management infrastructure (exchange and agbot).

2. Install the Horizon agent on your edge cluster and configure it to point to your Horizon exchange.

3. As part of the infrasctucture installation process for IBM Edge Computing Manager a file called `agent-install.cfg` was created that contains the values for `HZN_ORG_ID` and the exchange and css url values. Locate this file and set those environment variables in your shell now:

```bash
eval export $(cat agent-install.cfg)
```

 - **Note**: if for some reason you disconnected from ssh or your command line closes, run the above command again to set the required environment variables.

4. In addition to the file above, an API key associated with your Horizon instance would have been created, set the exchange user credentials, and verify them:

```bash
export HZN_EXCHANGE_USER_AUTH=iamapikey:<horizon-API-key>
hzn exchange user list
```

5. Choose an ID and token for your edge node, create it, and verify it:

```bash
export HZN_EXCHANGE_NODE_AUTH="<choose-any-node-id>:<choose-any-node-token>"
hzn exchange node create -n $HZN_EXCHANGE_NODE_AUTH
hzn exchange node confirm
```

6. If you have not done so already, unregister your node before moving on:
 ```bash
hzn unregister -f
```

## <a id=using-operator-policy></a> Using the Operator Example Edge Service with Deployment Policy

In the following steps you will deploy the `ibm.operator` to your edge cluster. This operator will then create three pods running the `ibm.helloworld` service. 

1. To execute the `hzn` commands below, `exec` inside the agent container
```bash
kubectl exec -it agent-<agent-pod-id> -n openhorizon-agent /bin/bash
```

2. Get the required node policy file on your edge cluster:
```bash
wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/operator/simple-operator/deploy/horizon/node.policy.json
```

- Below is the `node_policy.json` file you obtained in step one, above:

```json
{
  "properties": [
    { "name": "openhorizon.service", "value": "ibm.operator" }
  ],
  "constraints": [
  ]
}
```

- It provides values for one `properties` (`openhorizon.service`), that will effect which services get deployed to this edge node, and states no `constraints`.

3. Register your Node Policy with this policy
```bash
hzn register --policy node.policy.json
```

4. When the registration completes, use the following command to review the Node Policy:
```bash
hzn policy list
```

5. The edge device will make an agreement with one of the Horizon agreement bots (this typically takes about 15 seconds). Repeatedly query the agreements of this device until the `agreement_finalized_time` and `agreement_execution_start_time` fields are filled in:
```bash
hzn agreement list
```

- Below is the `service.policy.json` that has been published into the Exchange when it was created:
```json
{
    "properties": [],
    "constraints": [
        "openhorizon.arch == amd64"
    ]
}
```

- Below is the `business.policy.json` that has been published into the Exchange when it was created:
```json
{
  "mycluster/ibm.operator_1.0.0": {
    "owner": "mycluster/operator1",
    "label": "ibm.operator Business Policy for amd64",
    "description": "A super-simple sample Horizon Business Policy",
    "service": {
      "name": "ibm.operator",
      "org": "mycluster",
      "arch": "*",
      "serviceVersions": [
        {
          "version": "1.0.0",
          "priority": {},
          "upgradePolicy": {}
        }
      ],
      "nodeHealth": {}
    },
    "constraints": [
      "openhorizon.service == ibm.operator"
    ],
    "userInput": [
      {
        "serviceOrgid": "mycluster",
        "serviceUrl": "ibm.operator",
        "serviceVersionRange": "[0.0.0,INFINITY)",
        "inputs": [
          {
            "name": "HW_WHO",
            "value": "from the cluster!"
          }
        ]
      }
    ],
    "created": "2020-05-21T17:18:34.956Z[UTC]",
    "lastUpdated": "2020-05-21T19:50:56.937Z[UTC]"
  }
}
```

6. Veryfy that the `simple-operator` deployment is up and runing:
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

7. Verify that the operator is running successfully by checking its logs:
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

8. Verify that the operator sucessfully deployed the `ibm.helloworld` service and the environment variables were passed into the pod:
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

9. Unregister your edge node (which will also stop the myhelloworld service):
```bash
hzn unregister -f
```

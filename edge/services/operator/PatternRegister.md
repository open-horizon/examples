# Horizon Operator Example Edge Service

## Preconditions for Developing Your Own Service

1. If you have not already done so, complete the steps in this section:

  - [Preconditions for Using the Operator Example Edge Service](README.md#preconditions)
  
## <a id=using-operator-pattern></a> Using the Operator Example Edge Service with Deployment Pattern

1. Register your edge node with Horizon to use the `ibm.operator` pattern:

  ```bash
  hzn register -p IBM/pattern-ibm.operator-amd64 -s ibm.operator --serviceorg IBM -u $HZN_EXCHANGE_USER_AUTH
  ```
 - **Note**: using the `-s` flag with the `hzn register` command will cause Horizon to wait until agreements are formed and the service is running on your edge node to exit, or alert you of any errors encountered during the registration process. 

  2. Veryfy that the `simple-operator` deployment is up and runing:
  ```bash 
  kubectl get pod -n openhorizon-agent
  ```

- If everything deployed correctly you will see the operator pod in addition to three `example-ibmserviceoperator` pods running similar to following output
  ```bash 
  NAME                                          READY   STATUS    RESTARTS   AGE
  agent-6d8b8895f-bpwm9                         1/1     Running   0          2d21h
  example-ibmserviceoperator-7d6849c487-5pmcb   1/1     Running   0          88s
  simple-operator-5cd47878fc-gjcl6              1/1     Running   0          96s
  ```

**Note:** If you are attempting to run this service on an **OCP edge cluster** and the operator does not start you may have to grant the operator the privileges it requires to execute with the following command:
 ```bash
 oc adm policy add-scc-to-user privileged -z simple-operator -n openhorizon-agent
 ```

3. Verify that the operator is running successfully by checking its logs:
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

4. Verify that the operator sucessfully deployed the `ibm.helloworld` service and the environment variables were passed into the pod:
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

5. Unregister your edge node (which will also stop the myhelloworld service):
  ```bash
  hzn unregister -f
  ```

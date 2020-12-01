# Creating Your Own Operator Edge Service

Follow the steps on this page to create your first ansible operator that deploys an edge service to a cluster. With this operator we will expose the nginx service so you can `curl` it externally using a route on OCP, or thru the clusters IP address if you are using a microk8s or k3s edge cluster.

## Preconditions for Developing Your Own Service

1. If you have not already done so, complete the steps in these sections:

   - [Preconditions for Using the Operator Example Edge Service](README.md#preconditions)
   - [Using the Operator Example Edge Service with Deployment Policy](README.md#using-operator-policy)

2. If you are using macOS as your development host, configure Docker to store credentials in `~/.docker`:

   - Open the Docker **Preferences** dialog
   - Uncheck **Securely store Docker logins in macOS keychain**

3. Install [operator-sdk](https://github.com/operator-framework/operator-sdk/tree/v0.17.x#prerequisites) and all its prerequisites. This example was created using version `0.17`, and at has not been tested on any other versions.

4. Install the Kubenetes CLI [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

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

   * On **Linux**:

   ```bash
   sudo apt install -y git jq
   ```

  * On **macOS**:

   ```bash
   brew install git jq
   ```

## <a id=build-publish-your-op> Building and Publishing Your Own Operator Example Edge Service

In order to deploy a containerized edge service to an edge cluster, a software developer first has to build a Kubernetes Operator that deploys the containerized edge service in a Kubernetes cluster. There are several options when writing a Kubernetes operator. This example will guide through creating an ansible operator. The following steps were originally performed on the cluster host machine. 

1. Create a new operator application and generate a default directory layout based on the input name, and move into the created default oeprator directory:
   ```bash
   operator-sdk new my-operator --type=ansible --api-version=my.operator.com/v1alpha1 --kind=MyOperator
   cd my-operator/
   ```

**Note:** If you are following these steps to create an operator that will deploy your own service, you should modify the names used above, **not** modify things manually in the generated files. 

The above command will give you an empty ansible operator. At the very least you will need to define a `deployment`, `service`, and a set of tasks to deploy your service. If you look in the [ansible-role-files](https://github.com/open-horizon/examples/tree/master/edge/services/hello-operator/ansible-role-files) directory you can see these three files and how they are used to deploy the `nginxinc/nginx-unprivileged` image, and how they expose the port `8080`. There are several methods for configuring cluster network traffic, which can can read about in the [OCP traffic overview documentation](https://docs.openshift.com/container-platform/4.6/networking/configuring_ingress_cluster_traffic/overview-traffic.html). For this service we've going to create a service of type `NodePort` and and `expose` our service with a route.

2. Obtain the `deployment`, `service`, and task file and move them into the `my-operator` roles directory with the following commands:
   ```bash
   wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/nginx-operator/ansible-role-files/deployment.j2 && mv deployment.j2 roles/myoperator/templates/ 
   wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/nginx-operator/ansible-role-files/service.j2 && mv service.j2 roles/myoperator/templates/
   wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/nginx-operator/ansible-role-files/route.j2 && mv route.j2 roles/myoperator/templates/
   wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/nginx-operator/ansible-role-files/main.yml && mv main.yml roles/myoperator/tasks/
   ```

3. Build the operator image:
   ```bash
   operator-sdk build docker.io/$DOCKER_HUB_ID/my.operator_amd64:1.0.0
   docker push docker.io/$DOCKER_HUB_ID/my.operator_amd64:1.0.0
   ```

4. In the `deploy/operator.yaml` file replace `"REPLACE_IMAGE"` with the operator image name you build and pushed in the previous step:
   ```bash
   image: "<docker-hub-id>/my.operator_amd64:1.0.0"
   ```

5. By default the operator does not have the permission to create routes, however, with the following command you can add the lines needed to the `deploy/role.yaml` file so the operator can expose the `nginx` service with a route:
   ```bash
   echo "- apiGroups:
     - route.openshift.io
     resources:
     - routes
     - routes/custom-host
     verbs:
     - get
     - list
     - watch
     - patch
     - update
     - create
     - delete" >> deploy/role.yaml
   ```

6. Apply the required resources to run the operator 
   ```bash
   kubectl apply -f deploy/crds/my.operator.com_myoperators_crd.yaml
   kubectl apply -f deploy/service_account.yaml
   kubectl apply -f deploy/role.yaml
   kubectl apply -f deploy/role_binding.yaml
   kubectl apply -f deploy/operator.yaml
   kubectl apply -f deploy/crds/my.operator.com_v1alpha1_myoperator_cr.yaml
   ```

7. Ensure the operator pod and the deployed service pod are running:
   ```bash
   kubectl get pods
   ```

If everything deployed correctly you should see an output similar to the following after around 60 seconds:
   ```
   NAME                           READY   STATUS    RESTARTS   AGE
   nginx-7d5598fb56-vw6lz         1/1     Running   0          12s
   my-operator-55c6f56c47-b6p7c   1/1     Running   0          48s
   ```

8. Check that the service is up:
   ```bash
   kubectl get service
   ```

If everything deployed correctly you should see an output similar to the following after around 60 seconds:
   ```
   NAME                  TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)             AGE
   nginx                 NodePort    172.30.37.113    <none>        80:30080/TCP        24h
   ```

If you are using an **OCP edge cluster** you will need to `curl` the service using the exposed `route`.
9. Get the exposed route name:
   ```bash
   kubectl get route -n openhorizon-agent
   ```
  
If the route was exposed correctly you should see an output similar to the following:
   ```bash
   NAME          HOST/PORT                                                    PATH   SERVICES   PORT   TERMINATION   WILDCARD
   nginx-route   nginx-route-openhorizon-agent.apps.apollo5.cp.fyre.ibm.com          nginx      8080                 None
   ```

10. `curl` the service to test if it is functioning correctly:
   **OCP edge cluster** substitute the above `HOST/PORT` value:
      ```bash
      curl nginx-route-openhorizon-agent.apps.apollo5.cp.fyre.ibm.com
      ```
   
   **k3s or microk8s edge cluster**:
      ```bash
      curl <external-ip-address>:30080
      ```

If the service is running you should see following `Welcome to nginx!` output:
   ```
   <!DOCTYPE html>
   <html>
   <head>
   <title>Welcome to nginx!</title>
   <style>
       body {
           width: 35em;
         margin: 0 auto;
         font-family: Tahoma, Verdana, Arial, sans-serif;
      }
   </style>
   </head>
   <body>
   <h1>Welcome to nginx!</h1>
   <p>If you see this page, the nginx web server is successfully installed and
   working. Further configuration is required.</p>

   <p>For online documentation and support please refer to
   <a href="http://nginx.org/">nginx.org</a>.<br/>
   Commercial support is available at
   <a href="http://nginx.com/">nginx.com</a>.</p>

   <p><em>Thank you for using nginx.</em></p>
   </body>
   </html>
   ```

11. Delete the resources to stop your operator pod and service 
   ```bash
   kubectl delete crd myoperators.my.operator.com
   kubectl delete deployment my-operator
   kubectl delete service my-operator-metrics
   kubectl delete serviceaccount my-operator
   kubectl delete rolebinding my-operator
   kubectl delete role my-operator
   kubectl delete route nginx-route
   ```

**Note:** if any pods are stuck in the `Terminating` state after running the previous commands you can force delete them with the following command:
   ```bash
   kubectl -n <namespace> delete pods --grace-period=0 --force <pod_name(s)>
   ```

12. Create a tar archive that contains the files inside the operators `deploy/` directory:
   ```
   tar -zcvf operator.tar.gz deploy/*
   ```

## <a id=publish-op-service> Publish Your Operator Example Edge Service

**Note:** The following commands were performed from a separate machine, **not** the cluster host. 

1. Create a new working directory for a new horizon project:
   ```bash
   mkdir my-operator && cd my-operator/
   hzn dev service new -V 1.0.0 -s my-first-operator -c cluster
   ```

2. Transfer the `operator.tar.gz` archive you created in the previous section to your `my-operator/` directory.

3. Edit the `horizon/service.definition.json` file to point to the operator's yaml archive created in the previous step. Assuming it's in the `my-operator/` directory you can make it the following:
   ```
   "operatorYamlArchive": "../operator.tar.gz"
   ```

4. Publish your operator service:
   ```bash
   hzn exchange service publish -f horizon/service.definition.json
   ```

## <a id=publish-op-policy> Create a Deployment Policy For Your Operator Example Edge Service

1. Create a `deployment.policy.json` file to deploy your operator service to an edge cluster:
```bash
cat << 'EOF' > horizon/deployment.policy.json
{
  "label": "$SERVICE_NAME Deployment Policy",
  "description": "A super-simple sample Horizon Deployment Policy",
  "service": {
    "name": "$SERVICE_NAME",
    "org": "$HZN_ORG_ID",
    "arch": "*",
    "serviceVersions": [
      {
        "version": "$SERVICE_VERSION",
        "priority":{}
      }
    ]
  },
  "properties": [
  ],
  "constraints": [
    "example == my-operator"
  ],
  "userInput": [
  ]
}
EOF
```

Notice we have given this deployment policy the following constraint: `"example == my-operator"`

2. Publish your deployment policy:
   ```bash
   hzn exchange deployment addpolicy -f horizon/deployment.policy.json policy-my-operator
   ```

3. Back on your cluster host, create a `node.policy.json` file:
```bash
cat << 'EOF' > node.policy.json
{
  "properties": [
    { "name": "example", "value": "my-operator" }
  ]
}
EOF
```

4. Register your edge cluster with your new node policy:
   ```bash
   hzn register -u $HZN_EXCHANGE_USER_AUTH
   cat node.policy.json | hzn policy update -f-
   hzn policy list
   ```

Notice the node policy contains the property that matches the constraint specified by the deployment policy we created, which means your operator service will begin deploying to your edge cluster momentarily.

5. Check to see the agreement has been created (this can take approximately 15 seconds):
   ```bash
   hzn agreement list
   ```

6. Once you see an `agreement_execution_start_time`, you should start to see the operator pod and service deployment pod begin to start up:
   ```bash
   kubectl get pods -n openhorizon-agent
   ```

If everything deployed correctly you should see an output similar to the following:
   ```
   NAME                                   READY   STATUS    RESTARTS   AGE
   agent-dd984ff96-jmmdl                  1/1     Running   0          1d
   nginx-7d5598fb56-vw6lz                 1/1     Running   0          12s
   my-operator-55c6f56c47-b6p7c           1/1     Running   0          48s
   ```
7. Check that the service is up:
   ```bash
   kubectl get service -n openhorizon-agent
   ```

If everything deployed correctly you should see an output similar to the following after around 60 seconds:
   ```
   NAME                  TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)             AGE
   nginx                 NodePort    172.30.37.113    <none>        80:30080/TCP        24h
   ```

If you are using an **OCP edge cluster** you will need to `curl` the service using the exposed `route`.
8. Get the exposed route name:
   ```bash
   kubectl get route -n openhorizon-agent
   ```
  
If the route was exposed correctly you should see an output similar to the following:
   ```bash
   NAME          HOST/PORT                                                    PATH   SERVICES   PORT   TERMINATION   WILDCARD
   nginx-route   nginx-route-openhorizon-agent.apps.apollo5.cp.fyre.ibm.com          nginx      8080                 None
   ```

9. `curl` the service to test if it is functioning correctly:
   **OCP edge cluster** substitute the above `HOST/PORT` value:
      ```bash
      curl nginx-route-openhorizon-agent.apps.apollo5.cp.fyre.ibm.com
      ```
   
   **k3s or microk8s edge cluster**:
      ```bash
      curl <external-ip-address>:30080
      ```

If the service is running you should see following `Welcome to nginx!` output:
   ```
   <!DOCTYPE html>
   <html>
   <head>
   <title>Welcome to nginx!</title>
   <style>
       body {
           width: 35em;
         margin: 0 auto;
         font-family: Tahoma, Verdana, Arial, sans-serif;
      }
   </style>
   </head>
   <body>
   <h1>Welcome to nginx!</h1>
   <p>If you see this page, the nginx web server is successfully installed and
   working. Further configuration is required.</p>

   <p>For online documentation and support please refer to
   <a href="http://nginx.org/">nginx.org</a>.<br/>
   Commercial support is available at
   <a href="http://nginx.com/">nginx.com</a>.</p>

   <p><em>Thank you for using nginx.</em></p>
   </body>
   </html>
   ```

10. Unregister your edge cluster:
   ```
   hzn unregister -f
   ```

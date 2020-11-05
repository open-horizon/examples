# Creating Your Own Operator Edge Service

Follow the steps on this page to create your first ansible operator that deploys an edge service to a cluster that exposes a "Hello World" service on `<NODE_IP>:30007`.

## Preconditions for Developing Your Own Service

1. If you have not already done so, complete the steps in these sections:

   - [Preconditions for Using the Operator Example Edge Service](README.md#preconditions)
   - [Using the Operator Example Edge Service with Deployment Policy](README.md#using-operator-policy)

2. If you are using macOS as your development host, configure Docker to store credentials in `~/.docker`:

   - Open the Docker **Preferences** dialog
   - Uncheck **Securely store Docker logins in macOS keychain**

3. Install [operator-sdk](https://github.com/operator-framework/operator-sdk/tree/v0.17.x#prerequisites) and all its prerequisites. This example was created using version `0.17`, and at has not been tested on any other versions.

4. Install the Kubenetes CLI [kubectl]{https://kubernetes.io/docs/tasks/tools/install-kubectl/}

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
   operator-sdk new hello-operator --type=ansible --api-version=hello.operator.com/v1alpha1 --kind=HelloOperatorService
   cd hello-operator/
   ```

2. The above command will give you a blank ansible operator. At the very least you will need to define a `deployment`, `service`, and a set of tasks to deploy a service. For this example I have created these files to deploy the `openhorizon/mosquito.helloworld_amd64:1.0.0` image. Obtain them and move them into the `hellooperatorservice` roles directory with the following commands:
   ```bash
   wget deployment.j2 && mv deployment.j2 roles/hellooperatorservice/templates/ 
   wget service.j2 && mv service.j2 roles/hellooperatorservice/templates/
   wget main.yml && mv main.yml froles/hellooperatorservice/tasks/
   ```

3. Build the operator image:
   ```bash
   operator-sdk build docker.io/$DOCKER_HUB_ID/myhello.operator_amd64:1.0.0
   docker push docker.io/$DOCKER_HUB_ID/myhello.operator_amd64:1.0.0
   ```

4. In the `deploy/operator.yaml` file replace `"REPLACE_IMAGE"` with the operator image name you build and pushed in the previous step:
   ```
   image: "<docker-hub-id>/myhello.operator_amd64:1.0.0"
   ```

5. Apply the required resources to run the operator 
   ```bash
   kubectl apply -f deploy/crds/hello.operator.com_hellooperatorservices_crd.yaml
   kubectl apply -f deploy/service_account.yaml
   kubectl apply -f deploy/role.yaml
   kubectl apply -f deploy/role_binding.yaml
   kubectl apply -f deploy/operator.yaml
   kubectl apply -f deploy/crds/hello.operator.com_v1alpha1_hellooperatorservice_cr.yaml
   ```

6. Ensure the operator pod and the deployed service pod are running:
   ```bash
   kubectl get pods
   ```

If everything deployed correctly you should see an output similar to the following after around 60 seconds:
   ```
   NAME                                   READY   STATUS    RESTARTS   AGE
   hello-operator-6cf57bc5f-8cmjw         1/1     Running   0          89s
   mosquito-helloworld-7f7cb95db5-tq6m6   1/1     Running   0          76s
   ```

7. To test the service is functioning correctly you can curl the service using one of two methods:
   ```bash
   curl -sS <INTERNAL_IP>:8000 | jq .

   - or externally - 

   curl -sS <NODE_IP>:30007 | jq .
   ```

If the service is running you should see output similar to the following:
   ```json
   {
     "Hello": "10.22.29.174"
   }
   ```

8. Delete the resources to stop your operator pod and service 
   ```bash
   kubectl delete crd hellooperatorservices.hello.operator.com;
   kubectl delete deployment hello-operator;
   kubectl delete service hello-operator-metrics;
   kubectl delete serviceaccount hello-operator;
   kubectl delete rolebinding hello-operator;
   kubectl delete role hello-operator
   ```

**Note:** if any pods are stuck in the `Terminating` state after running the previous commands you can force delete them with the following command:
   ```bash
   kubectl -n <namespace> delete pods --grace-period=0 --force <pod_name(s)>
   ```

9. Create a tar archive that contains the files inside the operators `deploy/` directory:
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
    "org": "mycluster",
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
    "example == my-hello-operator"
  ],
  "userInput": [
  ]
}
EOF
```

Notice we have given this deployment policy the following constraint: `"example == my-hello-operator"`

2. Publish your deployment policy:
   ```bash
   hzn exchange deployment addpolicy -f horizon/deployment.policy.json policy-my-hello-operator
   ```

3. Back on your cluster host, create a `node.policy.json` file:
```bash
cat << 'EOF' > node.policy.json
{
  "properties": [
    { "name": "example", "value": "my-hello-operator" }
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
   hello-operator-6c5f8c4458-6ggwx        1/1     Running   0          24s
   mosquito-helloworld-7bccc7668c-x9qf7   1/1     Running   0          7s
   ```

7. You can now test the service is funcitoning correctly again by curl-ing the service using of of two methods:
   ```
   curl -sS <INTERNAL_IP>:8000 | jq .

   - OR externally with the following - 

   curl -sS <NODE_IP>:30001 | jq .
   ```

If the service is running you should see output similar to the following:
   ```json
   {
     "Hello": "10.22.29.174"
   }
   ```

8. Unregister your edge cluster:
   ```
   hzn unregister -f
   ```

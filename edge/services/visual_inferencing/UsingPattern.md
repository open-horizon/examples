# Horizon Object Detection and Classification Example Edge Service

## Preconditions for Developing Your Own Service

1. If you have not already done so, complete the steps in this section:

  - [Preconditions for Using the Operator Example Edge Service](UsingPolicy.md#preconditions)
  
## <a id=using-operator-pattern></a> Using the Object Detection and Classification Example Edge Service with Deployment Pattern

1. Register your edge node with Horizon to use the yolo pattern:
- if your edge device **does not** have a GPU, run the following command:
  ```bash
  hzn register -p IBM/pattern-ibm.yolocpu -s ibm.yolocpu --serviceorg IBM
  ```
- if your edge device **does** have a GPU, run the following command:
  ```bash
  hzn register -p IBM/pattern-ibm.yolocuda -s ibm.yolocuda --serviceorg IBM
  ```
 - **Note**: using the `-s` flag with the `hzn register` command will cause Horizon to wait until agreements are formed and the service is running on your edge node to exit, or alert you of any errors encountered during the registration process. 

2. Veryfy that the `yolo` service deployment is up and runing:
  ```bash 
  sudo docker ps
  ```
3. You can now navigate to http://0.0.0.0:5200 to confirm the object detection and classification is working as intended.

4. Unregister your edge node (which will also stop the myhelloworld service):
```bash
hzn unregister -f
```

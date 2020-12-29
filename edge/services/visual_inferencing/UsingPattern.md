# Horizon Object Detection and Classification Example Edge Service

## Preconditions for Developing Your Own Service

1. If you have not already done so, complete the steps in this section:

  - [Preconditions for Using the Operator Example Edge Service](UsingPolicy.md#preconditions)
  
## <a id=using-operator-pattern></a> Using the Object Detection and Classification Example Edge Service with Deployment Pattern

1. Get the user input file for the yolo object detection and classification sample:
- if your edge device **does not** have a GPU, run the following command:
  ```bash
  wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/visual_inferencing/yolocpu/horizon/userinput.json
  ```
- if your edge device **does** have a GPU, run the following command:
  ```bash
  wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/visual_inferencing/yolocuda/horizon/userinput.json
  ```
Below is the `userinput.json` file you got in the previous step:
  ```json
  {
    "services": [
        {
            "org": "IBM",
            "url": "yolocpu",
            "variables": {
                "EVTSTREAMS_API_KEY": "$EVTSTREAMS_API_KEY",
                "EVTSTREAMS_BROKER_URL": "$EVTSTREAMS_BROKER_URL",
                "EVTSTREAMS_CERT_ENCODED": "$EVTSTREAMS_CERT_ENCODED",
                "EVTSTREAMS_TOPIC": "$EVTSTREAMS_TOPIC",
                "CAM_URL": "https://github.com/MegaMosquito/achatina/raw/master/shared/restcam/mock.jpg"
            }
        }
    ]
  }
  ```
If you register with the pattern using the `userinput.json` file as is you can visit `http://0.0.0.0:5200` and see the classification being done on a [sample image](https://github.com/open-horizon/examples/tree/achatina/edge/services/visual_inferencing#object-detection-and-classification). Or, you can replace the `CAM_URL` value with your own feed and do the inferencing with that as well.  

2. Register your edge node with Horizon to use the yolo pattern:
- if your edge device **does not** have a GPU, run the following command:
  ```bash
  hzn register -p IBM/pattern-ibm.yolocpu -s ibm.yolocpu --serviceorg IBM -f userinput.json
  ```
- if your edge device **does** have a GPU, run the following command:
  ```bash
  hzn register -p IBM/pattern-ibm.yolocuda -s ibm.yolocuda --serviceorg IBM -f userinput.json
  ```
 - **Note**: using the `-s` flag with the `hzn register` command will cause Horizon to wait until agreements are formed and the service is running on your edge node to exit, or alert you of any errors encountered during the registration process. 

3. Veryfy that the `yolo` service deployment is up and runing:
  ```bash 
  sudo docker ps
  ```
4. You can now navigate to http://0.0.0.0:5200 to confirm the object detection and classification is working as intended. (This can take a couple minutes)

5. Unregister your edge node (which will also stop the myhelloworld service):
```bash
hzn unregister -f
```

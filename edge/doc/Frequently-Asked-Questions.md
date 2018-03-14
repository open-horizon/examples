# Frequently Asked Questions

### **How is "Edge with Watson IoT Platform" related to "Horizon" and "BlueHorizon"?**

**Blue Horizon** refers to the "citizen scientist" instance of Horizon that is found at [https://bluehorizon.network](https://bluehorizon.network).

**Edge with Watson IoT Platform** refers to IBM's edge solution integrated with the Watson IoT Platform, which is a different instance of Horizon.

Both of these Horizon instances are created using the [open-horizon](https://github.com/open-horizon) software.

**Horizon** simply refers to any instance of the open-horizon software.

### **Is the "Edge with Watson IoT" software open sourced?**

No, but Edge with WioTP is based on the [open-horizon](https://github.com/open-horizon) project which has been open sourced.  Many example programs found in the open source project will work in Edge + WIoTP.

### **How can I develop and deploy my edge software using IBM's Edge with WIoTP?**

At a high level the steps needed to develop software using Edge with WIoTP are:
 - write code (almost any language, any libraries, any modern linux)
 - containerize it in one or more docker container(s)
 - self-sign your containers with your cryptographic signing private key
 - publish your containers along with any needed configuration variables to the Horizon Exchange

To deploy your software you will:
 - define an Edge Node "type" in Watson IoT Platform
 - create unique instance IDs for each Edge Node (do this with the web GUI or programmatically using the REST API)

Visit each of your Edge Nodes and:
 - install the Linux packages containing the Edge with WIoTP Linux software
 - install your cryptographic signing public key
 - enter its type and unique ID and provide any configuration variables that are unique to the node
 - execute a command to register it with the Horizon Exchange

Soon after each Edge Node is registered, a Horizon AgBot should discover it and take on the responsibility for managing it's software from then on.  You should never need to again visit that Edge Node.

### **What Edge Node hardware platforms does Edge with WIoTP support?**

Currently debian package binaries for Edge with Watson IoT are provided for armhf (e.g., Raspberry Pi, ARM 32bit), ARM-64bit, and x86-64bit architectures).

### **Can I run any Linux distro on my Edge Nodes usign Edge with WIoTP?**

Yes, and no.  You may develop your software to run on any modern Linux distro within your docker container.  That is, you may base your docker container on any of the distros available for your architecture using the dockerfile FROM statement.  However, the Horizon software is currently only provided as a debian package.  The Edge Node host must therefore natively run a suitable debian variant (e.g., debian, ubuntu, raspbian, etc.).  That is, your code can run in any Linux distro context on top of the host's kernel inside your container, while the host will be beside your code running a debian variant on top of that kernel.

### **Which programming languages and evironments does Edge with WIoTP support?**

Almost any language and any software libraries, and any Linux variant can be built into your docker containers for deployment on any Horizon Edge Node.  If you are able to build and run your code in a docker container, then most likely it can be run on Horizon.

Note that if your software requires special hardware or operating systems servcies access to run, then you may need to specify "docker run" arguments to support that access in the "deployment" section of the definition file for your container.

### **Can I try Edge with Watson IoT Platform for free?**

Yes, but there are predefined limits on free accounts.  For example, the number of Edge Nodes is restricted and there is a monthly cap on the amount of data you can send.

### **Is there detailed documentation for the REST APIs provided by the compoinents in Edge with WIoTP?**

Yes!

The documentation for the Horizon Exchange's REST API is here:

 * [https://exchange.bluehorizon.network/api/api](https://exchange.bluehorizon.network/api/api)

The documentation for the Horizon Agent REST API is here:

 * [https://github.com/open-horizon/anax/blob/master/doc/api.md](https://github.com/open-horizon/anax/blob/master/doc/api.md)

The documentation for the Watson IoT Platform REST API is here:

 * [https://console.bluemix.net/docs/services/IoT/devices/api.html#api](https://console.bluemix.net/docs/services/IoT/devices/api.html#api)

### **Does Edge with WIoTP use Kubernetes?**

No. Edge with WIoTP has been implemented without using Kubernetes.

### **Does Edge with WIoTP use MQTT?**

Yes, and no.  Edge with WIoTP is set up to make it easy for your code to use MQTT, but the use of MQTT in your code is not required.

### **Is it possible to bypass TLS and communicate insecurely between my code and the Watson IoT Platform?**

Yes, but in general this is not recommended since it could be possible for third parties to eavesdrop on your communications.  To make TLS communications only "optional" (as opposed to the default where TLS is required) visit the Watson IoT Platform "dashboard" page for your organization (replacing "YOURORG" with your 6-character organization ID):
```
https://YOURORG.internetofthings.ibmcloud.com/dashboard/
```
Tap the (gear icon) "Settings" menu in the panel on the left, then under "Security" select "Connection Security", and use the pop-up menu there to select "TLS Optional".

### **Can my Edge Node code communicate directly with the WIoTP cloud MQTT instead of connecting through the "edge-connector" locally on the Edge Node?**

Yes, but doing that has some disadvantages:
 1) You will need to manage the storing and forwarding of events, whenever the Edge Node may lose connectivity to the Internet (as often occurs for IoT devices) -- this is handled automatically when you use the "edge-connector"
 2) You will need to store on the Edge Node, your app's API key and secret API token in order to communicate securely with the Watson IoT Platform in the cloud.  This unnecessairily exposes you to risk of those credentials being compromised on the Edge Node.
 3) You will need to create one app API key and secret API token pair for each Edge node that runs this workload.  You will not have this complication if you communicate through the "edge-connector"

### **How long does it normally take after I register my Edge Node before agreements are formed, and the corresponding containers start running?**

Typcially registration itself takes only a few seconds.  Once registration has completed, and "hzn node list" shows the state is "configured" then the Horizon AgBots in the Watson IoT Platform will be able to discover your Edge Node and begin to propose agreements.

Those agreements will normally be received and accepted in less than a minute, but in unusual circumstances they may take up to about 10 minutes to finalize and appear in the "hzn agreement list" output.

Once an agreement is finalized, the corresponding containers should begin running very shortly afterward.  Of course the Horizon Agent must first "docker pull" each container image, and it must also verify its cryptographic signature with the Horizon Exchange.  So those parts of the process are somewhat dependent upon the *size* of the containers.  Smaller containers, like those based on a micro Linux distribution like Alpine, often download in just a few seconds.

Once the containers for the agreement are downloaded and verified, an appropriate docker network will be created for them and they will be run.  That step also typically takes only seconds.  Once the containers are running, you will be able to see them listed in the "docker ps" output.

### **Something went wrong!  How can I completely remove the software and everything else related to Edge with WIoTP from my host, so I can restart the Quick Start Developer Guide from the beginning?**

If you have already registered your Edge Node, it is best to begin by unregistering it and removing it from the Horizon Registry:
```
hzn unregister -f -r
```

After that you can go ahead and cleanup all of the Horizon software:
```
sudo apt purge -y horizon horizon-cli horizon-ui horizon-wiotp
```

### **Is there in Edge with WIoTP a kind of dashboard to visualize the agreements, workloads, microservices that are on a particular Edge Node?**

Not yet, but this is coming soon.  The Blue Horizon instance of Horizon already has a web UI showing the status of its Horizon Edge Nodes and that code will soon be ported over to the Edge with WIoTP instance of Horizon.  Stay tuned...

### **During development I used short names for my microservices, like "cpu", but when the container is run by Horizon, the container name includes the long text string of the corresponding agreement ID.  How can my workload reach my microservice at runtime inside Horizon since I won't know that long name?**

When Horizon creates the private network for your microservcie and workload, it also applies a docker "--net-alias" with the short name, assuming you have specified that name in the "deployment" section of definition file for the container (i.e., the name you use for the container in the "services" array is deployed as its "--net-alias").


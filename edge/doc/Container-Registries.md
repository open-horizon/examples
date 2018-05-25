## How To Use The IBM Cloud Container Registry With Horizon

**Note:** not all of this functionality is in wiotp prod yet, so you can only do this in our bluehorizon hybrid environment at the moment. See [these instructions](https://docs.google.com/document/d/1_dIH79AKo_ngzbW9teE_x0sRfc8wpfikBqHugQQjKME/edit#heading=h.e9qakpsaxtok) for setting that up.

### Install and Use the IBM Cloud CLI Tool on x86 (Including Mac)
Begin by installing the IBM Cloud CLI (this must be done on an x86 host).  Instructions are available at the link below for Linux, MacOS, and Windows hosts: [https://console.bluemix.net/docs/cli/reference/bluemix_cli/get_started.html](https://console.bluemix.net/docs/cli/reference/bluemix_cli/get_started.html)  Homebrew for MacOS shortcut: `brew install homebrew/cask/ibm-cloud-cli`

### Login To IBM Cloud
Use the IBM Cloud CLI bx command login:

	`bx login`
	
(Note that IBM employees should instead login with:

	`bx login -sso)`

### Configure Organization
Once you have successfully logged-in to BlueMix, tell the IBM Cloud CLI which of your organization ID and space names you wish to use:

	`bx target -o <org ID> -s <space ID>`
	
For example:

	`bx target -o glendarling@us.ibm.com -s wiotp`
	
Or for our department's shared space (associated with Mike's master account):

	`bx target -o Hovitos -s dev`

If you don’t know your org ID and space name, use a web browser to login to the IBM cloud with the link below, and find your org ID and space at the top of that page:
	[https://console.bluemix.net/](https://console.bluemix.net/)

If it is your first time logging in to the IBM Cloud, you will need to create a new org and space.

### Install Container Registry Plugin
You will then need to add the BlueMix container registry plugin onto the IBM Cloud CLI, as follows:

	`bx plugin install container-registry -r Bluemix`
	
### Create Your Private Namespace
Use this container registry plugin to create a private namespace in the container registry under your account:

	`bx cr namespace-add <your namespace name>`
	
For example:

	`bx cr namespace-add glendarling`
	
Then finally, login to that Docker registry namespace by using the command:

	`bx cr login`

For more sub-commands, see the [bx cr CLI documentation](https://console.bluemix.net/docs/services/Registry/registry_cli.html) or run:

	`bx cr --help`

At this point you will be logged in to your private namespace in the IBM Cloud container registry and will be able to use registry paths like the one below in docker commands:

	`registry.ng.bluemix.net/<namespace>/<arch/image>:<version>`

### Create a Read/Write Token For Publishing Services and Patterns
The golang docker api package that hzn exchange service publish uses doesn't support the specific flavor of authentication that bx cr login uses. In addition, you can't use bx cr login on non-x86 systems if you are developing and publishing there. Both of these problems can be overcome by creating a read/write token for yourself. With that you can use docker login (instead of bx cr login) and then use all of the normal docker  and hzn commands.

To create the read/write token:

	`bx cr token-add --description '...' --non-expiring --readwrite`
	
The output should contain a line that begins with “Token”, then whitespace, then a long string of characters.  That long string is your token.  You can ignore everything else in that output.  To use this token with docker login, use a user name of “token” (i.e., literally those 5 characters) and then pass the token created above as the password:

	`docker login -u token -p 'WJh3NDJ2…' registry.ng.bluemix.net`
	

Or you could alternatively [create an API key](https://console.bluemix.net/docs/services/Registry/registry_tokens.html#registry_access).

### Create Read-Only Token For Horizon Users
To enable Horizon edge nodes to pull your service container images, you need to create a read-only token that Horizon can use on your edge nodes to access to your private namespace in the IBM Cloud Container Registry.  The bx command below will create such a token:

	`bx cr token-add --description '...' --non-expiring`

The output should contain a line that begins with “Token”, then whitespace, then a long string of characters.  That long string is your token.  You'll use this token in the next section.

### Build And Publish Your Images (On Any Machine)
Build your docker images as you normally would, and name them in this form:

	`registry.ng.bluemix.net/<your-namespace>/<arch>/<name>:<tag>`
	
When you are ready to publish your service to horizon, use the hzn command like this:

	`hzn exchange service publish -k <private-key-file> -K <public-key-file> -r "registry.ng.bluemix.net:<readonly-token>" -f service.definition.json`

That command will do several things:
Push the docker images to your ibm cloud container registry, and get the digest of the image in the process
Sign the digest and the deployment information with the private key/cert
Put the service metadata (including the signature) into the exchange
Put the key/cert into the exchange under the service definition so horizon edge nodes can automatically retrieve it to verify signatures when needed
Put your IBM Cloud container registry read-only token into the exchange under the service definition so horizon edge nodes can automatically retrieve it when needed

**Note:** this publish command is only available for services, not microservices or workloads.

### Using Your Service on Horizon Edge Nodes
The only thing special you or anyone else needs to do to use this service in a pattern on an edge node is to verify that /etc/horizon/anax.json contains these values and then restart anax:

	`"TrustCertUpdatesFromOrg": true,`
	`"TrustDockerAuthFromOrg": true,`
	`"ServiceUpgradeCheckIntervalS": 300`
    
(This is a temporary thing until the horizon-wiotp deb package sets these.)

After that create a pattern for your gw type that includes your service, and register your edge node with it. The edge node Horizon agent will automatically get the readonly token for this service from the exchange, and do the equivalent of docker login so it can access the service images.

If for debug you want to manually log your edge node into the container registry, run:

	`docker login -u token -p <token> registry.ng.bluemix.net`

### Additional Commands and Further Reading
Once you have push your images to the IBM Cloud registry, you can list all your images with this command:

	`bx cr images`
	
The image-rm command may be used to remove images from your namespace.

For more information, please see [using API keys and tokens with bx cr](https://console.bluemix.net/docs/services/Registry/registry_tokens.html#registry_access).

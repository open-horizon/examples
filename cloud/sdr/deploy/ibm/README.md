# sdr-auto

> Deployment automation for SDR PoC cloud part on IBM Cloud

## Overview

The tool helps you to create your deployment of SDR Cloud services. It creates and configures the following IBM Cloud services:

- Watson Speech to Text;

- Watson Natural Language Understanding;

- Compose for PostgreSQL;

- Event Streams;

- Functions;

- Cloud Foundry application.

## Prerequisites

The automation tool requires you to have the following prerequisites:

- `make`;
- `curl`;
- `jq`;
- IBM Cloud CLI `ibmcloud`;
- The interactive terminal for working with Postgres `psql`;
- `sed`;
- `grep`;
- `cut`;
- `npm` the package manager (v. 6.x);
- Node.js (v. 8.x or 10.x);
- a [token](https://www.mapbox.com/help/how-access-tokens-work/) from Mapbox;

You also need to be logged in with IBM Cloud CLI and target your CLI for an organization and space where you would like the tool to deploy and use services. Please make sure you have necessary permissions to create and access all required IBM Cloud services for SDR.

### Prerequisites installation

If you don't have installed, on Ubuntu you can install the prerequisites using the following commands:

- Command line utilities:

```
apt-get update
apt-get install -y make curl jq
```

- [IBM CLoud CLI](https://console.bluemix.net/docs/cli/index.html#overview):

```
curl -sL https://ibm.biz/idt-installer | bash
```

To login to your account and target specific organization and space, please run:

> You can create your API by following [the instructions](https://console.bluemix.net/docs/iam/userid_keys.html#userapikey).

```
export IC_PLATFORM_KEY='<your_api_key>'
ibmcloud login -a api.ng.bluemix.net -o <your_organization> -s <your_space> --apikey $IC_PLATFORM_KEY
```

- The interactive PostgreSQL terminal:

```
apt-get install -y postgresql-client
```

- Nodejs and npm manager by following [the instructions](https://www.digitalocean.com/community/tutorials/how-to-install-node-js-on-ubuntu-18-04) and issuing:

```
curl -sL https://deb.nodesource.com/setup_8.x -o nodesource_setup.sh
bash nodesource_setup.sh
apt update
apt install nodejs
apt install build-essential
```

> A build step for the UI app requires extra memory. We found that it works fine on a VM with 8 GB of RAM.

## Configuration

You can either use default values for services or define your own. Whatever you choose, in order to avoid any collisions with other working services in your account, please make sure there’s no other services with the same names you used with the tool.

### Quick setup

1. Please follow the instructions for prerequisites install.

2. For a quick setup, please specify a prefix `SERVICE_PREFIX` in the `deploy.sh` script to a value of your preference. All service and credentials names will start with it. Default is `sdr-poc`

3. Specify the following environment variables:

	- `UI_APP_USER` - email for login to the SDR UI application (i.e. user@mail.com)
	- `UI_APP_PASSWORD` - password for the SDR UI application
	- `MAPBOX_TOKEN` - [Mapbox token](https://www.mapbox.com/help/how-access-tokens-work/)

4. Check your configuration either with `make config` or `./deploy.sh --config`

5. Make sure you have all prerequisities ready with `make prereqs` or `./deploy.sh --prereqs`

6. Run `make deploy-all` or `./deploy.sh --install=all` to deploy all the SDR cloud services.

### Advanced setup

The `deploy.sh` script allows you to configure the following values:

1. Event Streams (former Message Hub) instance:

	- `MH_INSTANCE` - Event Streams instance name;
	- `MH_INSTANCE_CREDS` - Event Streams credentials name;
	- `MH_SDR_TOPIC` - SDR topic;
	- `MH_SDR_TOPIC_PARTIONS` - partitions number for the topic;
	- `MH_RESPONSE_RETRY` - time in seconds to retry when checking if the instance is ready;

2. Watson Speech-to-Text service:

	- `STT_INSTANCE` - service instance name;
	- `STT_INSTANCE_CREDS` - service credentials name;
	- `STT_INSTANCE_PLAN` - service plan (`lite`, **`standard`**, etc.);

3. Watson Natural Language Understanding service:

	- `NLU_INSTANCE` - service instance name;
	- `NLU_INSTANCE_CREDS` - service credentials name;
	- `NLU_INSTANCE_PLAN` - service plan (`lite`, **`standard`**, etc.);

4. Compose for PostgreSQL Database service:

	- `DB_INSTANCE` - database service instance name;
	- `DB_INSTANCE_CREDS` - service credentials name;
	- `DB_INSTANCE_PLAN` - service plan (**`Standard`**);
	- `DB_INSTANCE_RETRY` - check DB instance readiness interval, seconds;
	- `DB_INSTANCE_TIMEOUT` - timeout for DB instance readiness check;
	- `DB_NAME` - SDR database name (`sdr`);

5. Cloud Functions service:

	- `FUNC_PACKAGE` - package name for the action;
	- `FUNC_MH_FEED` - feed with messages from an Event Streams instance, please use default value;
	- `FUNC_TRIGGER` - trigger name;
	- `FUNC_ACTION` - action name;
	- `FUNC_ACTION_CODE` - path to the action code;
	- `FUNC_RULE` - rule name;

6. SDR UI application:

	- `UI_SRC_PATH` - path to the root of the Node.js UI app;
	- `UI_APP_NAME` - application name (aslo used in a URL generated for the application `UI_APP_NAME-<random-generated word>-<random-generated-word>.mybluemix.net`)

## Usage

1. Using `make`.

The `Makefile` has the following targets which help to create services:

```
make help
help				: Display help
 prereqs			: Check prerequisites
 config				: Display current configuration
 deploy-es			: Create and configure Event Streams Instance
 teardown-es		: Delete Event Streams instance
 deploy-db			: Create and configure Compose for PostgreSQL instance
 teardown-db		: Delete Compose for PostgreSQL instance
 deploy-stt			: Create Watson Speech-To-Text instance
 teardown-stt		: Delete Watson Speech-To-text instance
 deploy-nlu			: Create Watson Natural Language Understanding instance
 teardown-nlu		: Delete Watson Natural Language Understanding instance
 deploy-func		: Create and configure functions
 teardown-func		: Delete functions
 deploy-ui			: Deploy UI application
 teardown-ui		: Delete UI application
 deploy-all			: Create all instances
 teardown-all		: Delete all instances
 ```
 
`make deploy-all` - to deploy the cloud part of SDR PoC, use.

`make teardown-all` - to delete all the instances of the SDR PoC cloud part.

If a target creates a service that depends on other services (e.g. Cloud Functions `deploy-func`, UI app `deploy-ui`), it also deploys those service it depends on first. Targets that don’t depend on other services just create a single service without any dependencies (e.g. Watson STT `deploy-stt`, Watson NLU `deploy-nlu`, Event Streams `deploy-es`, Compose for PostgreSQL `deploy-db`).

On the other hand, if a target tears down a service (Watson STT `teardown-stt`, Watson NLU `teardown-nlu`, Event-Streams `teardown-es`, Compose for PostgreSQL `teardown-db`) other SDR parts depend on, it deletes those dependant services first.

2. Using the `deploy.sh` script directly.

```
deploy.sh [ [-i|--install=<component>] || [-u|--uninstall=<component>]] -- deploying cloud part for SDR PoC

where:
	-p | --prereqs					- check for prerequisites
	-c | --config					- show current configuration
	-i= | --install=[component]			- install [component]
	-u= | --uninstall=[component]			- uninstall [component]

Example: ./deploy.sh --install=all

Supported components are:
stt nlu db es func ui all
```

Some examples:

`./deploy.sh --install=all` - create all SDR PoC services

`./deploy.sh --uninstall=all` - delete all SDR PoC services

`./deploy.sh --install=db` - create DB
# sdr-auto

> Deployment automation for SDR PoC cloud part on IBM Cloud

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
- `npm` the package manager;
- Node.js;
- define environment variable `MAPBOX_TOKEN` with a token from Mapbox;

You also need to be logged in with IBM Cloud CLI and targeted your CLI for organization and space where you would like it to deploy to and use services from.

## Configuration

You can either use default values for services or define your own. Whatever you choose, in order to avoid any collisions with other working services in your account, please make sure there’s no other services with the same names you used with the tool.

The `Makefile` allows you to configure the following values:

1. Event Streams (former Message Hub) instance:

	- `MH_INSTANCE` - Event Streams instance name;
	- `MH_INSTANCE_CREDS` - Event Streams credentials name;
	- `MH_SDR_TOPIC` - SDR topic;
	- `MH_SDR_TOPIC_PARTIONS` - partitions number for the topic;
	- `MH_RESPONSE_RETRY` - time in seconds to retry when checking if the instance is ready;

2. Watson Speech-to-Text service:

	- `STT_INSTANCE` - service instance name;
	- `STT_INSTANCE_CREDS` - service credentials name;
	- `STT_INSTANCE_PLAN` - service plan (**`lite`**, `standard`, etc.);

3. Watson Natural Language Understanding service:

	- `NLU_INSTANCE` - service instance name;
	- `NLU_INSTANCE_CREDS` - service credentials name;
	- `NLU_INSTANCE_PLAN` - service plan (**`lite`**, `standard`, etc.);

4. Compose for PostgreSQL Database service:

	- `DB_INSTANCE` - database service instance name;
	- `DB_INSTANCE_CREDS` - service credentials name;
	- `DB_INSTANCE_PLAN` - service plan (**`Standard`**);
	- `DB_INSTANCE_RETRY` - check DB instance readiness interval, seconds;
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

The `Makefile` has the following targets which help to create services:

```
make help
 help				: Display help
 prereqs			: Check for prerequisites
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



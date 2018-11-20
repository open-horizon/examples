.PHONY:	help deploy-stt teardown-stt deploy-nlu teardown-nlu \
		deploy-db teardown-db deploy-all teardown-all prereqs config \
		deploy-es teardown-es deploy-func teardown-func deploy-ui teardown-ui

SHELL:=bash

## help				: Display help
help : Makefile
	@sed -n 's/^##//p' $<

## prereqs			: Check prerequisites
prereqs :
	./deploy.sh --prereqs

## config				: Display current configuration
config:
	./deploy.sh --config

## deploy-es			: Create and configure Event Streams Instance
deploy-es :
	./deploy.sh --install=es

## teardown-es			: Delete Event Streams instance
teardown-es :
	./deploy.sh --uninstall=es

## deploy-db			: Create and configure Compose for PostgreSQL instance
deploy-db : 
	./deploy.sh --install=db

## teardown-db			: Delete Compose for PostgreSQL instance
teardown-db :
	./deploy.sh --uninstall=db

## deploy-stt			: Create Watson Speech-To-Text instance
deploy-stt :
	./deploy.sh --install=stt

## teardown-stt			: Delete Watson Speech-To-text instance
teardown-stt :
	./deploy.sh --uninstall=stt

## deploy-nlu			: Create Watson Natural Language Understanding instance
deploy-nlu :
	./deploy.sh --install=nlu

## teardown-nlu			: Delete Watson Natural Language Understanding instance
teardown-nlu :
	./deploy.sh --uninstall=nlu

## deploy-func			: Create and configure functions
deploy-func :
	./deploy.sh --install=func

## teardown-func			: Delete functions
teardown-func :
	./deploy.sh --uninstall=func

## deploy-ui			: Deploy UI application
deploy-ui :
	./deploy.sh --install=ui

## teardown-ui			: Delete UI application
teardown-ui :
	./deploy.sh --uninstall=ui

## deploy-all			: Create all instances
deploy-all :
	./deploy.sh --install=all

## teardown-all			: Delete all instances
teardown-all :
	./deploy.sh --uninstall=all


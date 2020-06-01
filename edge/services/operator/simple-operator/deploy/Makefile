# Make targets for building the IBM example operator edge service

# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them
-include horizon/.hzn.json.tmp.mk

# Default ARCH to the architecture of this machines (as horizon/golang describes it)
export ARCH ?= $(shell hzn architecture)

# Configurable parameters passed to serviceTest.sh in "test" target
export MATCH:='says: Hello'
export TIME_OUT:=20

# Target for travis to test new PRs
test-all-arches: test

# Run and verify the service
test:
	hzn dev service verify

# Publish the service to the Horizon Exchange for the current architecture
publish-service:
	hzn exchange service publish -O -f horizon/service.definition.json

# Target for travis to publish service and pattern after PR is merged  
publish: publish-service publish-service-policy publish-deployment-policy
	hzn exchange pattern publish -f horizon/pattern.json

# Build, run and verify, if test succeeds then publish (for the current architecture)
test-publish: test publish-service

# Build/test/publish the service to the Horizon Exchange for 3 architectures and publish a deployment pattern for those architectures
publish-all-arches: test-publish
	hzn exchange pattern publish -f horizon/pattern.json

# target for script - overwrite and pull insitead of push docker image
publish-service-overwrite:
	hzn exchange service publish -O -P -f horizon/service.definition.json

# Publish Service Policy target for exchange publish script
publish-service-policy:
	hzn exchange service addpolicy -f horizon/service.policy.json $(HZN_ORG_ID)/$(SERVICE_NAME)_$(SERVICE_VERSION)_$(ARCH)

# Publish Business Policy target for exchange publish script
publish-deployment-policy:
	hzn exchange deployment addpolicy -f horizon/deployment.policy.json $(HZN_ORG_ID)/$(SERVICE_NAME)_$(SERVICE_VERSION)

# new target for icp exchange to run on startup to publish only
publish-only: publish-service-overwrite publish-service-policy
	hzn exchange pattern publish -f horizon/pattern.json

# This imports the variables from horizon/hzn.cfg. You can ignore these lines, but do not remove them.
horizon/.hzn.json.tmp.mk: horizon/hzn.json
	@ hzn util configconv -f $< | sed 's/=/?=/' > $@

.PHONY: build build-all-arches test publish-service build-test-publish publish-all-arches clean clean-all-archs

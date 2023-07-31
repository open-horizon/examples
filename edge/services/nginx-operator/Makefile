# Make targets for building the IBM example operator edge service

# This imports the variables from horizon/hzn.json. You can ignore these lines, but do not remove them
-include horizon/.hzn.json.tmp.mk

# Default ARCH to the architecture of this machines (as horizon/golang describes it)
export ARCH ?= amd64
export SERVICE_VERSION ?= 2.0.0

# extra helper operator env vars
OPERATOR_KIND_LC = $(shell echo $(OPERATOR_KIND) | tr A-Z a-z)
OPERATOR_IMAGE = $(DOCKER_HUB_ID)/$(OPERATOR_GROUP_NAME)_$(ARCH):$(SERVICE_VERSION)

init:
	if [ ! -d "$(OPERATOR_GROUP_NAME)" ]; then \
		echo ""; \
		echo "${BGREEN}Initializing files using Operator SDK in $(OPERATOR_GROUP_NAME) ... ${NC}"; \
		mkdir "$(OPERATOR_GROUP_NAME)"; \
		cd "$(OPERATOR_GROUP_NAME)"; \
		operator-sdk init --plugins $(OPERATOR_TYPE) --domain "$(OPERATOR_DOMAIN)" --group $(OPERATOR_GROUP_NAME) --version "$(OPERATOR_API_VERSION)" --kind $(OPERATOR_KIND) --generate-role ; \
		sed -i -e 's/env:/env:\n        - name: WATCH_NAMESPACE\n          valueFrom:\n            fieldRef:\n              fieldPath: metadata.namespace/' ./config/manager/manager.yaml ; \
		sed -i -z -e 's/\s\{2\}- apiGroups:\n\s\{6\}- ""\n\s\{4\}resources:\n/  - apiGroups:\n      - ""\n    resources:\n      - services\n/' ./config/rbac/role.yaml ; \
		sed -i -z -e 's/\s\{2\}- apiGroups:\n\s\{6\}- ""\n\s\{4\}resources:\n/  - apiGroups:\n      - ""\n    resources:\n      - routes\n/' ./config/rbac/role.yaml ; \
		sed -i -e 's/metadata:/metadata:\n  namespace: $(OPERATOR_NAMESPACE)/' "./config/samples/$(OPERATOR_GROUP_NAME)_$(OPERATOR_API_VERSION)_$(OPERATOR_KIND_LC).yaml" ; \
		sed -i -e 's/spec:/spec:\n  size: 1/' "./config/samples/$(OPERATOR_GROUP_NAME)_$(OPERATOR_API_VERSION)_$(OPERATOR_KIND_LC).yaml" ; \
		cd config/default; kustomize edit set namespace $(OPERATOR_NAMESPACE) ; \
	else \
		echo ""; \
		echo "${BGREEN}Initialized files already exists in $(OPERATOR_GROUP_NAME) ... ${NC}"; \
	fi

nginx-files:
	echo ""; \
	echo "${BGREEN}Gathering the nginx deployment and service files... ${NC}"; \
	cd "$(OPERATOR_GROUP_NAME)"; \
	wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/nginx-operator/ansible-role-files/deployment.j2; mv deployment.j2 roles/$(OPERATOR_KIND_LC)/templates/. ; \
	wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/nginx-operator/ansible-role-files/service.j2; mv service.j2 roles/$(OPERATOR_KIND_LC)/templates/. ; \
	wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/nginx-operator/ansible-role-files/route.j2; mv route.j2 roles/$(OPERATOR_KIND_LC)/templates/. ; \
	wget https://raw.githubusercontent.com/open-horizon/examples/master/edge/services/nginx-operator/ansible-role-files/main-route.yml; mv main-route.yml roles/$(OPERATOR_KIND_LC)/tasks/main.yml

build:
	@echo ""
	@echo "${BGREEN}Building operator docker image... ${NC}"; \
	cd $(OPERATOR_GROUP_NAME); make docker-build IMG="$(OPERATOR_IMAGE)"

push:
	@echo ""
	@echo "${BGREEN}Pushing operator image into docker... ${NC}"; \
	cd $(OPERATOR_GROUP_NAME); make docker-push IMG="$(OPERATOR_IMAGE)"

fix:
	@echo ""
	@echo "${BGREEN}Creating operator-tar to publish in IEAM ... ${NC}"; \
	cd $(OPERATOR_GROUP_NAME); sed -i -e 's/namespace: $(OPERATOR_NAMESPACE)/namespace: openhorizon-agent/' "./config/samples/$(OPERATOR_GROUP_NAME)_$(OPERATOR_API_VERSION)_$(OPERATOR_KIND_LC).yaml" ; \
	cd config/default; kustomize edit set namespace openhorizon-agent

deploy:
	@echo ""
	@echo "${BGREEN}Deploying operator... ${NC}"; \
	cd $(OPERATOR_GROUP_NAME); make deploy IMG="$(OPERATOR_IMAGE)"; \
	kubectl apply -f config/samples/$(OPERATOR_GROUP_NAME)_$(OPERATOR_API_VERSION)_$(OPERATOR_KIND_LC).yaml

undeploy:
	@echo ""
	@echo "${BGREEN}Deleting operator... ${NC}"; \
	cd $(OPERATOR_GROUP_NAME); kubectl delete -f config/samples/$(OPERATOR_GROUP_NAME)_$(OPERATOR_API_VERSION)_$(OPERATOR_KIND_LC).yaml; \
	make undeploy IMG="$(OPERATOR_IMAGE)"

# Save operator tar in horizon directory and update operatorYamlArchive name
tar:
	@echo ""
	@echo "${BGREEN}Creating operator-tar to publish in IEAM ... ${NC}"; \
	sed -i -e 's/namespace: $(OPERATOR_NAMESPACE)/namespace: openhorizon-agent/' "./$(OPERATOR_GROUP_NAME)/config/samples/$(OPERATOR_GROUP_NAME)_$(OPERATOR_API_VERSION)_$(OPERATOR_KIND_LC).yaml" ; \
	mkdir -p tmp; \
	cd $(OPERATOR_GROUP_NAME)/config/default; kustomize edit set namespace openhorizon-agent ; cd ../ ; \
	cd manager/ && kustomize edit set image controller="$(OPERATOR_IMAGE)" ; cd ../../.. ; \
	kustomize build $(OPERATOR_GROUP_NAME)/config/default > tmp/$(OPERATOR_GROUP_NAME).yaml; \
	cp $(OPERATOR_GROUP_NAME)/config/samples/$(OPERATOR_GROUP_NAME)_$(OPERATOR_API_VERSION)_$(OPERATOR_KIND_LC).yaml tmp/. ; \
	tar -zcvf operator.tar.gz tmp/* ; rm -r tmp; 

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
	hzn exchange service publish -O -P --public=true -f horizon/service.definition.json

# Publish Service Policy target for exchange publish script
publish-service-policy:
	hzn exchange service addpolicy -f horizon/service.policy.json $(HZN_ORG_ID)/$(SERVICE_NAME)_$(SERVICE_VERSION)_$(ARCH)

# Publish Deployment Policy target for exchange publish script
publish-deployment-policy:
	hzn exchange deployment addpolicy -f horizon/deployment.policy.json $(HZN_ORG_ID)/policy-$(SERVICE_NAME)_$(SERVICE_VERSION)

# new target for icp exchange to run on startup to publish only
publish-only:
	ARCH=amd64 $(MAKE) publish-service-overwrite
	ARCH=amd64 $(MAKE) publish-service-policy
	ARCH=ppc64le $(MAKE) publish-service-overwrite
	ARCH=ppc64le $(MAKE) publish-service-policy
	hzn exchange pattern publish -f horizon/pattern-all-arches.json

# This imports the variables from horizon/hzn.cfg. You can ignore these lines, but do not remove them.
horizon/.hzn.json.tmp.mk: horizon/hzn.json
	@ hzn util configconv -f $< | sed 's/=/?=/' > $@

.PHONY: build build-all-arches test publish-service build-test-publish publish-all-arches clean clean-all-archs

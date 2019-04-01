# Make targets for building the IBM example helloworld edge service

# This imports the variables from horizon/hzn.env. You can ignore these lines, but do not remove them.
TMP_MK=horizon/hzn.env.tmp.mk
-include $(TMP_MK)


build:
	docker build -t $(DOCKER_IMAGE) -f ./Dockerfile.$(ARCH) .
	@rm -f $(TMP_MK)

build-all-arches:
	ARCH=amd64 $(MAKE) build
	ARCH=arm $(MAKE) build
	ARCH=arm64 $(MAKE) build
	@rm -f $(TMP_MK)

publish-service:
	hzn exchange service publish -f horizon/service.definition.json -k $(HZN_PRIVATE_KEY_FILE) -K $(HZN_PUBLIC_KEY_FILE)
	@rm -f $(TMP_MK)

publish-all-arches:
	ARCH=amd64 $(MAKE) publish-service
	ARCH=arm $(MAKE) publish-service
	ARCH=arm64 $(MAKE) publish-service
	hzn exchange pattern publish -f horizon/pattern/pattern-helloworld-all-arches.json -p pattern-helloworld
	@rm -f $(TMP_MK)

clean:
	-docker rmi $$DOCKER_IMAGE 2> /dev/null || :
	@rm -f $(TMP_MK)

clean-all-archs:
	ARCH=amd64 $(MAKE) clean
	ARCH=arm $(MAKE) clean
	ARCH=arm64 $(MAKE) clean
	@rm -f $(TMP_MK)

# Just to experiment with make
testmake:
	@echo "testmake: $(ARCH)"
	@rm -f $(TMP_MK)

# This imports the variables from horizon/hzn.env. You can ignore these lines, but do not remove them.
# Note: this created artifact needs to be deleted before make is run again, in case the user wants to override a variable like ARCH
$(TMP_MK): FORCE
	@echo 'set -x' | cat - horizon/hzn.env | bash 2>&1 | sed -e 's/^+ export /export /' -e t -e d | sed "s/'//g" > $@

FORCE:

.PHONY: build clean

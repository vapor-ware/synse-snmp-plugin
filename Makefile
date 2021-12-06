#
# Synse SNMP Plugin
#

PLUGIN_NAME    := snmp
PLUGIN_VERSION := 2.2.3
IMAGE_NAME     := vaporio/snmp-plugin
BIN_NAME       := synse-snmp-plugin

GIT_COMMIT     ?= $(shell git rev-parse --short HEAD 2> /dev/null || true)
GIT_TAG        ?= $(shell git describe --tags 2> /dev/null || true)
BUILD_DATE     := $(shell date -u +%Y-%m-%dT%T 2> /dev/null)
GO_VERSION     := $(shell go version | awk '{ print $$3 }')

PKG_CTX := github.com/vapor-ware/synse-sdk/sdk
LDFLAGS := -w \
	-X ${PKG_CTX}.BuildDate=${BUILD_DATE} \
	-X ${PKG_CTX}.GitCommit=${GIT_COMMIT} \
	-X ${PKG_CTX}.GitTag=${GIT_TAG} \
	-X ${PKG_CTX}.GoVersion=${GO_VERSION} \
	-X ${PKG_CTX}.PluginVersion=${PLUGIN_VERSION}


.PHONY: build
build:  ## Build the plugin binary
	go build -ldflags "${LDFLAGS}" -o ${BIN_NAME} || exit

.PHONY: build-linux
build-linux:  ## Build the plugin binarry for linux amd64
	GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o ${BIN_NAME} . || exit

.PHONY: clean
clean:  ## Remove temporary files
	go clean -v || exit
	rm -rf dist

.PHONY: deploy
deploy:  ## Run a local deployment of the plugin with Synse Server
	docker-compose -f compose.yml up -d || exit

.PHONY: dep
dep:  ## Verify and tidy gomod dependencies
	go mod verify || exit
	go mod tidy || exit

.PHONY: docker
docker:  ## Build the production docker image locally
	docker build -f Dockerfile \
		--label "org.label-schema.build-date=${BUILD_DATE}" \
		--label "org.label-schema.vcs-ref=${GIT_COMMIT}" \
		--label "org.label-schema.version=${PLUGIN_VERSION}" \
		-t ${IMAGE_NAME}:latest . || exit

.PHONY: docker-dev
docker-dev:  ## Build the development docker image locally
	docker build -f Dockerfile.dev -t ${IMAGE_NAME}:dev-${GIT_COMMIT} . || exit

.PHONY: fmt
fmt:  ## Run goimports on all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file" || exit ; done

.PHONY: github-tag
github-tag:  ## Create and push a tag with the current plugin version
	git tag -a ${PLUGIN_VERSION} -m "${PLUGIN_NAME} plugin version ${PLUGIN_VERSION}"
	git push -u origin ${PLUGIN_VERSION}

.PHONY: lint
lint:  ## Lint project source files
	golint -set_exit_status ./pkg/... || exit

.PHONY: version
version:  ## Print the version of the plugin
	@echo "${PLUGIN_VERSION}"

.PHONY: help
help:  ## Print usage information
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.DEFAULT_GOAL := help

.PHONY: integration-test
unit-test: test

.PHONY: start-snmp-emulators
start-snmp-emulators: ## Start emulators for functional tests.
	# Start the SNMP emulator for the pxgmsups (Eaton UPS) in a docker container in the background.
	# Tests run on the local machine.
	docker-compose -f ./emulator/ups/pxgms_ups/test_snmp.yml down || true
	docker-compose -f ./emulator/ups/pxgms_ups/test_snmp.yml build
	docker-compose -f ./emulator/ups/pxgms_ups/test_snmp.yml up -d
	# Start the SNMP emulator for the pxgmsups (Eaton UPS) in a docker container in the background.
	# Tests run on the local machine.
	docker-compose -f ./emulator/ups/tripplite_ups/test_snmp.yml down || true
	docker-compose -f ./emulator/ups/tripplite_ups/test_snmp.yml build
	docker-compose -f ./emulator/ups/tripplite_ups/test_snmp.yml up -d

.PHONY: stop-snmp-emulators
stop-snmp-emulators: ## Shutdown the emulators.
	docker-compose -f ./emulator/ups/tripplite_ups/test_snmp.yml down
	docker-compose -f ./emulator/ups/pxgms_ups/test_snmp.yml down

.PHONY: run-tests
run-tests: ## Run the tests. Requires the emulators to be up.
	go test -cover -v ./... || (echo TESTS FAILED $$?; docker-compose -f ./emulator/ups/pxgms_ups/test_snmp.yml kill; exit 1)

.PHONY: test-dev-box ## Start the emulators on your dev box, run tests, stop the emulators.
test-dev-box: start-snmp-emulators run-tests stop-snmp-emulators  ## Start emulators, run all tests, stop emulators.

.PHONY: test
test: run-tests ## CI hooks into this. CI starts its own emulators.

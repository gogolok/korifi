
# Image URL to use all building/pushing image targets
IMG_SSR ?= cloudfoundry/korifi-statefulset-runner:latest
# ENVTEST_K8S_VERSION refers to the version of kubebuilder assets to be downloaded by envtest binary.
ENVTEST_K8S_VERSION = 1.23
CLUSTER_NAME ?= "e2e"

# Setting SHELL to bash allows bash commands to be executed by recipes.
# This is a requirement for 'setup-envtest.sh' in the test target.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
export GOBIN = $(shell pwd)/bin
export PATH := $(shell pwd)/bin:$(PATH)

webhooks-file = ../helm/korifi/statefulset-runner/manifests.yaml
.PHONY: manifests
manifests: bin/controller-gen bin/yq
	controller-gen \
		paths="./..." \
		webhook \
		rbac:roleName=korifi-statefulset-runner-appworkload-manager-role \
		output:rbac:artifacts:config=../helm/korifi/statefulset-runner \
		output:webhook:artifacts:config=../helm/korifi/statefulset-runner

	yq -i 'with(.metadata; .annotations["cert-manager.io/inject-ca-from"]="{{ .Release.Namespace }}/korifi-statefulset-runner-serving-cert")' $(webhooks-file)
	yq -i 'with(.metadata; .name="korifi-statefulset-runner-" + .name)' $(webhooks-file)
	yq -i 'with(.webhooks[]; .clientConfig.service.namespace="{{ .Release.Namespace }}")' $(webhooks-file)
	yq -i 'with(.webhooks[]; .clientConfig.service.name="korifi-statefulset-runner-" + .clientConfig.service.name)' $(webhooks-file)

.PHONY: generate
generate: bin/controller-gen
	controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."

.PHONY: test
test: manifests generate
	../scripts/run-tests.sh

bin:
	mkdir -p bin

bin/controller-gen: bin
	go install sigs.k8s.io/controller-tools/cmd/controller-gen

bin/yq: bin
	go install github.com/mikefarah/yq/v4@latest

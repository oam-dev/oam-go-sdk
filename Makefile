
# Image URL to use all building/pushing image targets
IMG ?= controller:latest
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

all: manager

# Run tests
test: generate fmt vet manifests
	go test ./apis/... ./pkg/... -coverprofile cover.out

# Build manager binary
examples: generate fmt vet
	go build -o bin/example1 pkg/examples/framework/main.go
	go build -o bin/example2 pkg/examples/extendworkload/main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	go run ./main.go

# Install CRDs into a cluster
install: manifests
	kubectl apply -f config/crd/bases

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests
	kubectl apply -f config/crd/bases
	kustomize build config/default | kubectl apply -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./apis/core.oam.dev/...;./controllers/..." output:crd:artifacts:config=config/crd/bases

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths=./apis/core.oam.dev/...
	./hack/update-client-gen.sh

# Build the docker image
docker-build:
	go mod vendor
	docker build . -t ${IMG}

# Push the docker image to docker hub
# Usage: make docker-build docker-push IMG=fireeye2018/hydra-controller:latest
docker-push:
	docker push ${IMG}

# find or download controller-gen
# download controller-gen if necessary
# make sure use controller-gen v0.2.1
controller-gen:
ifeq (, $(shell which controller-gen))
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.4
CONTROLLER_GEN=$(shell go env GOPATH)/bin/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
# Copyright © 2023 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PKG = github.com/openshift/lvm-driver
EXE_NAME = lvmdriver
GIT_COMMIT = $(shell git rev-parse HEAD)
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
DRIVER_VERSION = v0.0.0
LDFLAGS = -X ${PKG}/pkg/lvmdriver.driverVersion=${DRIVER_VERSION} -X ${PKG}/pkg/lvmdriver.gitCommit=${GIT_COMMIT} -X ${PKG}/pkg/lvmdriver.buildDate=${BUILD_DATE}
OS ?= linux
ARCH ?= amd64

# IMAGE Build info 
IMAGE_REGISTRY ?= quay.io
REGISTRY_NAMESPACE ?= lvms_dev
IMAGE_TAG ?= $(GIT_COMMIT)
IMAGE_NAME ?= lvm-driver
IMAGE_REPO ?= $(IMAGE_REGISTRY)/$(REGISTRY_NAMESPACE)/$(IMAGE_NAME)
IMG ?= $(IMAGE_REPO):$(IMAGE_TAG)
IMAGE_BUILD_CMD ?= $(shell command -v docker 2>&1 >/dev/null && echo docker || echo podman)

all: build

.PHONY: build
build: 
	GOOS=$(OS) GOARCH=$(ARCH) go build -a -ldflags "${LDFLAGS}" -o ./bin/lvm_driver ./cmd/lvmdriver/main.go

.PHONY: container
container:
	$(IMAGE_BUILD_CMD) build --platform=${OS}/${ARCH} -t ${IMG} .
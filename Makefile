# Copyright 2016 Red Hat, Inc.
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

SHELL := /bin/bash -euo pipefail

# Use the native vendor/ dependency system
export GO111MODULE := on

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

DOCKER_REPO ?= docker.io/kaasops
DOCKER_IMAGE_NAME ?= config-reloader
DOCKER_IMAGE_TAG ?= develop
BINARY=config-reloader

LDFLAGS := -extldflags '-static'

.PHONY: build-local
build-local: clean
	GOARCH=$(GOARCH) GOOS=$(GOOS) CGO_ENABLED=0 \
		go build -ldflags="$(LDFLAGS)" -o out/$(BINARY) cmd/configreloader/main.go

.PHONY: build
build: clean
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 \
		go build -ldflags="$(LDFLAGS)" -o out/$(BINARY) cmd/configreloader/main.go

.PHONY: run
run: build-local
	./out/$(BINARY)

.PHONY: clean
clean:
	rm -rf ./out

.PHONY: docker
docker:
	docker build -t $(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

.PHONY: push
push:
	docker push $(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(TAG)
SHELL := /bin/bash -euo pipefail

# Use the native vendor/ dependency system

DOCKER_REPO ?= docker.io/kaasops
DOCKER_IMAGE_NAME ?= config-reloader
DOCKER_IMAGE_TAG ?= develop

.PHONY: docker
docker:
	docker build -t $(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

.PHONY: push
push:
	docker push $(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(TAG)
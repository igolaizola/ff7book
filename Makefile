#!/bin/bash

SHELL             = /bin/bash
PLATFORMS        ?= linux/amd64 darwin/amd64 darwin/arm64 windows/amd64
IMAGE_PREFIX     ?= igolaizola
REPO_NAME        ?= ff7book
COMMIT_SHORT     ?= $(shell git rev-parse --verify --short HEAD)
VERSION          ?= $(shell git describe --tags --exact-match 2>/dev/null || echo $(COMMIT_SHORT))

# Build the binaries for the current platform
.PHONY: build
build:
	os=$$(go env GOOS); \
	arch=$$(go env GOARCH); \
	PLATFORMS="$$os/$$arch" make app-build

# Build the binaries
# Example: PLATFORMS=linux/amd64 make app-build
.PHONY: app-build
app-build:
	@for platform in $(PLATFORMS) ; do \
		os=$$(echo $$platform | cut -f1 -d/); \
		arch=$$(echo $$platform | cut -f2 -d/); \
		arm=$$(echo $$platform | cut -f3 -d/); \
		arm=$${arm#v}; \
		ext=""; \
		if [ "$$os" == "windows" ]; then \
			ext=".exe"; \
		fi; \
		file=./bin/$(REPO_NAME)-$(COMMIT_SHORT)-$$(echo $$platform | tr / -)$$ext; \
		GOOS=$$os GOARCH=$$arch GOARM=$$arm CGO_ENABLED=0 \
		go build \
			-a -x -tags netgo,timetzdata -installsuffix cgo -installsuffix netgo \
			-ldflags " \
				-X main.version=$(VERSION) \
				-X main.commit=$(COMMIT_SHORT) \
				-X main.date=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
			" \
			-o $$file \
			./cmd/$(REPO_NAME); \
		if [ $$? -ne 0 ]; then \
			exit 1; \
		fi; \
		chmod +x $$file; \
	done

# Build the docker image
# Example: PLATFORMS=linux/amd64 make docker-build
.PHONY: docker-build
docker-build:
	rm -rf bin; \
	@platforms=($(PLATFORMS)); \
	platform=$${platforms[0]}; \
	if [[ $${#platforms[@]} -ne 1 ]]; then \
    	echo "Multi-arch build not supported"; \
		exit 1; \
	fi; \
	docker build --platform $$platform -t $(IMAGE_PREFIX)/$(REPO_NAME):$(COMMIT_SHORT) .; \
	if [ $$? -ne 0 ]; then \
		exit 1; \
	fi

# Build the docker images using buildx
# Example: PLATFORMS="linux/amd64 darwin/amd64 windows/amd64" make docker-buildx
.PHONY: docker-buildx
docker-buildx:
	@platforms=($(PLATFORMS)); \
	platform=$$(IFS=, ; echo "$${platforms[*]}"); \
	docker buildx build --platform $$platform -t $(IMAGE_PREFIX)/$(REPO_NAME):$(COMMIT_SHORT) .

# Clean binaries
.PHONY: clean
clean:
	rm -rf bin/README.*
	rm -rf bin/$(REPO_NAME)-*

# Zip the binaries
.PHONY: zip
zip: clean build
	cp README.pdf bin/; \
	cd bin; \
	zip -r $(REPO_NAME)-$(shell date -u +'%Y%m%d-%H%M').zip README.pdf $(REPO_NAME)-*; \

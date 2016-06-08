# Project
PROJECT?=tasks
ORGANIZATION?=czertbytes
REPOSITORY?=github.com
DOCKER_SOURCE?=Docker

# Build
BUILDER_IMAGE?=czertbytes/golang-builder:latest
GO_BUILD_PARAMS?=-a -installsuffix -cgo
GO_BUILD_CMD?=go build $(GO_BUILD_PARAMS) -o $(DOCKER_SOURCE)/bin/$(PROJECT) cmd/tasks/main.go

# Build Linux
GO_BUILD_LINUX_PARAMS?=-a -installsuffix -cgo
GO_BUILD_LINUX_CMD?=go build $(GO_BUILD_LINUX_PARAMS) -o $(DOCKER_SOURCE)/bin/$(PROJECT)-linux cmd/tasks/main.go

BUILD_TAGS?=latest

# Test packages
# ignores "can't load package: package" errors produced by go list command
GO_TEST_PACKAGES?=$(shell go list ./... 2>/dev/null | grep -v /vendor/)

all: test image

build:
	mkdir -p $(DOCKER_SOURCE)/bin
	$(GO_BUILD_CMD)

build-linux:
	mkdir -p $(DOCKER_SOURCE)/bin
	docker run --rm \
		-v $(PWD):/go/src/$(REPOSITORY)/$(ORGANIZATION)/$(PROJECT) \
		-w /go/src/$(REPOSITORY)/$(ORGANIZATION)/$(PROJECT) \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-e CGO_ENABLED=1 \
		$(BUILDER_IMAGE) \
		$(GO_BUILD_LINUX_CMD)

test:
	go test -v -timeout 60s $(GO_TEST_PACKAGES)

image: build-linux
	docker build -t $(ORGANIZATION)/$(PROJECT):latest $(DOCKER_SOURCE)

run:
	docker-compose up

clean:
	rm -rf $(DOCKER_SOURCE)/bin

.PHONY: all

REPO_NAME ?= mon

DOCKER_IMAGE_TAG ?=
DOCKER_IMAGE ?= monserve:$(DOCKER_IMAGE_TAG)

TEST_BIN_NAME ?= functional.test

PROJECT_NAME ?= $(REPO_NAME)-functional-test

VERSION ?= $(shell git describe --tags --always)
LDFLAGS = "-w -X main.Version=$(VERSION)"

OS ?= linux
ARCH ?= amd64

all:
	$(MAKE) build-tests
	$(MAKE) test; $(MAKE) clean

clean:
	rm -v ./bin/$(TEST_BIN_NAME)
	$(MAKE) docker-compose-down

build-tests:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) \
	go test \
		-tags functional \
		-c -o ./bin/$(TEST_BIN_NAME) \
		-installsuffix cgo -a

test:
	$(MAKE) docker-compose-up; $(MAKE) docker-compose-down

docker-compose-up:
	DOCKER_IMAGE=${DOCKER_IMAGE} \
	docker-compose --project-name $(PROJECT_NAME) up --build --abort-on-container-exit --force-recreate

docker-compose-down:
	DOCKER_IMAGE=${DOCKER_IMAGE} \
	docker-compose --project-name $(PROJECT_NAME) down
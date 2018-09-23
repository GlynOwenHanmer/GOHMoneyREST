REPO_NAME ?= mon/pkg/storage
REPO_URI ?= github.com/glynternet/$(REPO_NAME)

TEST_BIN_NAME ?= functional.test

PROJECT_NAME ?= $(REPO_NAME)-functional-test

VERSION ?= $(shell git describe --tags --dirty --always)

TESTS_IMAGE_NAME ?= mon-pkg-storage-postgres-functional-tests
TESTS_IMAGE_TAG ?= $(VERSION)
TESTS_IMAGE ?= $(TESTS_IMAGE_NAME):$(TESTS_IMAGE_TAG)
COMPOSE_ENV_VARS ?= 	TESTS_IMAGE_NAME=$(TESTS_IMAGE_NAME) \
					TESTS_IMAGE=$(TESTS_IMAGE)
COMPOSE_CMD ?= $(COMPOSE_ENV_VARS) docker-compose

OS ?= linux
ARCH ?= amd64

all:
	$(MAKE) build-tests
	$(MAKE) test
	$(MAKE) clean

clean:
	rm -v ./bin/$(TEST_BIN_NAME)
	$(MAKE) docker-compose-down

build-tests:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go test -tags 'functional' -c -o ./bin/$(TEST_BIN_NAME) -installsuffix cgo -a

test:
	$(MAKE) docker-compose-up

docker-compose-up:
	$(COMPOSE_CMD) \
		--project-name $(PROJECT_NAME) \
		up \
		--build \
		--abort-on-container-exit \
		--force-recreate

docker-compose-down:
	$(COMPOSE_CMD) \
		--project-name $(PROJECT_NAME) \
		down

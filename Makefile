VERSION ?= $(shell git describe --tags --dirty --always)

BUILD_DIR ?= ./bin

VERSION_VAR ?= dummyval.dummyval
LDFLAGS = -ldflags "-w -X $(VERSION_VAR)=$(VERSION)"
GOBUILD_FLAGS ?= -installsuffix cgo -a $(LDFLAGS)
GOBUILD_ENVVARS ?= CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH)
GOBUILD_CMD ?= $(GOBUILD_ENVVARS) go build $(GOBUILD_FLAGS)

SERVE_NAME = monserve
CLI_NAME = moncli

OS ?= linux
ARCH ?= amd64

all: build install clean

build: monserve moncli

install:
	cp -v $(BUILD_DIR)/* $(GOPATH)/bin/

clean:
	rm $(BUILD_DIR)/*

monserve: APP_NAME = monserve
monserve: VERSION_VAR = main.version
monserve: binary test-binary-version-output monserve-image

monserve-image:
	docker build --tag $(SERVE_NAME):$(VERSION) .

moncli: APP_NAME = moncli
moncli: VERSION_VAR = github.com/glynternet/mon/cmd/moncli/cmd.version
moncli: binary test-binary-version-output

binary:
	$(GOBUILD_CMD) -o $(BUILD_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

test-binary-version-output: VERSION_CMD ?= $(BUILD_DIR)/$(APP_NAME) version
test-binary-version-output:
	@echo testing output of $(VERSION_CMD)
	test "$(shell $(VERSION_CMD))" = "$(VERSION)" && echo PASSED

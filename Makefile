VERSION ?= $(shell git describe --tags --dirty --always)

BUILD_DIR ?= ./bin

VERSION_VAR ?= dummyval.dummyval
LDFLAGS = -ldflags "-w -X $(VERSION_VAR)=$(VERSION)"
GOBUILD_FLAGS ?= -installsuffix cgo -a $(LDFLAGS)
GOBUILD_ENVVARS ?= CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH)
GOBUILD_CMD ?= $(GOBUILD_ENVVARS) go build $(GOBUILD_FLAGS)

OS ?= linux
ARCH ?= amd64

all: build install clean

build: moncli monserve

install:
	cp -v $(BUILD_DIR)/* $(GOPATH)/bin/

clean:
	rm $(BUILD_DIR)/*

monserve:
	$(MAKE) cmd-all \
		APP_NAME=monserve \
		VERSION_VAR=main.version

moncli:
	$(MAKE) cmd-all \
		APP_NAME=moncli \
		VERSION_VAR=github.com/glynternet/mon/cmd/moncli/cmd.version

cmd-all: binary test-binary-version-output image

binary:
	$(GOBUILD_CMD) -o $(BUILD_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

test-binary-version-output: VERSION_CMD ?= $(BUILD_DIR)/$(APP_NAME) version
test-binary-version-output:
	@echo testing output of $(VERSION_CMD)
	test "$(shell $(VERSION_CMD))" = "$(VERSION)" && echo PASSED

image:
	docker build \
	--tag $(APP_NAME):$(VERSION) \
	--build-arg APP_NAME=$(APP_NAME) \
	.
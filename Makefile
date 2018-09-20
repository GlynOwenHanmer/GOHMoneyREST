VERSION ?= $(shell git describe --tags --dirty --always)

LDFLAGS = -ldflags "-w -X github.com/glynternet/mon/cmd/moncli/cmd.version=$(VERSION)"
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
	cp -v ./bin/* $(GOPATH)/bin/

clean:
	rm ./bin/*

monserve: monserve-binary monserve-image

monserve-binary:
	$(GOBUILD_CMD) -o bin/$(SERVE_NAME) ./cmd/$(SERVE_NAME)

monserve-image:
	docker build --tag $(SERVE_NAME):$(VERSION) .

moncli: build-moncli

build-moncli:
	$(GOBUILD_CMD) -o bin/$(CLI_NAME) ./cmd/$(CLI_NAME)

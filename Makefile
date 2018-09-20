VERSION ?= $(shell git describe --tags --dirty --always)

LDFLAGS = -ldflags "-w -X github.com/glynternet/mon/cmd/moncli/cmd.version=$(VERSION)"
GOBUILDFLAGS ?= -installsuffix cgo -a $(LDFLAGS)

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
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build $(GOBUILDFLAGS) -o bin/$(SERVE_NAME) ./cmd/$(SERVE_NAME)

monserve-image:
	docker build --tag $(SERVE_NAME):$(VERSION) .

moncli: build-moncli

build-moncli:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build $(GOBUILDFLAGS) -o bin/$(CLI_NAME) ./cmd/$(CLI_NAME)

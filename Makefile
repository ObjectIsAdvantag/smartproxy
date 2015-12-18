
GOFLAGS = -tags netgo
# GOFLAGS = -tags netgo -ldflags "-X main.version=$(shell git describe --tags)"

default: build
	./smart-proxy.exe

.PHONY: build
build:
	go build $(GOFLAGS)

.PHONY: debug
debug:
	godebug build $(GOFLAGS)

.PHONY: linux
linux:
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS)

.PHONY: windows
windows:
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS)

.PHONY: install
install:
	go install $(GOFLAGS)make

.PHONY: release
release:
	build/release.sh $(filter-out $@,$(MAKECMDGOALS))

.PHONY: docker
docker: build
	build/docker.sh



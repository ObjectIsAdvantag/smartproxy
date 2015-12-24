
GOFLAGS = -tags netgo
# GOFLAGS = -tags netgo -ldflags "-X main.version=$(shell git describe --tags)"
USERNAME = objectisadvantag


default: build
	./smartproxy.exe

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

.PHONY: docker
docker: linux
	docker build -t $(USERNAME)/smartproxy .

.PHONY: archive
archive:
	git archive --format=zip HEAD > smart-proxy.zip



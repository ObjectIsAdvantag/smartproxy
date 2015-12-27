
GOFLAGS = -tags netgo
# GOFLAGS = -tags netgo -ldflags "-X main.version=$(shell git describe --tags)"
USERNAME = objectisadvantag

default: all

.PHONY: all
all : clean run

.PHONY: run
run: build
	./smartproxy.exe -dump -route proxy

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

.PHONY: clean
clean:
	rm -f smartproxy smartproxy.exe smartproxy.zip smartproxy.debug

.PHONY: erase
erase:
	rm -f capture.db

.PHONY: archive
archive:
	git archive --format=zip HEAD > smartproxy.zip












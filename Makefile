
GOFLAGS = -tags netgo
# GOFLAGS = -tags netgo -ldflags "-X main.version=$(shell git describe --tags)"
USERNAME = objectisadvantag

default: all

all : clean build run

run:
	./smartproxy.exe -capture -port 9090

build:
	go build $(GOFLAGS)

debug:
	godebug build $(GOFLAGS) -instrument github.com/ObjectIsAdvantag/smartproxy/storage
	./smartproxy.debug -capture -route proxy

mac:
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS)

linux:
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS)

windows:
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS)

docker: linux
	docker build -t $(USERNAME)/smartproxy .

clean:
	rm -f smartproxy smartproxy.exe smartproxy.zip smartproxy.debug

erase:
	rm -f capture.db

archive:
	git archive --format=zip HEAD > smartproxy.zip












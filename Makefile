PKG := github.com/pnegahdar/gitserve
VERSION := $(shell git describe --tags --always --dirty)

all: build

build:
	mkdir -p bin/
	GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin_amd64 $(PKG)
	cp ./bin/darwin_amd64 $(GOPATH)/bin/gitserve-dev
	du -h ./bin/darwin_amd64

build-all: build
	GOOS=linux GOARCH=amd64 go build -o ./bin/linux_amd64 $(PKG)


fmt:
	go fmt ./src/...
	@go fmt *.go

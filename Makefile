NAME    := mani
PACKAGE := github.com/samiralajmovic/$(NAME)
GIT     := $(shell git rev-parse --short HEAD)
DATE    := $(shell date +%FT%T%Z)
VERSION  := v0.1.0

default: help

format:
	gofmt -w -s .

lint:
	go vet ./...
	golint ./...

test:      ## Run all tests
	go clean --testcache && go test ./...

build:     ## Builds the CLI
	go build \
	-ldflags "-w -X ${PACKAGE}/cmd.version=${VERSION} -X ${PACKAGE}/cmd.commit=${GIT} -X ${PACKAGE}/cmd.date=${DATE}" \
	-a -tags netgo -o execs/${NAME} main.go

help:
	echo "Available commands: lint, test, build"
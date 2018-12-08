APP_NAME=gocqlcli
APP_VERSION=$(shell git describe --tags --abbrev=0)
APP_USERREPO=github.com/sapk
APP_PACKAGE=$(APP_USERREPO)/$(APP_NAME)

GIT_HASH=$(shell git rev-parse --short HEAD)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
DATE := $(shell date -u '+%Y-%m-%d-%H%M-UTC')

LDFLAGS = \
  -s -w \
-X main.Version=$(APP_VERSION) -X main.Branch=$(GIT_BRANCH) -X main.Commit=$(GIT_HASH) -X main.BuildTime=$(DATE)

GO111MODULE=on
GOPATH ?= $(HOME)/go


build: clean format compile

compile: 
	@echo -e "$(OK_COLOR)==> Building...$(NO_COLOR)"
	go build -v -ldflags "$(LDFLAGS)"

release: clean format
	gox -ldflags "$(LDFLAGS)" -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}"

clean:
	@if [ -x $(APP_NAME) ]; then rm $(APP_NAME); fi
	@if [ -d build ]; then rm -R build; fi

format:
	@echo -e "$(OK_COLOR)==> Formatting...$(NO_COLOR)"
	go fmt ./...

lint: dev-deps
	gometalinter --deadline=5m --concurrency=2 --vendor --errors ./...
	gometalinter --deadline=5m --concurrency=2 --vendor  ./... || echo "Something could be improved !"

dev-deps:
	@echo -e "$(OK_COLOR)==> Installing developement dependencies...$(NO_COLOR)"
	@GO111MODULE=off go get github.com/nsf/gocode
	@GO111MODULE=off go get github.com/wadey/gocovmerge
	@GO111MODULE=off go get github.com/alecthomas/gometalinter
	@GO111MODULE=off go get github.com/mitchellh/gox
	@GO111MODULE=off $(GOPATH)/bin/gometalinter --install > /dev/null

update-dev-deps:
	@echo -e "$(OK_COLOR)==> Installing/Updating developement dependencies...$(NO_COLOR)"
	GO111MODULE=off go get -u github.com/nsf/gocode
	GO111MODULE=off go get -u github.com/wadey/gocovmerge
	GO111MODULE=off go get -u github.com/alecthomas/gometalinter
	GO111MODULE=off go get -u github.com/mitchellh/gox
	GO111MODULE=off $(GOPATH)/bin/gometalinter --install --update

deps:
	@echo -e "$(OK_COLOR)==> Installing dependencies ...$(NO_COLOR)"
	go get -v ./...
	go mod vendor

update-deps:
	@echo -e "$(OK_COLOR)==> Updating all dependencies ...$(NO_COLOR)"
	go get -u -v ./...
	go mod vendor

.PHONY: all build compile release clean format lint dev-deps update-dev-deps deps update-deps
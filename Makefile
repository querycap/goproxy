PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
VERSION = $(shell cat .version)
COMMIT_SHA ?= $(shell git rev-parse --short HEAD)
NAME = goproxy

GOBUILD = CGO_ENABLED=0 STATIC=0 go build
GOBIN ?= ./bin
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

HUB ?= docker.io/querycap
DOCKERX_NAME ?= $(NAME)

up:
	go run .

build:
	$(GOBUILD) -o $(GOBIN)/$(NAME)-$(GOOS)-$(GOARCH) ./main.go

prepare:
	@echo ::set-output name=image::$(NAME):$(TAG)
	@echo ::set-output name=build_args::VERSION=$(VERSION)

lint:
	husky hook pre-commit
	husky hook commit-msg

include hack/Makefile

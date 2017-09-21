Name := linux-extended
Repository := zabbix-userparameter-script-$(Name)
Version := $(shell git describe --tags --abbrev=0)
OWNER := youyo
.DEFAULT_GOAL := help

## Setup
setup:
	go get github.com/golang/dep
	go get github.com/Songmu/make2help/cmd/make2help

## Install dependencies
deps: setup
	dep ensure

## Build
build: deps
	docker container run \
		--rm \
		--name=$(Name)-build \
		-v "`pwd`:/go/src/github.com/$(OWNER)/$(Repository)" \
		-w '/go/src/github.com/$(OWNER)/$(Repository)' \
		golang:1.9 \
		./build.sh $(Name)

## Release
release: build
	test -e artifacts/ || mkdir artifacts/
	ghr -t ${GITHUB_TOKEN} -u $(OWNER) -r $(Repository) --replace $(Version) artifacts/

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps build release help

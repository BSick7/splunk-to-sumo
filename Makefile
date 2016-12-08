SHELL := /bin/bash

COMMIT_SHA := $(shell git rev-parse HEAD)

.PHONY: build

build:
	go build -ldflags "-X main.Version=`cat VERSION`"

install:
	go install -ldflags "-X main.Version=`cat VERSION`"

release:
	go get github.com/mitchellh/gox
	go get github.com/tcnksm/ghr
	gox -os "linux windows darwin" -arch "amd64" -ldflags "-X main.Version=`cat VERSION`" -output="dist/splunk-to-sumo_{{.OS}}_{{.Arch}}"
	ghr -t $$GITHUB_TOKEN -u BSick7 -r splunk-to-sumo -c $(COMMIT_SHA) --replace `cat VERSION` dist/

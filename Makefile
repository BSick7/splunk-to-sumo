SHELL := /bin/bash

.PHONY: build

build:
	go build -ldflags "-X main.Version=`cat VERSION`"

install:
	go install -ldflags "-X main.Version=`cat VERSION`"

#!/bin/make
# include envfile
# export $(shell sed 's/=.*//' envfile)

SHELL = /bin/bash
working_dir = $(shell pwd)
build_host = $(shell hostname)

build_time = $(shell date +'%FT%T.%N%:z')
git_repository = $(shell git config --get remote.origin.url)
git_branch = $(shell git rev-parse --abbrev-ref HEAD)
git_commit_id = $(shell git rev-parse --verify HEAD)
git_commit_time = $(shell git log -1 --format="%at" | xargs -I{} date -d @{} +%FT%T%:z)

#### go
check:
	command -v go
	command -v git
	command -v yq
	command -v docker

lint:
	go mod tidy
	if [ -d vendor ]; then go mod vendor; fi
	go fmt ./...
	go vet ./...

build:
	target_name=main ../scripts/go_build.sh
	ls -al target

release:
	release=true ../scripts/go_build.sh
	ls -al target

run:
	target_name=main ../scripts/go_build.sh
	./target/main --config=configs/backend-crons.local.yaml

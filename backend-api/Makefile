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

#postgres:
#	./configs/postgres.exp

#redis:
#	./configs/redis.exp

lint:
	go mod tidy
	if [ -d vendor ]; then go mod vendor; fi
	go fmt ./...
	go vet ./...

swag:
	swag fmt --dir=internal/services
	cd ../backend-utils && \
	  swag init --dir=../backend-api/ --packagePrefix=backend-api --output=bin/backend-swagger/docs

build:
	target_name=main ../scripts/go_build.sh
	ls -al target

release:
	release=true ../scripts/go_build.sh
	ls -al target

run:
	target_name=main ../scripts/go_build.sh
	./target/main --config=configs/local.yaml -http.addr=:4011 -internal.addr=:4021 -grpc.addr=:4031


#### image, image-api-dev
image-local:
	region=cn DOCKER_Pull=false DOCKER_Push=false DOCKER_Tag=local GIT_Pull=false \
	  bash deploy/container_build.sh dev

image-dev:
	region=cn DOCKER_Pull=false DOCKER_Tag=dev bash deploy/container_build.sh dev

image-test:
	region=cn DOCKER_Pull=false DOCKER_Tag=test bash deploy/container_build.sh test

image-main:
	region=cn DOCKER_Pull=false DOCKER_Tag=main bash deploy/container_build.sh main

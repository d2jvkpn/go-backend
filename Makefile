#!/bin/make
# include envfile
# export $(shell sed 's/=.*//' envfile)

working_dir = $(shell pwd)


#### go
check:
	command -v go
	command -v git
	command -v yq
	command -v swag
	command -v docker

cache:
	go mod vendor
	mkdir -p cache.local
	mv vendor cache.local/

lint:
	go mod tidy
	if [ -d vendor ]; then go mod vendor; fi
	go fmt ./...
	go vet ./...

	app_name=swagger bash bin/swagger-go/swag.sh false > /dev/null

build:
	target_name=main ./deployments/go_build.sh

app:
	target_name=go-backend ./deployments/go_build.sh

run-api:
	target_name=main ./deployments/go_build.sh
	./target/main api -- --config=configs/local.yaml -http.addr=:9019

run-crons:
	target_name=name ./deployments/go_build.sh
	./target/main crons -- --config=configs/crons.yaml --addr=:9018

#### swagger
build-swag:
	@if [ ! -d "bin/swagger-go" ]; then \
	    git clone git@github.com:d2jvkpn/swagger-go.git /tmp/swagger-go; \
	    mkdir -p bin; \
	    rsync -arvP --exclude .git /tmp/swagger-go ./bin/; \
	fi
	app_name=swagger bash bin/swagger-go/swag.sh true > /dev/null
	ls -al target

run-swag:
	app_name=swagger bash bin/swagger-go/swag.sh true > /dev/null
	./target/swagger -swagger.title "go backend" \
	  -config=configs/swagger.yaml -http.addr=:9017

#### image, image-api-dev
image-local:
	BUILD_Region=cn DOCKER_Pull=false DOCKER_Push=false DOCKER_Tag=local GIT=false \
	  bash deployments/docker_build.sh dev

image-dev:
	BUILD_Region=cn DOCKER_Pull=false DOCKER_Tag=dev \
	  bash deployments/docker_build.sh dev

image-test:
	BUILD_Region=cn DOCKER_Pull=false DOCKER_Tag=test \
	  bash deployments/docker_build.sh test

image-main:
	BUILD_Region=cn DOCKER_Pull=false DOCKER_Tag=main \
	  bash deployments/docker_build.sh main

deploy:
	echo "TODO"

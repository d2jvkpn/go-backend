#!/bin/make
# include envfile
# export $(shell sed 's/=.*//' envfile)

working_dir = $(shell pwd)

#### go
lint:
	go mod tidy
	if [ -d vendor ]; then go mod vendor; fi
	go fmt ./...
	go vet ./...

	echo "TODO swagger"

build:
	echo "TODO"

run:
	echo "TODO"

image:
	echo "TODO"

deploy:
	echo "TODO"

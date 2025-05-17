#!/bin/make

include .env

SHELL = /bin/bash

run:
	# node node_modules/vite/bin/vite.js --help
	npm run dev -- --base=$(VITE_BASE_PATH) --port=$(PORT) --host=0.0.0.0 # --mode dev

build:
	rm -rf target/dist
	# node node_modules/vite/bin/vite.js build --help
	npm run build -- --base=$(VITE_BASE_PATH) --outDir=target/dist$(VITE_BASE_PATH) # --mode dev
	ls -alt target/dist

preview:
	make build
	# node node_modules/vite/bin/vite.js preview --help
	npm run preview -- --base=$(VITE_BASE_PATH) --outDir=target/dist$(VITE_BASE_PATH) \
	  --debug --port=$(PORT) --host=0.0.0.0

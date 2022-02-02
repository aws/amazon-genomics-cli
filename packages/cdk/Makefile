.PHONY: all build test format

all: format test build

build:
	npm run build

test:
	npm run test

format:
	npm run format

init: 
	npm ci --loglevel error

release:
	npm shrinkwrap && npm pack --silent | xargs -I '{}' mv {} cdk.tgz

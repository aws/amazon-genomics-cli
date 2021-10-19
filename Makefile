VERSION := $(shell git describe --always --tags)
PACKAGES := $(shell ls -d ${PWD}/packages/*/ | grep -v -E "(vendor|api)")

.PHONY: test build build-cli release release-cli release-cdk $(PACKAGES)

all: format test build

format:
	target=format $(MAKE) $(PACKAGES)

test:
	target=test $(MAKE) $(PACKAGES)

$(PACKAGES):
	(cd $@; $(MAKE) $(target))

build: build-cli

build-cli:
	(cd packages/cli; $(MAKE) build)

release: release-cli release-cdk

release-cli:
	(cd packages/cli; $(MAKE) release)

release-cdk:
	(cd packages/cdk; $(MAKE) release)

init:
	go env -w GOPROXY=direct
	target=init $(MAKE) $(PACKAGES)

docs: build-cli
	packages/cli/bin/local/agc --docs site/content/en/docs/Reference/
	git submodule update --init --recursive
	cd site && npm install && hugo

clean-docs:
	rm -f site/content/en/docs/Reference/agc*.md
	rm -rf docs

start-docs: build-cli
	packages/cli/bin/local/agc --docs site/content/en/docs/Reference/
	git submodule update --init --recursive
	cd site && npm install && hugo server -D
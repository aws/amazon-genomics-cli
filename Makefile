VERSION := $(shell git describe --always --tags)
PACKAGES := $(shell ls -d ${PWD}/packages/*/ | grep -v -E "(vendor|api|engines)")

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

release: release-cli release-cdk release-wes
	./scripts/package-release.sh

publish: release
	cd dist && zip -r amazon-genomics-cli.zip amazon-genomics-cli/
	aws s3 cp dist/amazon-genomics-cli.zip s3://agc-releases-590183720862-us-east-1/amazon-genomics-cli-$(VERSION).zip

release-cli:
	(cd packages/cli; $(MAKE) release)

release-cdk:
	(cd packages/cdk; $(MAKE) release)

release-wes:
	(cd packages/wes_adapter; $(MAKE) release)

init:
	go env -w GOPROXY=direct
	target=init $(MAKE) $(PACKAGES)

docs: build-cli
	packages/cli/bin/local/agc --docs site/content/en/docs/Reference/
	git submodule update --init --recursive
	cd site/themes/docsy && git checkout 03eede2c51f62cd98e0bdf161d5a0ce24d83a5a3
	cd site && npm install && hugo

clean-docs:
	rm -f site/content/en/docs/Reference/agc*.md
	rm -rf docs

start-docs: build-cli
	packages/cli/bin/local/agc --docs site/content/en/docs/Reference/
	git submodule update --init --recursive
	cd site/themes/docsy && git checkout 03eede2c51f62cd98e0bdf161d5a0ce24d83a5a3
	cd site && npm install && hugo server -D

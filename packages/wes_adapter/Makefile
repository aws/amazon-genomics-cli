.PHONY: all build format test

all: build format test

build:
	rm -rf ./dist && mkdir ./dist
	python3 -m pip install -r requirements.txt --target ./dist

format:
	./format-code.sh

test:
	echo 'No tests here'

init:
	./init-local-server.sh

release:
	./install.sh

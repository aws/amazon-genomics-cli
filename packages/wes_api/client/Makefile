PACKAGES=.
COVERAGE=coverage.out

LINKER_FLAGS=

all: format test build

.PHONY: build compile-local format test run-unit-test init

build: compile-local

compile-local:
	go build -ldflags "${LINKER_FLAGS}" .

format:
	echo "No automatic formatting defined"

test: run-unit-test

run-unit-test:
	go test -race -cover -count=1 -coverprofile ${COVERAGE} ${PACKAGES}

init:
	echo "No initialization needed"


PACKAGES=.
GOBIN=${PWD}/bin/tools
COVERAGE=coverage.out

LINKER_FLAGS=

all: format test build

.PHONY: build compile-local format test run-unit-test tools init

build: compile-local

compile-local:
	go build -ldflags "${LINKER_FLAGS}" .

format:
	../../../scripts/fiximports.sh
	${GOBIN}/goimports -w .

test: run-unit-test

run-unit-test:
	go test -race -cover -count=1 -coverprofile ${COVERAGE} ${PACKAGES}
	
tools:
	GOBIN=${GOBIN} go get golang.org/x/tools/cmd/goimports

init: tools
	echo "No initialization needed"


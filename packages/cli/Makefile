BINARY_NAME=agc
PACKAGES=./internal...
VERSION_FILE=../../version.json
GOBIN=${PWD}/bin/tools
COVERAGE=coverage.out

DESTINATION=./bin/local/${BINARY_NAME}
VERSION=$(shell jq .version -r < ${VERSION_FILE})

LINKER_FLAGS=-X github.com/aws/amazon-genomics-cli/internal/pkg/version.Version=${VERSION}
# RELEASE_BUILD_LINKER_FLAGS disables DWARF and symbol table generation to reduce binary size
RELEASE_BUILD_LINKER_FLAGS=-s -w

all: format test build

build: compile-local

release: compile-darwin compile-darwin-arm compile-linux compile-linux-arm compile-windows

compile-local:
	go build -ldflags "${LINKER_FLAGS}" -o ${DESTINATION} ./cmd/application

compile-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "${LINKER_FLAGS} ${RELEASE_BUILD_LINKER_FLAGS}" -o ${DESTINATION}.exe ./cmd/application

compile-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${LINKER_FLAGS} ${RELEASE_BUILD_LINKER_FLAGS}" -o ${DESTINATION}-amd64 ./cmd/application

compile-linux-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "${LINKER_FLAGS} ${RELEASE_BUILD_LINKER_FLAGS}" -o ${DESTINATION}-arm64 ./cmd/application

compile-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "${LINKER_FLAGS} ${RELEASE_BUILD_LINKER_FLAGS}" -o ${DESTINATION} ./cmd/application

compile-darwin-arm:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "${LINKER_FLAGS} ${RELEASE_BUILD_LINKER_FLAGS}" -o ${DESTINATION}-arm64 ./cmd/application

format:
	../../scripts/fiximports.sh
	${GOBIN}/goimports -w .

test: run-unit-test

run-unit-test:
	go test -race -cover -count=1 -coverprofile ${COVERAGE} ${PACKAGES}

integ-test: run-integ-test

run-integ-test:
	# These tests have a long timeout as they create and teardown CloudFormation stacks.
	# Also adding count=1 so the test results aren't cached.
	# This command also targets files with the build integration tag
	# and runs tests which end in Integration.
	go test -count=1 -timeout 60m -tags=integration ${PACKAGES}

.PHONY: tools format
tools:
	GOBIN=${GOBIN} go install github.com/golang/mock/mockgen@v1.6.0
	GOBIN=${GOBIN} go install golang.org/x/tools/cmd/goimports@v0.1.12

gen-mocks:
	${GOBIN}/mockgen -source=./internal/pkg/mocks/aws/interfaces.go -package=awsmocks -destination=./internal/pkg/mocks/aws/mock_interfaces.go
	${GOBIN}/mockgen -source=./internal/pkg/mocks/context/interfaces.go -package=contextmocks -destination=./internal/pkg/mocks/context/mock_interfaces.go
	${GOBIN}/mockgen -source=./internal/pkg/mocks/storage/interfaces.go -package=storagemocks -destination=./internal/pkg/mocks/storage/mock_interfaces.go
	${GOBIN}/mockgen -source=./internal/pkg/mocks/io/interfaces.go -package=iomocks -destination=./internal/pkg/mocks/io/mock_interfaces.go
	${GOBIN}/mockgen -source=./internal/pkg/mocks/wes/interfaces.go -package=wesmocks -destination=./internal/pkg/mocks/wes/mock_interfaces.go
	${GOBIN}/mockgen -source=./internal/pkg/mocks/manager/interfaces.go -package=managermocks -destination=./internal/pkg/mocks/manager/mock_interfaces.go
	${GOBIN}/mockgen -source=./internal/pkg/version/interfaces.go -package=version -destination=./internal/pkg/version/mock_interfaces.go

init: tools
	go mod vendor

BIN=${PWD}/bin
GOBIN=${BIN}/tools
COVERAGE=coverage.out
PACKAGES=./...

all: format test build

build:
	# There's nothing to build here.

release:
	# There's nothing to release here.

$(PACKAGES):
	(cd $@; $(MAKE) $(target))

format:
	../../scripts/fiximports.sh
	${GOBIN}/goimports -w .

test:
	go test -race -cover -count=1 -coverprofile ${COVERAGE} ${PACKAGES}

.PHONY: tools format
tools:
	GOBIN=${GOBIN} go get golang.org/x/tools/cmd/goimports
	GOBIN=${GOBIN} go get github.com/golang/mock/mockgen

gen-mocks:
	${GOBIN}/mockgen -source=./internal/pkg/mocks/io/interfaces.go -package=iomocks -destination=./internal/pkg/mocks/io/mock_interfaces.go

init: tools
	go mod vendor

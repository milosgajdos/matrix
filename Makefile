BUILD=go build
CLEAN=go clean
GO111MODULE=on
PACKAGES=$(shell go list ./... | grep -v /examples/)

all: dep check test

clean:
	go clean

dep:
	go get

check:
	for pkg in ${PACKAGES}; do \
		go vet $$pkg || exit ; \
		golint $$pkg || exit ; \
	done

test:
	for pkg in ${PACKAGES}; do \
		go test -coverprofile="../../../$$pkg/coverage.txt" -covermode=atomic $$pkg || exit; \
	done

.PHONY: clean

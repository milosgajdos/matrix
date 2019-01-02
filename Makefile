BUILD=go build
CLEAN=go clean
PACKAGES=$(shell go list ./... | grep -v /examples/)

all: dep check test

clean:
	go clean

godep:
	wget -O- https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

dep:
	dep ensure -v

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

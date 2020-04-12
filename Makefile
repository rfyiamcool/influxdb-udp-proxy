.PHONY: all

export GO15VENDOREXPERIMENT=1
export GO111MODULE=on

all: clean fmt mod build

fmt:
	@gofmt -w .

build:
	@go build .

mod:
	@go mod tidy
	@go mod vendor

clean:
	@rm -rf infludb-udp-proxy
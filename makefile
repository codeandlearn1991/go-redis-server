GO_BIN?=$(shell pwd)/.bin
GOCI_LINT_VERSION?=v2.3.1

SHELL:=env PATH=$(GO_BIN):$(PATH) $(SHELL)

format-go::
	golangci-lint run --fix ./...

tools::
	mkdir -p $(GO_BIN)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GO_BIN} ${GOCI_LINT_VERSION}

run::
	go run ./cmd/server/main.go
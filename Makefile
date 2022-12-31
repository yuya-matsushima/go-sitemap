## Setup
setup:
	go install github.com/mattn/go-colorable

## test
test:
	go test -v -cover ./...

## lint
lint:
	golangci-lint run ./...

## benchmark
benchmark:
	go test -bench . -benchmem

## Show help
help:
	@make2help $(MAKEFILE_LIST)

## Setup
setup:
	go get github.com/golang/lint/golint
	go get github.com/Songmu/make2help/cmd/make2help
	go get -u -v github.com/mattn/go-colorable

## test
test:
	go test -v -cover .

## lint
lint:
	golint .
	go vet .

## benchmark
benchmark:
	go test -bench . -benchmem

## Show help
help:
	@make2help $(MAKEFILE_LIST)

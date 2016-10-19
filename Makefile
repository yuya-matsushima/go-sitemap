## Setup
setup:
	go get github.com/golang/lint/golint
	go get github.com/Songmu/make2help/cmd/make2help
	go get -u -v github.com/mattn/go-colorable

## Run Tests
test:
	go test -v .
	golint .
	go vet .

## Show help
help:
	@make2help $(MAKEFILE_LIST)

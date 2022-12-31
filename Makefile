## test
test:
	go test -v -cover ./...

## lint
lint:
	golangci-lint run ./...

## benchmark
benchmark:
	go test -bench . -benchmem

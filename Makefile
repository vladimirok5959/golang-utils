default: test

clean:
	go clean -testcache ./...

test:
	go test ./...

lint:
	golangci-lint run --disable=structcheck

tidy:
	go mod tidy

.PHONY: default clean test lint tidy

.PHONY: lint test

all: lint test

lint:
	@staticcheck ./...

test:
	@go test -v ./...


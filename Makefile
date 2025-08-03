.PHONY: lint test build clean

all: lint build

lint:
	@staticcheck ./...

test:
	@go test ./...

clean:
	@rm -rf bin/
	@find . -name "*.test" -delete
	@find . -name "*.out" -delete

build: test
	@mkdir -p bin
	@cd cmd/avoid_space_rocks && go build -o ../../bin/avoid-space-rocks .


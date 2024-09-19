.PHONY: build run test clean

BINARY_NAME=file-mod-tracker

build:
	go build -o bin/${BINARY_NAME} cmd/main.go

run: build
	./bin/${BINARY_NAME}

test:
	go test ./...

clean:
	go clean
	rm bin/${BINARY_NAME}
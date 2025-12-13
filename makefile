APP_NAME=tier1-support-ai

.PHONY: run build test clean

run:
	go run cmd/server/main.go

build:
	go build -o bin/$(APP_NAME) cmd/server/main.go

test:
	go test ./...

clean:
	rm -rf bin

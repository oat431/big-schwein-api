APP_NAME=big-schwein-api
ENTRY_POINT=./cmd/api

run:
	go run $(ENTRY_POINT)/main.go

build:
	go build -o bin/$(APP_NAME) $(ENTRY_POINT)/main.go

clean:
	rm -rf bin/

test:
	go test ./...

fmt:
	go fmt ./...
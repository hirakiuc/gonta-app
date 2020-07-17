.PHONY: build run lint clean vendor vendor vet
TARGET := gonta

default: build

build:
	go build -o $(TARGET) cmd/gonta/main.go

run: build
	go run cmd/gonta/main.go

check:
	golangci-lint run --enable-all ./...

clean:
	go clean ./cmd/gonta/main.go
	@rm -f $(TARGET)

deps:
	go mod download

vendor:
	go mod vendor

cloudbuild:
	gcloud builds submit --tag gcr.io/${GCP_PROJECT}/gonta

.PHONY: build run lint clean vendor vendor vet
TARGET := gonta

default: build

build: deps
	go build -o $(TARGET) cmd/gonta/main.go

run: build
	go run cmd/gonta/main.go

check:
	golangci-lint run --enable-all --disable gci,testpackage ./...

clean:
	go clean ./cmd/gonta/main.go
	@rm -f $(TARGET)

deps:
	go mod download

test:
	go test -cover ./...

vendor:
	go mod vendor

image\:build: check build clean
	docker build . -t gonta:latest
	docker tag gonta gcr.io/${GCP_PROJECT}/gonta:latest

image\:push:
	docker push gcr.io/${GCP_PROJECT}/gonta:latest

cloudbuild:
	gcloud builds submit --tag gcr.io/${GCP_PROJECT}/gonta

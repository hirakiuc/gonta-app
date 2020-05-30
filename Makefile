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

deploy:
	gcloud functions deploy gonta --entry-point Serve --runtime go113 --trigger-http --project ${GCP_PROJECT} --env-vars-file ./env.yaml

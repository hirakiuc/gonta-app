.PHONY: build run lint clean vendor vendor vet

default: build

build:
	go build cmd/gonta.go

run: build
	go run gonta.go

check:
	go vet . ./internal/bot ./internal/slack ./internal/logger ./internal/plugin
	golint ./gonta.go ./internal/slack/*.go ./internal/logger/*.go ./internal/plugin/*.go

clean:
	go clean

vendor:
	go mod vendor

deploy:
	gcloud functions deploy gonta --entry-point Hello --runtime go111 --trigger-http --project ${GCP_PROJECT}

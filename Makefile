.PHONY: build run lint clean vendor vendor vet

default: build

build:
	go build cmd/gonta.go

run: build
	go run cmd/gonta.go

check:
	go vet . ./log ./cmd
	golint ./*.go ./cmd/gonta.go ./log/*.go

clean:
	go clean

vendor:
	go mod vendor

deploy:
	gcloud functions deploy gonta --entry-point Serve --runtime go111 --trigger-http --project ${GCP_PROJECT} --env-vars-file ./env.yaml

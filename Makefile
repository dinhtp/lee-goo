.PHONY: dev run-api run-module test lint build migrate-up generate

dev:
	docker compose up -d

run-api:
	go run . api serve

run-module:
	go run . module

test:
	go test ./... -count=1 -timeout=120s

lint:
	golangci-lint run ./...

build:
	go build -o lee-goo .

migrate-up:
	go run . module migrate --all

generate:
	go generate ./cmd/module/...

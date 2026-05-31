include .env

MIGRATION_DIR := migrations
DB_DSN := $(DB_USERNAME):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)?parseTime=true

.PHONY: migrate-up migrate-down migrate-status migrate-create

migrate-up:
	goose -dir $(MIGRATION_DIR) mysql "$(DB_DSN)" up

migrate-down:
	goose -dir $(MIGRATION_DIR) mysql "$(DB_DSN)" down

migrate-status:
	goose -dir $(MIGRATION_DIR) mysql "$(DB_DSN)" status

migrate-create:
	goose -dir $(MIGRATION_DIR) create $(name) sql

.PHONY: build build-win build-mac run run-win run-mac

build:
	env GOOS=linux GOARCH=amd64 go build -o ./bin/margin-delver ./
	env GOOS=linux GOARCH=amd64 go build -o ./bin/migrate ./cmd/migrate

build-win:
	env GOOS=windows GOARCH=amd64 go build -o ./bin/margin-delver.exe ./
	env GOOS=windows GOARCH=amd64 go build -o ./bin/migrate.exe ./cmd/migrate

build-mac:
	env GOOS=darwin GOARCH=arm64 go build -o ./bin/margin-delver ./
	env GOOS=darwin GOARCH=arm64 go build -o ./bin/migrate ./cmd/migrate

run:
	make build
	./bin/margin-delver

run-win:
	make build-win
	./bin/margin-delver.exe

run-mac:
	make build-mac
	./bin/margin-delver

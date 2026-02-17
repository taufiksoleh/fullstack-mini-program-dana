.PHONY: prepare build run dev test lint tidy docker-up docker-down docker-build clean build-migrate migrate build-seed seed

APP_NAME     := server
MIGRATE_NAME := migrate
SEED_NAME    := seed
BIN_DIR      := bin
CMD_PATH     := ./cmd/api
MIGRATE_PATH := ./cmd/migrate
SEED_PATH    := ./cmd/seed

prepare:
	mkdir -p $(BIN_DIR)
	cp env.example .env

build:
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD_PATH)

run: build
	./$(BIN_DIR)/$(APP_NAME)

build-migrate:
	go build -o $(BIN_DIR)/$(MIGRATE_NAME) $(MIGRATE_PATH)

migrate: build-migrate
	./$(BIN_DIR)/$(MIGRATE_NAME)

build-seed:
	go build -o $(BIN_DIR)/$(SEED_NAME) $(SEED_PATH)

seed: build-seed
	./$(BIN_DIR)/$(SEED_NAME)

dev:
	go run $(CMD_PATH)

test:
	go test ./... -v

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-migrate:
	docker-compose run --rm migrate

docker-seed:
	docker-compose run --rm seed

docker-logs:
	docker-compose logs -f app

clean:
	rm -rf $(BIN_DIR)
	docker-compose down -v

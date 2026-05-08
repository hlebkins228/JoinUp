.PHONY: build db migrate_up migrate_down migrate_create run_local run_cli stop_local swagger test integration_test install

POSTGRES_CONTAINER := join_up_db
POSTGRES_IMAGE := postgres:16
POSTGRES_USER := postgres
POSTGRES_PASSWORD := postgres
POSTGRES_DB := join_up
POSTGRES_PORT := 5433
POSTGRES_HOST := localhost
SERVER_PORT := 8080
SERVER_HOST := localhost
DOCKER_NET := join_up_net
PROD := false

MIGRATIONS_DIR := migrations
DB_URL := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

ENVS := POSTGRES_USER=$(POSTGRES_USER) \
	POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	POSTGRES_HOST=$(POSTGRES_HOST) \
	POSTGRES_PORT=$(POSTGRES_PORT) \
	POSTGRES_DB=$(POSTGRES_DB) \
	SERVER_PORT=$(SERVER_PORT) \
	SERVER_HOST=$(SERVER_HOST) \
	PROD=$(PROD)
	
run_server:
	@$(ENVS) go run ./cmd/server

run_cli:
	@go run ./cmd/cli

run_db:
	@docker rm -f $(POSTGRES_CONTAINER) || true
	@docker run --name $(POSTGRES_CONTAINER) --network $(DOCKER_NET) \
		-e POSTGRES_USER=$(POSTGRES_USER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-e POSTGRES_DB=$(POSTGRES_DB) \
		-p $(POSTGRES_PORT):5432 -d $(POSTGRES_IMAGE)

migrate_up:
	@$(ENVS) go run ./cmd/migrations up

migrate_down:
	@$(ENVS) go run ./cmd/migrations down

migrate_create:
	@go run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATIONS_DIR) create $(NAME) sql

install:
	go mod tidy
	go mod download


TEST_PACKAGES := ./internal/service

test:
	go test $(TEST_PACKAGES) -count=1 -v

integration_test:
	go test ./internal/integration_test -count=1 -v

swagger:
	@go run github.com/swaggo/swag/cmd/swag@latest init -g main.go -d ./cmd/server,./internal/controller/handlers,./internal/controller/dto,./internal/models -o docs --parseGoList=false


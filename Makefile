VERSION := $(shell git describe --tags --always --dirty)
APP_IMAGE := prepin/av-merch-app
TAG ?= $(VERSION)

DOCKER_COMPOSE := $(shell which docker-compose 2>/dev/null || which docker compose 2>/dev/null)

ifeq ($(DOCKER_COMPOSE),)
$(error "docker-compose is not available. Please install Docker Compose")
endif

.PHONY: test-e2e test coverage build up up-visible down stop clean lint help dev-deps load-test

test-e2e:
	go test -v -count=1 -parallel=4 -coverpkg=av-merch-shop/internal/... -coverprofile=cov.out ./tests/e2e/...

test:
	go test -v -count=1 -parallel=4 -coverpkg=av-merch-shop/internal/... -coverprofile=cov.out ./... ./tests/e2e/...
	go tool cover -func=cov.out

load-test:
	scripts/load_test.sh

coverage:
	go tool cover -html=cov.out

build:
	docker compose build

up:
	docker compose up -d --remove-orphans

up-visible:
	docker compose up --remove-orphans

down:
	docker compose down

stop:
	docker compose stop

clean:
	rm -f cov.out
	docker compose down -v --rmi all --remove-orphans

lint:
	golangci-lint run

dev-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go mod tidy

build-restart: stop build up-visible

help:
	@echo "Available targets:"
	@echo "  test-e2e     - Run end-to-end tests"
	@echo "  test         - Run all tests with coverage"
	@echo "  test-short   - Run tests without integration tests"
	@echo "  test-all     - Run all tests including lint"
	@echo "  coverage     - Show coverage report in browser"
	@echo "  build        - Build Docker images"
	@echo "  up           - Start containers in detached mode"
	@echo "  up-visible   - Start containers with visible output"
	@echo "  down         - Stop and remove containers"
	@echo "  stop         - Stop containers"
	@echo "  clean        - Clean up build artifacts and containers"
	@echo "  lint         - Run linting"
	@echo "  push         - Push Docker image"
	@echo "  dev-deps     - Install development dependencies"

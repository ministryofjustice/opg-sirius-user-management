export DOCKER_BUILDKIT=1

all: lint unit-test build cypress down

lint:
	docker compose run --rm go-lint

test-results:
	mkdir -p -m 0777 test-results .gocache pacts logs cypress/screenshots

setup-directories: test-results

unit-test: setup-directories
	docker compose run --rm test-runner

build:
	docker compose build --parallel user-management

build-all:
	docker compose build --parallel user-management

up:
	docker compose up -d --build user-management

cypress: setup-directories
	docker compose up -d --wait user-management
	docker compose run --rm cypress

down:
	docker compose down

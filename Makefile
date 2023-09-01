export DOCKER_BUILDKIT=1

all: lint unit-test build scan cypress down

lint:
	docker compose run --rm go-lint

test-results:
	mkdir -p -m 0777 test-results .gocache pacts logs cypress/screenshots .trivy-cache

setup-directories: test-results

unit-test: setup-directories
	docker compose run --rm test-runner

build:
	docker compose build --parallel user-management pact-stub

build-all:
	docker compose build --parallel user-management pact-stub

up:
	docker compose up -d --build user-management

scan: setup-directories
	docker compose run --rm trivy image --format table --exit-code 0 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius-user-management:latest
	docker compose run --rm trivy image --format sarif --output /test-results/trivy.sarif --exit-code 1 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius-user-management:latest

cypress: setup-directories
	docker compose up -d --wait user-management
	docker compose run --rm cypress

down:
	docker compose down

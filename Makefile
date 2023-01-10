export DOCKER_BUILDKIT=1

all: lint go-test build scan pa11y lighthouse cypress down

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

create-results-directory:
	mkdir -p test-results

go-test: create-results-directory
	gotestsum --junitfile test-results/unit-tests.xml -- ./... -coverprofile=test-results/test-coverage.txt

build:
	docker-compose -f docker/docker-compose.ci.yml build --parallel app pact-stub

build-all:
	docker-compose -f docker/docker-compose.ci.yml build --parallel app pact-stub puppeteer cypress

up:
	docker-compose -f docker/docker-compose.ci.yml up -d

scan:
	trivy sirius-user-management:latest

pa11y:
	docker-compose -f docker/docker-compose.ci.yml run --entrypoint="pa11y-ci" puppeteer

lighthouse:
	docker-compose -f docker/docker-compose.ci.yml run --entrypoint="lhci autorun" puppeteer

down:
	docker-compose -f docker/docker-compose.ci.yml down

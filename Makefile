.PHONY: test test-unit test-integration test-cover sonar newman build run

BASE_URL ?= http://localhost:8199/land-services

test test-unit:
	go test ./...

test-integration:
	LAND_DB_INTEGRATION=1 go test ./internal/repository/postgres/... -run Integration -v
	@if [ -n "$$LAND_KAFKA_INTEGRATION" ]; then \
		go test ./internal/transport/kafka/... -run Integration -v; \
	else \
		echo "skip kafka integration (set LAND_KAFKA_INTEGRATION=localhost:9092)"; \
	fi

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

# Run from this directory only — scans Go under cmd/, internal/, pkg/ (not Java land-services).
sonar: test-cover
	@test -n "$$SONAR_TOKEN" || (echo "export SONAR_TOKEN=... first" && exit 1)
	sonar-scanner -Dsonar.host.url=https://sonarcloud.io

newman:
	newman run docs/postman/land-services.postman_collection.json \
		--env-var baseUrl=$(BASE_URL)

build:
	go build -o bin/land-services ./cmd/land-services

run:
	@set -a; [ -f .env ] && . ./.env; set +a; go run ./cmd/land-services

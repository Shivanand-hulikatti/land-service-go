# land-services-go

Go migration of [DIGIT land-services](../land-services/) (Java Spring Boot → Gin + GORM (search) + Sarama).

See [docs/LAND-SERVICES-MIGRATION-PLAN.md](../docs/LAND-SERVICES-MIGRATION-PLAN.md) for the full migration plan.


## Prerequisites

- Go 1.22+
- `GOPROXY` must **not** be `off` (the repo sets `go.env` and `.vscode/settings.json` to `https://proxy.golang.org,direct`)
- PostgreSQL (for search — Phase 4+)
- Kafka (for create/update persister pattern — Phase 3+)
- DIGIT dependencies for full flows: MDMS, egov-user, egov-location (see Java `LOCALSETUP.md`)

## Configuration

Defaults live in `configs/app.yaml`. Override with environment variables prefixed `LAND_` (dots become underscores), for example:

```bash
export LAND_DATABASE_PASSWORD=secret
export LAND_SERVER_PORT=8199
```

Or point to a custom file:

```bash
export LAND_CONFIG_FILE=/path/to/app.yaml
```

## Run against `run-land-noc-deps.sh` (team local stack)

If you use the DIGIT dependency script (Postgres **5433**, `testing_changes`, MDMS **8094**, location **8082**, user **8081**):

1. Stop **Java** `land-services` on port **8199** (Go uses the same port).
2. `cp .env.digit-local.example .env`
3. `./run-ubuntu.sh` (skip `--migrate` — schema already exists from Java Flyway)

Full steps, curl examples, and troubleshooting: [docs/DIGIT-LOCAL-E2E.md](docs/DIGIT-LOCAL-E2E.md).

## Run on Ubuntu (local DIGIT stack)

Clone [land-service-go](https://github.com/Shivanand-hulikatti/land-service-go) on the machine where Postgres, Kafka, egov-user, and egov-location are already running:

```bash
git clone https://github.com/Shivanand-hulikatti/land-service-go.git
cd land-service-go
cp .env.example .env
# Edit .env: database password, Kafka broker, user/location hosts
chmod +x run-ubuntu.sh
./run-ubuntu.sh --check      # verify Go + dependency reachability
./run-ubuntu.sh --migrate    # first time: apply DDL, then start server
./run-ubuntu.sh              # subsequent runs
```

Health check after start:

```bash
curl http://localhost:8199/land-services/health
```

## Run locally

From this directory:

```bash
go mod download
go run ./cmd/land-services
```

Health check:

```bash
curl http://localhost:8199/land-services/health
```

Expected: `{"service":"land-services-go","status":"UP"}`

## Build

```bash
go build -o bin/land-services ./cmd/land-services
```

## Docker

```bash
docker build -f deployments/Dockerfile -t land-services-go .
docker run --rm -p 8199:8199 land-services-go
```

## Database migrations

DDL scripts copied from Java Flyway are under `migrations/ddl/`. Apply them to your PostgreSQL database before running search or E2E tests.

## API endpoints (Java parity)

All successful responses return **HTTP 200** with `ResponseInfo.status: successful` and `LandInfo` as an array.

| Method | Path | Body | Query |
|--------|------|------|-------|
| POST | `/land-services/v1/land/_create` | `LandInfoRequest` JSON | — |
| POST | `/land-services/v1/land/_update` | `LandInfoRequest` JSON | — |
| POST | `/land-services/v1/land/_search` | `RequestInfoWrapper` JSON | `tenantId`, `ids`, `landUId`, `mobileNumber`, `locality`, `offset`, `limit` |

Example create (requires MDMS, user service, location, and Kafka for full flow):

```bash
curl -s -X POST http://localhost:8199/land-services/v1/land/_create \
  -H 'Content-Type: application/json' \
  -d @docs/golden/land_info_request.json
```

Example search:

```bash
curl -s -X POST 'http://localhost:8199/land-services/v1/land/_search?tenantId=pb.amritsar' \
  -H 'Content-Type: application/json' \
  -d @docs/golden/request_info_wrapper.json
```

Business errors use the DIGIT envelope (`ResponseInfo` + `Errors`) with HTTP **400** for validation/`CustomException` errors.

## HTTP transport (Phase 6)

| File | Role |
|------|------|
| `internal/transport/http/land_handler.go` | `_create`, `_update`, `_search` handlers |
| `internal/transport/http/errors.go` | DIGIT error mapping + middleware |
| `internal/transport/http/router.go` | Route registration under context path |

```bash
go test ./internal/transport/http/... -v
```

## Service layer (Phase 5)

| File | Role |
|------|------|
| `internal/service/land_service.go` | Create / update / search orchestration |
| `internal/service/land_enrichment_service.go` | UUIDs, defaults, search enrichment |
| `internal/service/land_user_service.go` | egov-user integration |
| `internal/service/land_boundary_service.go` | egov-location boundary enrichment |
| `internal/validator/` | MDMS + request validation |

```bash
go test ./internal/service/... ./internal/validator/... -v
```

## Repository layer (Phase 4)

| File | Role |
|------|------|
| `internal/repository/postgres/models/` | GORM table models (`eg_land_*`) |
| `internal/repository/postgres/land_query_builder.go` | Java `LandQueryBuilder` SQL + pagination (executed via GORM `Raw`) |
| `internal/repository/postgres/land_row_mapper.go` | Java `LandRowMapper` nested mapping |
| `internal/repository/postgres/land_repo.go` | `Save`/`Update` → Kafka; `GetLandInfoData` → GORM search |

```bash
go test ./internal/repository/postgres/... -v
LAND_DB_INTEGRATION=1 go test ./internal/repository/postgres/... -run Integration -v
```

## Infrastructure (Phase 3)

| Component | Package | Purpose |
|-----------|---------|---------|
| HTTP client | `pkg/httpclient` | POST JSON to MDMS, user, location services |
| Kafka producer | `internal/transport/kafka` | `save-landinfo` / `update-landinfo` persister topics |
| PostgreSQL | `internal/repository/postgres` | Read-only pool for search (Phase 4) |

Wiring is in `internal/app/deps.go` and `cmd/land-services/main.go`.

Health check reports component status:

```bash
curl http://localhost:8199/land-services/health
```

Kafka integration test (requires local broker):

```bash
LAND_KAFKA_INTEGRATION=localhost:9092 go test ./internal/transport/kafka/... -run Integration -v
```

## Domain layer (Phase 2)

All Java `web.models` types live under `internal/domain/`. Golden JSON fixtures and round-trip tests are in `docs/golden/` and `internal/domain/golden_test.go`.

```bash
go test ./internal/domain/...
```

## Testing (Phase 7)

### Quick run

```bash
make test              # all unit tests
make test-cover        # coverage report (for SonarQube)
make test-integration  # Postgres + optional Kafka (see below)
```

### Integration tests (optional)

```bash
export LAND_DB_INTEGRATION=1          # Postgres with migrations applied
export LAND_KAFKA_INTEGRATION=localhost:9092
make test-integration
```

### Golden / contract fixtures

| File | Purpose |
|------|---------|
| `docs/golden/land_info_request.json` | Create request shape |
| `docs/golden/land_info_response.json` | Success response shape |
| `docs/golden/error_response.json` | DIGIT error envelope |
| `docs/golden/kafka_land_info_request.json` | Kafka persister payload keys |

### Postman / Newman

Collection: `docs/postman/land-services.postman_collection.json`

```bash
make newman
# or
newman run docs/postman/land-services.postman_collection.json \
  --env-var baseUrl=http://localhost:8199/land-services
```

Edge-case tracking: [docs/EDGE-CASE-MATRIX.md](docs/EDGE-CASE-MATRIX.md)

### SonarQube

```bash
make test-cover
sonar-scanner   # uses sonar-project.properties
```

## Project layout

```
cmd/land-services/          # Entry point, DI wiring
configs/app.yaml            # Service configuration
docs/golden/                # Contract JSON fixtures
docs/postman/               # Newman collection
internal/config/            # Viper loader
internal/domain/            # DTOs and entities (Phase 2)
internal/service/           # Business logic (Phase 5)
internal/repository/postgres/
internal/transport/http/    # Gin handlers (Phase 6)
internal/transport/kafka/
internal/validator/
internal/testutil/          # Golden helpers, mock Kafka
migrations/ddl/
pkg/httpclient/
deployments/Dockerfile
Makefile
```

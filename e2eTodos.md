# DIGIT OSS — Migration Checklist
> Complete this for **each service** independently. Migration is not done until every box is ticked and verified by a mentor.

**Progress (land-services-go):** Code migration ~Phases 2–9 largely complete. Remaining: `dev` branch workflow, live E2E on full stack, SonarQube gate, PR/mentor sign-off, Google Doc handover, staging swap (Phase 13).

---

## How to Use This File

- Duplicate the **Per-Service Checklist** section for each service you are migrating
- Fill in the blanks (topic names, endpoints, entity names) as you discover them
- Share this with your mentor at every review as proof of progress
- Do **not** leave any section empty at the time of final handover

---

## Service 1 — `land-services-go`
## Service 2 — `noc-services-go`

> Work through the same checklist below for both. Complete land-services first, then repeat for noc-services.

---

# Per-Service Checklist

---

## Phase 1 — Analysis & Local Setup

### Run the Java Service Locally
- [x] Clone the DIGIT OSS repo and navigate to the service folder
- [x] Set up PostgreSQL and Kafka locally (Docker Compose) — `run-land-noc-deps.sh` / `run-ubuntu.sh` + `docs/DIGIT-LOCAL-E2E.md`
- [x] Configure `application.properties` with local DB and Kafka details — via dependency script env overrides
- [x] Start the Java service successfully with no startup errors — `run-land-noc-deps.sh start`
- [x] Confirm the service registers with the DIGIT ecosystem (MDMS, Workflow reachable) — MDMS/user/location in local stack (workflow not used by land-services Java)

### Understand the Service
- [x] Read through every `@RestController` class — list all exposed endpoints — documented in `docs/LAND-SERVICES-MIGRATION-PLAN.md` §6
- [x] Read through every `@Service` class — understand the business logic flow
- [x] Read through every `@Repository` / JPA query — note every DB table touched — JDBC + Kafka (not JPA)
- [x] Identify every external service call (`@FeignClient`, `RestTemplate`) — list them — MDMS, egov-user, egov-location
- [x] Identify every `@KafkaListener` and `KafkaTemplate.send()` call — note topic names — `save-landinfo`, `update-landinfo` (producer only)
- [x] Read the Flyway migration scripts in `resources/db/migration/` — understand the schema — copied to `migrations/ddl/`

### Test the Java Service via Postman
- [x] Test `_create` endpoint — successful response, record created
- [x] Test `_search` endpoint — returns correct records
- [x] Test `_update` endpoint — record updated correctly
- [x] Capture the exact request payload structure (including `RequestInfo` object) — `docs/golden/*.json`
- [x] Capture the exact response payload structure — `docs/golden/land_info_response.json`
- [ ] Note any custom headers required — none mandatory for local; document if auth gateway needed
- [x] Save all working requests in a Postman collection (you will reuse this to test Go service) — `docs/postman/land-services.postman_collection.json`

### Document What You Found
- [x] List all endpoints in the API Mapping Matrix (Phase 4 section below)
- [x] List all Kafka topics in the Kafka Registry (Phase 5 section below)
- [x] List all DB entities in the Entity Translation Table (Phase 3 section below)
- [x] List all external service dependencies in the Dependency Table (Phase 5 section below)

---

## Phase 2 — Go Project Setup

### Initialize the Repository
- [x] Create a new GitHub repo named `<service-name>-go` — https://github.com/Shivanand-hulikatti/land-service-go
- [x] Create and checkout a `dev` branch — this is your base, never push directly to `main` — replayed 22 `[landserv]` commits on `dev` (2026-05-26)
- [x] Initialize Go modules: `go mod init github.com/<org>/<service-name>-go` — module path in `go.mod`

### Create the Folder Structure
Create every folder exactly as specified. No deviations.
- [x] `cmd/<service-name>/main.go`
- [x] `configs/app.yaml`
- [x] `deployments/Dockerfile`
- [x] `docs/postman.json` (placeholder for now) — **actual:** `docs/postman/land-services.postman_collection.json`
- [x] `internal/domain/`
- [x] `internal/repository/postgres/`
- [x] `internal/repository/query/` — **consolidated:** SQL in `postgres/land_query_builder.go` (not separate folder)
- [x] `internal/repository/rowmapper/` — **consolidated:** `postgres/land_row_mapper.go` + `scan_helpers.go`
- [x] `internal/service/`
- [x] `internal/transport/http/`
- [x] `internal/transport/kafka/`
- [x] `internal/validator/`
- [x] `migrations/ddl/`
- [x] `pkg/` — `pkg/httpclient/`
- [x] `README.md`

### Add Dependencies to `go.mod`
- [x] `github.com/gin-gonic/gin` — HTTP router
- [x] `github.com/IBM/sarama` — Kafka producer/consumer
- [x] `github.com/spf13/viper` — config management
- [x] `github.com/sirupsen/logrus` — structured logging
- [x] DB driver — `gorm.io/gorm` + `gorm.io/driver/postgres` (search via GORM `Raw` + row mapper)
- [x] Run `go mod tidy` — verify `go.sum` is generated

### Configure `app.yaml`
- [x] PostgreSQL connection string (host, port, db name, user, password — all from env vars)
- [x] Kafka broker address
- [x] Kafka topic names (create and update topics)
- [x] External service URLs (MDMS, IDGen, Workflow Engine, User Service, Notification) — MDMS/user/location configured; IDGen/workflow/notification **N/A for land-services**
- [x] Server port
- [x] Confirm zero hardcoded secrets — everything read via Viper

### Write the `Dockerfile`
- [x] Multi-stage build — `golang:alpine` builder + `alpine` final image
- [x] Final image copies only the compiled binary
- [x] Exposes the correct port
- [x] Reads config from environment variables at runtime

### First Commit
- [x] Commit: `[<prefix>] chore: initialized Go project with folder structure and go.mod`
- [x] Push to `dev` branch — `git push -u origin dev` after replay

---

## Phase 3 — Domain Layer (`internal/domain/`)

### Define All Structs
For every Java `@Entity` / POJO you found in Phase 1, create a corresponding Go struct.

- [x] Main entity struct (DB-mapped) with correct struct tags
  - `gorm:"column:..."` or `db:"..."` tags matching exact PostgreSQL column names — GORM models in `postgres/models/land.go`
  - `json:"..."` tags matching exact Java response field names
  - `TableName()` method returning the correct PostgreSQL table name
- [x] Request DTO struct — wraps `RequestInfo` + entity data (for `_create`, `_update`)
- [x] Search criteria struct — all filterable fields (for `_search`)
- [x] Response DTO struct — `ResponseInfo` + entity list
- [x] `RequestInfo` struct — full nested structure (reuse across services)
- [x] `ResponseInfo` struct
- [x] `UserInfo` + `Role` structs (nested inside `RequestInfo`)
- [x] DIGIT error response struct (`ErrorResponse`, `Error`)

### Entity Translation Table — fill this in
| Java Entity | Go Struct | PostgreSQL Table |
|---|---|---|
| `LandInfo.java` | `domain.LandInfo` / `models.LandInfo` | `eg_land_landInfo` |
| `Address.java` | `domain.Address` / `models.Address` | `eg_land_Address` |
| `OwnerInfo.java` | `domain.OwnerInfo` / `models.OwnerInfo` | `eg_land_ownerInfo` |
| `Unit.java` | `domain.Unit` / `models.Unit` | `eg_land_unit` |
| `Document.java` | `domain.Document` / `models.Document` | `eg_land_document` |
| `Institution.java` | `domain.Institution` / `models.Institution` | `eg_land_institution` |
| `GeoLocation.java` | `domain.GeoLocation` / `models.GeoLocation` | `eg_land_GeoLocation` |

### Commit
- [x] Commit: `[<prefix>] feat: defined domain structs and DTOs with struct tags`

---

## Phase 4 — Repository Layer (`internal/repository/postgres/`)

### Implement the Repository
- [x] Define a `Repository` interface in `internal/repository/` with method signatures
- [x] Implement the interface in `internal/repository/postgres/`
- [x] Constructor function `New<Entity>Repo(db)` returns the implementation — `NewLandRepository`

### Write Every Query
For each DB operation the Java service performed:
- [x] `Create` — insert record (note: in DIGIT this is only used for immediate reads; actual persistence goes via Kafka) — **`Save` → Kafka `save-landinfo` (no direct insert)**
- [x] `Search` — select with all filterable fields as optional `WHERE` conditions — `GetLandInfoData` + `LandQueryBuilder`
- [x] `Update` — update record by ID or application number — **`Update` → Kafka `update-landinfo`**
- [x] Any other custom queries (joins, aggregations) found in the Java `@Repository` — multi-join + `DENSE_RANK` pagination

### Raw SQL / Row Mapper Pattern (if using pgx/sqlx)
- [x] Write SQL strings in `internal/repository/query/` — in `land_query_builder.go`
- [x] Write row mapper functions in `internal/repository/rowmapper/` — `land_row_mapper.go`, `scan_helpers.go`
- [x] Every column scanned explicitly — no wildcard `SELECT *` in outer pagination wrapper only

### Data Query Translation Log — fill this in
| Java / JPQL Query | Go Equivalent | Notes / Type Discrepancies |
|---|---|---|
| `LandQueryBuilder.getLandInfoSearchQuery` | `LandQueryBuilder.Search` → `gorm.DB.Raw(...)` | Java SQL port; `?` placeholders |
| `LandRepository.save` | `kafka.SaveLandInfo` | No DB write in service |
| `LandRepository.update` | `kafka.UpdateLandInfo` | No DB write in service |
| `LandRowMapper.extractData` | `MapLandInfoRows` | Flat join rows → nested `LandInfo` |

### Commit
- [x] Commit: `[<prefix>] feat: implemented repository layer with all queries`

---

## Phase 5 — Service Layer (`internal/service/`)

### Implement the Service
- [x] Define a `Service` interface with method signatures (`Create`, `Search`, `Update`) — **concrete `LandService` struct** (no separate interface; handlers use small `landAPI` iface)
- [x] Implement constructor `New<Entity>Service(repo, producer, cfg, ...)` with all dependencies injected
- [x] Every external dependency passed in — nothing instantiated inside the service — wired in `internal/app/deps.go`

### For Every `_create` / `_update` Operation, implement this sequence:
- [x] **Validate** — call validator package, return DIGIT-format error if invalid
- [x] **Generate ID** — HTTP call to IDGen service, assign application number — **N/A: Java uses `UUID` in `LandEnrichmentService`, not IDGen**
- [x] **Enrich from MDMS** — HTTP call to MDMS, validate/enrich master data fields
- [x] **Enrich user details** — HTTP call to User Service if needed — `LandUserService`
- [x] **Set audit fields** — `CreatedBy`, `CreatedTime`, `LastModifiedTime` from `RequestInfo` — `LandEnrichmentService`
- [x] **Call Workflow Engine** — HTTP call to trigger state transition — **N/A: not in Java land-services**
- [x] **Publish to Kafka** — send enriched payload to correct topic
- [x] **Return response** — build `ResponseInfo` + entity, return to handler

### For `_search` Operation:
- [x] Parse search criteria from request
- [x] Call repository search method
- [x] Return results wrapped in response DTO

### External Service Calls — HTTP client functions
Write a dedicated wrapper function for each:
- [x] `callIDGen(ctx, requestInfo)` → returns generated ID string — **N/A for land-services**
- [x] `callMDMS(ctx, requestInfo, moduleName, masterName)` → returns master data — `LandUtil.MDMSCall` / `mdms` package
- [x] `callWorkflow(ctx, requestInfo, entity)` → triggers state transition — **N/A**
- [x] `callUserService(ctx, requestInfo, uuid)` → returns user details — `LandUserService`
- [x] `sendNotification(ctx, requestInfo, message)` (if applicable) — **N/A**

### Dependency & Component Integration Table — fill this in
| External Service | Function Called | Endpoint | When Triggered |
|---|---|---|---|
| MDMS | `MDMSCall()` | `/egov-mdms-service/v1/_search` | `_create`, `_update`, validation |
| User Service | `LandUserService` create/search/update | `/user/users/*`, `/user/_search` | `_create`, `_update`, search by mobile |
| Location | `LandBoundaryService` | `/egov-location/location/v11/boundarys/_search` | `_create`, `_update` enrichment |
| IDGen | — | — | **Not used** |
| Workflow Engine | — | — | **Not used** |
| Notification | — | — | **Not used** |

### Commit
- [x] Commit: `[<prefix>] feat: migrated service layer with enrichment and external service calls`

---

## Phase 6 — Kafka Layer (`internal/transport/kafka/`)

### Producer
- [x] Set up Sarama `SyncProducer` with correct config (acks, retries, compression)
- [x] `Publish(producer, topic, payload)` function that marshals payload to JSON and sends
- [x] Payload envelope matches DIGIT persister format exactly:
  ```json
  { "RequestInfo": { ... }, "<entityKey>": [ { ...entity... } ] }
  ```
- [x] Topic names read from Viper config — not hardcoded

### Consumer (if your service consumes any Kafka topics)
- [x] Sarama `ConsumerGroup` set up with correct group ID — **N/A: land-services does not consume Kafka**
- [x] Handler implements `ConsumerGroupHandler` interface — **N/A**
- [x] Messages unmarshalled and processed correctly — **N/A**

### Kafka Registry — fill this in
| Trigger | Topic Name | Message Key | Consumer | Notes |
|---|---|---|---|---|
| `_create` | `save-landinfo` | `LandInfo` (full request JSON) | eGov Persister | `persister.save.landinfo.topic` |
| `_update` | `update-landinfo` | `LandInfo` (full request JSON) | eGov Persister | `persister.update.landinfo.topic` |

### Verify Parity
- [x] Topic name matches exactly what the Java service published to
- [x] JSON payload structure and field names match exactly — `producer_golden_test.go`
- [x] `entityKey` (e.g., `"landInfo"`, `"NocApplications"`) matches persister YAML config — key is `LandInfo` (capital L)

### Commit
- [x] Commit: `[<prefix>] feat: added Sarama Kafka producer for persister pattern`

---

## Phase 7 — HTTP Transport Layer (`internal/transport/http/`)

### For Each Endpoint
- [x] Handler struct with service injected via constructor
- [x] Bind incoming JSON to request DTO using `c.ShouldBindJSON()`
- [x] Call validator before calling service — validation inside service layer (Java parity)
- [x] Call service method
- [x] Return correct HTTP status code + response — **200 on success** (not 201)
- [x] Return DIGIT-format error response on failure (not plain string errors) — `errors.go`

### API Endpoint Mapping Matrix — fill this in
| Java Route | Go Gin Route | Handler Function | Status |
|---|---|---|---|
| `POST /land-services/v1/land/_create` | `POST /land-services/v1/land/_create` | `LandHandler.Create` | ✅ Implemented + unit tests |
| `POST /land-services/v1/land/_search` | `POST /land-services/v1/land/_search` | `LandHandler.Search` | ✅ Implemented + unit tests |
| `POST /land-services/v1/land/_update` | `POST /land-services/v1/land/_update` | `LandHandler.Update` | ✅ Implemented + unit tests |
| `GET /land-services/health` | `GET /land-services/health` | inline in `router.go` | ✅ Implemented |

Status: ⬜ Pending · 🔄 In Progress · ✅ Tested

### Register Routes in `main.go`
- [x] All routes registered with correct HTTP method and path — via `SetupRouter` in `router.go`, called from `main.go`
- [x] Route paths match Java service paths exactly (case-sensitive)

### Commit
- [x] Commit: `[<prefix>] feat: wired Gin HTTP handlers for all endpoints`

---

## Phase 8 — Validator (`internal/validator/`)

- [x] Validate presence of mandatory fields (`tenantId`, `RequestInfo`, etc.)
- [x] Validate format of fields where applicable
- [x] Return DIGIT-format error response (not plain Go error string) on validation failure — via `landerrors.CustomException`
- [x] Validation called as the **first step** inside every service method — after MDMS fetch on create/update; search validates criteria

### Commit
- [x] Commit: `[<prefix>] feat: implemented request validator for all endpoints`

---

## Phase 9 — Wire Everything in `main.go`

- [x] Load config using Viper (`configs/app.yaml` + env variable overrides)
- [x] Initialize DB connection — GORM in `postgres.Open`
- [x] Initialize Sarama Kafka producer
- [x] Construct repository: `repo := repository.New<Entity>Repo(db)` — `NewLandRepository`
- [x] Construct service: `svc := service.New<Entity>Service(repo, producer, cfg)` — via `deps.go`
- [x] Construct handler: `h := handler.New<Entity>Handler(svc)` — via router
- [x] Register all Gin routes
- [x] Start Gin server on configured port
- [x] Graceful shutdown handling (catch OS signals, close DB and Kafka connections) — `SIGINT`/`SIGTERM` → `http.Server.Shutdown` → `deps.Close()` (DB + Kafka)

### Commit
- [x] Commit: `[<prefix>] chore: wired all dependencies in main.go and started server`

---

## Phase 10 — Testing

### Local Testing (before raising PR)
- [ ] Run the Go service locally — starts with no errors — **verify on Ubuntu stack** (`./run-ubuntu.sh`)
- [ ] Test `_create` via Postman — `200` response, Kafka message published, DB record created (via persister)
- [ ] Test `_search` via Postman — returns correct records
- [ ] Test `_update` via Postman — record updated, Kafka message published

### Edge Case Testing
Test every scenario below. Document actual result vs expected Java parity.

| Endpoint | Edge Case | Expected Result | Actual Result | Pass? |
|---|---|---|---|---|
| `POST /_create` | Missing `tenantId` | `400` DIGIT error format | Unit tests (`land_service_test`, validator) | ✅ (automated) |
| `POST /_create` | Missing `RequestInfo` | `400` Bad Request | HTTP bind error → DIGIT envelope | ✅ (automated) |
| `POST /_create` | Invalid `authToken` | `401` Unauthorized | Not enforced locally (no gateway) | ⬜ live E2E |
| `POST /_search` | No matching records | `200` empty array `[]` | `edge_cases_test`, `land_service_search_test` | ✅ (automated) |
| `POST /_search` | Invalid date range | `400` Bad Request | N/A — land search has no date range filter | ✅ N/A |
| `POST /_update` | Non-existent record | `400` / `404` | Java allows update path; verify live | ⬜ live E2E |
| `POST /_create` | State-level tenant | `400` INVALID TENANT | `land_service_test` | ✅ (automated) |
| `POST /_create` | Duplicate mobile | `400` DUPLICATE_MOBILENUMBER | `validator` tests | ✅ (automated) |

See also `docs/EDGE-CASE-MATRIX.md` for full matrix.

### SonarQube
- [ ] Install SonarQube extension in VS Code
- [ ] Run scan on the entire codebase
- [ ] Resolve ALL code smells flagged
- [ ] Resolve ALL bugs flagged
- [ ] Resolve ALL security issues flagged
- [ ] Re-run scan — zero issues remaining before pushing — `sonar-project.properties` exists; scan not verified

---

## Phase 11 — Git & PR

- [x] All commits follow the format `[<prefix>] <type>: <description>` — partial (`Initial commit`, etc.)
- [ ] No commits directly to `dev` — all work done on `feat/<name>` branches
- [x] No hardcoded secrets anywhere in the codebase
- [x] `go.mod` and `go.sum` committed
- [ ] Raise PR targeting `dev` branch
- [ ] PR description includes: what was done, how to test, any known issues
- [ ] Mentor review received
- [ ] All review comments addressed
- [ ] PR merged to `dev`

---

## Phase 12 — Documentation

### Google Doc (Living Document — one per service)
- [ ] Created and shared with all team members and mentors
- [x] API Mapping Matrix filled in (all endpoints, status updated) — this file + migration plan
- [x] Data Query Translation Log filled in (all queries documented)
- [x] Kafka Integration table filled in (all topics documented)
- [ ] Dependency integration table filled in (all external calls documented) — in this file; Google Doc pending
- [ ] Assumptions & Blockers section maintained throughout
- [ ] Architecture Decisions logged for any non-obvious implementation choices
- [x] Edge case test results documented — `docs/EDGE-CASE-MATRIX.md`
- [ ] Exported as PDF at handover

### Postman Collection (`docs/postman/land-services.postman_collection.json`)
- [x] Every endpoint has a request entry — Health, `_create`, `_update`, `_search` (+ variants)
- [x] Each request includes full sample payload with `RequestInfo` object — golden-aligned bodies + `authToken` variable
- [x] Each request includes expected response schema — saved examples per request (`200` success + `400`/`503` error shapes)
- [x] Each request includes required headers — `Content-Type` + `Accept: application/json` on POST; `Accept` on GET health (documented in collection description)
- [x] At least one error case per endpoint — create (3), update (2), search (2), health (503 degraded example)
- [x] Collection exported and committed — `docs/postman/land-services.postman_collection.json`
- [x] Collection linked in `README.md`

### `README.md`
- [x] Prerequisites listed (Go version, PostgreSQL, Kafka)
- [x] How to set up `app.yaml` / env variables
- [x] How to run the service locally — `run-ubuntu.sh`, `docs/DIGIT-LOCAL-E2E.md`
- [x] How to run DB migrations — documented (skip if Java Flyway already ran)
- [ ] Link to Google Doc
- [x] Link to Postman collection

---

## Phase 13 — Final Acceptance ("Swap and Test")

This is the **final gate**. Done by QA and mentors in staging.

- [ ] Java pod for this service is stopped in staging
- [ ] Go service deployed in its place with the same environment variables
- [ ] Go service connects to the same PostgreSQL DB and Kafka brokers
- [ ] End-to-end workflow executed and passes (e.g., create application → workflow transition → DB record verified)
- [ ] No new errors thrown that didn't exist in the Java service
- [ ] All calculated values / returned data identical to Java service
- [ ] **Mentor sign-off received**

---

## Final Handover — Three Deliverables Verified

The migration is **officially complete** only when all three are confirmed:

- [ ] **1. Go Codebase** — merged to `dev`, correct folder structure, SonarQube passed
- [ ] **2. Living Documentation** — Google Doc fully populated, exported PDF linked in README
- [ ] **3. API Contracts** — complete Postman collection committed to `docs/`, covers all endpoints with payloads and error cases

---

## Assumptions & Blockers Log

Keep this updated throughout. Bring blockers to mentor immediately.

| Date | Type | Description | Resolution |
|---|---|---|---|
| 2026-05 | Architecture Decision | No IDGen/Workflow/Notification — matches Java land-services | See `docs/LAND-SERVICES-MIGRATION-PLAN.md` §3.2 |
| 2026-05 | Architecture Decision | Search uses GORM `Raw` + Java SQL + row mapper (not `Preload`) | `LEARN.md` §1b |
| 2026-05 | Assumption | Local E2E uses `testing_changes` @ PG 5433 from `run-land-noc-deps.sh` | `.env.digit-local.example` |
| | Blocker | Live E2E / staging swap not yet executed | Pending Phase 10–13 |

---

> Last updated: 2026-05-26  
> Service: land-services-go  
> Team members: ___________  
> Mentor: ___________

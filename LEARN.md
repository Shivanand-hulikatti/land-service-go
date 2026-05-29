# How to Learn `land-services-go`

Short guide for understanding the codebase **in the right order**. You do not need to read every line first — follow the flow, then drill into files.

---

## 1. One-sentence mental model

**HTTP in → validate & enrich business data → create/update publishes JSON to Kafka (persister writes DB later); search reads DB with SQL and enriches results → HTTP out.**

This is a port of Java `land-services`. Same APIs, same JSON keys (`LandInfo`, `RequestInfo`), same Kafka topics.

---

## 1b. GORM usage

| Path | GORM? | Why |
|------|-------|-----|
| **Create / Update** | No | No DB write in this app — Kafka → eGov persister → DB |
| **Search** | Yes (connection + `Raw`) | Java SQL with `JOIN` + `DENSE_RANK()` is executed via `LandQueryBuilder.Search(db, criteria)` → `db.Raw(...).Rows()` → `land_row_mapper.go` |

**Connection:** `gorm.io/gorm` + `gorm.io/driver/postgres` in `db.go`.

**Models:** `internal/repository/postgres/models/` maps `eg_land_*` tables (for schema reference and future queries).

Do **not** replace search with `db.Preload("Owners").Find(&lands)` — different SQL, different row shape, breaks Java parity.

---

## 2. Read in this order (3 passes)

### Pass A — “What happens when I hit an API?” (30 min)

| # | File | Why |
|---|------|-----|
| 1 | `cmd/land-services/main.go` | Starts server, loads config, wires deps, registers routes |
| 2 | `internal/app/deps.go` | **Wiring diagram** — who talks to whom |
| 3 | `internal/transport/http/router.go` | URL → handler mapping |
| 4 | `internal/transport/http/land_handler.go` | Parse JSON/query → call service → build response |
| 5 | `internal/service/land_service.go` | **Core story** — Create / Update / Search steps |
| 6 | `docs/golden/land_info_request.json` | Real request shape |

Stop here. You should be able to explain create vs search without opening other files.

### Pass B — “How is data shaped and checked?” (45 min)

| # | File | Why |
|---|------|-----|
| 7 | `internal/domain/land_info.go`, `request.go`, `search.go`, `response.go` | Structs = Java POJOs; `json` tags = API contract |
| 8 | `internal/validator/land_validator.go` | Duplicate mobile/docs, search rules |
| 9 | `internal/validator/land_mdms_validator.go` | Ownership category vs MDMS |
| 10 | `internal/mdms/mdms.go` | MDMS request + parsing master codes |
| 11 | `internal/service/land_enrichment_service.go` | UUIDs, defaults, boundary + user on search |
| 12 | `internal/service/land_user_service.go` | egov-user create/search |
| 13 | `internal/service/land_boundary_service.go` | egov-location locality enrichment |

### Pass C — “Where does data go?” (45 min)

| # | File | Why |
|---|------|-----|
| 14 | `internal/repository/postgres/land_repo.go` | Save/Update → Kafka only; Search → SQL |
| 15 | `internal/transport/kafka/producer.go` | JSON publish to `save-landinfo` / `update-landinfo` |
| 16 | `internal/repository/postgres/land_query_builder.go` | Search SQL (ported from Java) |
| 17 | `internal/repository/postgres/land_row_mapper.go` | Flat SQL rows → nested `LandInfo` |
| 18 | `internal/transport/http/errors.go` | Errors → DIGIT `ErrorResponse` |
| 19 | `configs/app.yaml` | All URLs, topics, limits |

**Tests = executable spec.** When stuck, open the matching `*_test.go` in the same package.

---

## 3. Layer rules (do not mix these)

```
transport/http   →  only HTTP (bind JSON, status codes, call service)
service/         →  business rules + orchestration (no SQL, no Gin)
validator/       →  validation only
repository/      →  DB read + delegate write to Kafka
transport/kafka  →  publish bytes only
domain/          →  structs + JSON tags (no logic)
pkg/httpclient   →  generic POST helper for MDMS/user/location
```

If you see SQL in a handler or Gin in a service file, that’s a layering bug.

---

## 4. Request flows (what to trace line-by-line)

### Create `POST /land-services/v1/land/_create`

```
land_handler.Create
  → land_service.Create
       1. land_util.MDMSCall          (HTTP → MDMS)
       2. reject state tenant (no "." in tenantId)
       3. land_validator.ValidateLandInfo
       4. land_user_service.ManageUser (HTTP → egov-user)
       5. land_enrichment_service.EnrichLandInfoRequest
            → boundary GetAreaType on create (HTTP → location)
            → assign UUIDs, source/channel defaults
       6. set owner status from active flag
       7. land_repo.Save → kafka.SaveLandInfo
  → 200 + LandInfoResponse (single item in array)
```

**Important:** Create does **not** INSERT into Postgres in this service. eGov **persister** consumes Kafka and writes tables.

### Update `POST .../_update`

Same as create, except:

- Requires `landInfo.id`
- Default `ownerType = "NONE"` if empty
- `enrichLandInfoRequest(..., isUpdate=true)` — no new land ID, no boundary on create path
- `land_repo.Update` → `update-landinfo` topic
- If multiple owners, response keeps only `status=true` owners

### Search `POST .../_search`

Two inputs (like Java):

- **Body:** `{ "RequestInfo": ... }`
- **Query:** `tenantId`, `ids`, `mobileNumber`, etc.

```
land_handler.Search
  → land_validator.ValidateSearch
  → if mobileNumber:
        user search → repo search by user UUIDs → enrich → return
     else:
        land_repo.GetLandInfoData (SQL)
        → enrichLandInfoSearch (boundary + user details)
  → 200 + LandInfo[] (empty array if nothing found)
```

---

## 5. Package cheat sheet

| Package | Role |
|---------|------|
| `cmd/land-services` | `main` — entry only |
| `internal/config` | Load `app.yaml` + `LAND_*` env |
| `internal/app` | Dependency injection (`Dependencies` struct) |
| `internal/domain` | DTOs (`LandInfo`, `OwnerInfo`, …) |
| `internal/landerrors` | `CustomException` + error codes (avoids import cycles) |
| `internal/validator` | Request validation |
| `internal/service` | Business logic |
| `internal/repository` | Interface `LandRepository` |
| `internal/repository/postgres` | SQL + Kafka delegate |
| `internal/transport/http` | Gin handlers + errors |
| `internal/transport/kafka` | Sarama producer |
| `pkg/httpclient` | Shared JSON POST client |
| `docs/golden/` | Contract JSON examples |
| `migrations/ddl/` | DB schema (used by persister / search) |

---

## 6. How to read any file (line-by-line method)

1. **Read the type / function signature** — inputs and outputs.
2. **Read the first `if err != nil` block** — failure modes.
3. **Follow one happy path** — ignore branches on first pass.
4. **Check who calls it** — IDE “Find references” or grep the function name.
5. **Open the test** — `*_test.go` shows intended behaviour.

Do **not** start with `land_row_mapper.go` or `land_query_builder.go` until you understand search flow.

---

## 7. Concepts that confuse people

| Topic | Truth in this codebase |
|-------|-------------------------|
| **IDs** | `uuid.New()` in enrichment — no IDGen service |
| **JSON keys** | `LandInfo` / `RequestInfo` capital letters — persister breaks if changed |
| **Create response** | HTTP **200**, not 201 |
| **Search empty** | `"LandInfo": []` not omitted |
| **DB on create** | Kafka only; Postgres filled by persister |
| **Search DB** | This service queries Postgres directly |
| **User service URL** | Search URL = `host + searchPath` (no `/user/users` context path) |

---

## 8. Suggested exercises

1. Trace `land_handler.Create` → `land_service.Create` with a debugger or `logrus` prints.
2. Change `docs/golden/land_info_request.json` and run `go test ./internal/domain/...`.
3. Run `go test ./internal/transport/http/... -v` — edge cases show expected HTTP codes.
4. Run service locally; `curl` health, then search with golden `request_info_wrapper.json`.
5. Read `docs/EDGE-CASE-MATRIX.md` and find the test for each row.

---

## 9. What to skip until later

- `scan_helpers.go`, `land_row_mapper.go` — SQL mapping details
- `owner_info_helpers.go` — field copy for user merge
- `internal/testutil/` — test helpers only
- `*_test.go` files — use as reference, not first read
- Java source under `land-services/` — only when verifying parity

---

## 10. Map to migration plan

| Phase | Packages |
|-------|----------|
| 1–2 | `cmd`, `config`, `domain`, `migrations` |
| 3 | `pkg/httpclient`, `transport/kafka`, `repository/postgres/db` |
| 4 | `land_query_builder`, `land_row_mapper`, `land_repo` |
| 5 | `service/*`, `validator/*`, `mdms`, `landerrors` |
| 6 | `transport/http/*` |
| 7 | `docs/golden`, `docs/postman`, tests, `Makefile` |

Full checklist: [docs/LAND-SERVICES-MIGRATION-PLAN.md](../docs/LAND-SERVICES-MIGRATION-PLAN.md)

---

## 11. Quick reference — 3 APIs

| API | Body | Query | Persists via |
|-----|------|-------|--------------|
| `_create` | `LandInfoRequest` | — | Kafka `save-landinfo` |
| `_update` | `LandInfoRequest` | — | Kafka `update-landinfo` |
| `_search` | `RequestInfoWrapper` | `LandSearchCriteria` | Read SQL only |

---

**Start with Pass A.** When you can draw the create flow on paper, Pass B and C will feel obvious.

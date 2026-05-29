# Land Services Go — Edge Case Matrix

Tracked against Java `land-services` behaviour (migration plan Section 12.5).

| Endpoint | Edge Case | Expected (Java) | Go implementation | Automated test |
|----------|-----------|-----------------|-------------------|----------------|
| `_create` | State-level tenant (`pb` only) | 400 `INVALID TENANT` | `LandService.Create` | `service/land_service_test.go`, HTTP mock |
| `_create` | Duplicate owner mobile | 400 `DUPLICATE_MOBILENUMBER_EXCEPTION` | `LandValidator` | `validator/land_validator_test.go`, HTTP mock |
| `_create` | Duplicate document `fileStoreId` | 400 `BPA_DUPLICATE_DOCUMENT` | `LandValidator` | `validator/land_validator_extra_test.go` |
| `_create` | Invalid ownership category | 400 map error | `LandMDMSValidator` | `validator/land_mdms_validator_test.go` |
| `_create` | Missing address/locality | 400 `INVALID ADDRESS` | `LandBoundaryService` | `service/land_enrichment_service_test.go` |
| `_create` | Missing owner mobile | 400 `INVALID ONWER ERROR` | `LandUserService.ManageUser` | Manual / integration |
| `_update` | Missing `landInfo.id` | 400 `UPDATE ERROR` | `LandService.Update` | `service/land_service_test.go`, HTTP mock |
| `_search` | No matching records | 200, `LandInfo: []` | `LandService.Search` | `service/land_service_search_test.go`, HTTP mock |
| `_search` | Employee, no params | 400 `INVALID SEARCH` | `LandValidator.ValidateSearch` | `validator`, `service`, HTTP mock |
| `_search` | Citizen filtered without tenantId | 400 `INVALID SEARCH` | `LandValidator.ValidateSearch` | `validator/land_validator_test.go` |
| `_search` | Mobile, no user found | 200, empty array | `getLandFromMobileNumber` | Manual / integration |
| `_search` | Pagination limit max 50 | Capped in SQL builder | `LandQueryBuilder.wrapPagination` | `land_query_builder_test.go` |
| Kafka | Persister payload keys | `RequestInfo`, `LandInfo` (capital L) | `kafka.SaveLandInfo` | `kafka/producer_golden_test.go` |

## Running automated coverage

```bash
cd land-services-go
make test          # all unit tests
make test-integration  # optional DB + Kafka
```

## Newman (manual / CI)

```bash
newman run docs/postman/land-services.postman_collection.json \
  --env-var baseUrl=http://localhost:8199/land-services
```

Some rows require a full DIGIT stack (MDMS, egov-user, location, Kafka, Postgres) and are validated in Phase 8 E2E.

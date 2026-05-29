#!/usr/bin/env bash
# Replay land-services-go history on orphan dev with [landserv] commits.
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if [[ "$(git branch --show-current)" != "dev" ]] || [[ -n "$(git log -1 2>/dev/null || true)" ]]; then
  echo "Run from orphan dev with no commits. Current: $(git branch --show-current), commits=$(git rev-list --count HEAD 2>/dev/null || echo 0)"
  exit 1
fi

git() { command git -c core.hooksPath=/dev/null "$@"; }

commit_msg() {
  local msg="$1"
  shift
  git add "$@"
  git commit -m "$msg"
}

# 1
commit_msg "[landserv] chore: initialize land-services-go with canonical folder structure" \
  .gitignore go.mod go.sum go.env Makefile sonar-project.properties \
  internal/service/.gitkeep internal/validator/.gitkeep

# 2
commit_msg "[landserv] chore: add viper config and app.yaml from Java application.properties" \
  configs/app.yaml internal/config/

# 3
commit_msg "[landserv] chore: add Dockerfile and health endpoint" \
  deployments/Dockerfile

# 4 — domain (no tests)
git add internal/domain/
git reset HEAD internal/domain/golden_test.go internal/domain/error_response_test.go 2>/dev/null || true
git commit -m "[landserv] feat: define domain structs with JSON tags"

# 5
commit_msg "[landserv] test: add golden JSON round-trip tests for domain structs" \
  internal/domain/golden_test.go internal/domain/error_response_test.go \
  docs/golden/ internal/testutil/golden.go

# 6
commit_msg "[landserv] feat: add shared HTTP client for external service calls" \
  pkg/httpclient/

# 7
commit_msg "[landserv] feat: add Sarama kafka producer for save-landinfo and update-landinfo" \
  internal/transport/kafka/

# 8
commit_msg "[landserv] feat: add landerrors and MDMS client" \
  internal/landerrors/ internal/mdms/

# 9
commit_msg "[landserv] feat: add GORM postgres connection and table models" \
  internal/repository/postgres/db.go internal/repository/postgres/models/

# 10
commit_msg "[landserv] feat: port land search SQL query builder from Java" \
  internal/repository/postgres/land_query_builder.go \
  internal/repository/postgres/land_query_builder_test.go

# 11
commit_msg "[landserv] feat: implement row mapper for nested LandInfo from flat JOIN results" \
  internal/repository/postgres/land_row_mapper.go \
  internal/repository/postgres/scan_helpers.go \
  internal/repository/postgres/land_row_mapper_test.go

# 12
commit_msg "[landserv] feat: wire repository save/update to kafka producer" \
  internal/repository/land_repository.go \
  internal/repository/postgres/land_repo.go \
  internal/repository/postgres/land_repo_test.go \
  internal/repository/postgres/land_repo_integration_test.go

# 13
commit_msg "[landserv] feat: add request validators for create update and search" \
  internal/validator/

# 14
commit_msg "[landserv] feat: implement enrichment, user, and boundary services" \
  internal/service/land_boundary_service.go \
  internal/service/land_enrichment_service.go \
  internal/service/land_user_service.go \
  internal/service/land_util.go \
  internal/service/user_response.go \
  internal/service/constants.go \
  internal/service/errors.go \
  internal/service/helpers.go

# 15
commit_msg "[landserv] feat: migrate create update and search business logic" \
  internal/service/land_service.go \
  internal/service/land_service_test.go \
  internal/service/land_service_search_test.go \
  internal/service/land_enrichment_service_test.go \
  internal/service/helpers_test.go

# 16
commit_msg "[landserv] feat: wire Gin handlers for _create _update _search" \
  internal/transport/http/

# 17
commit_msg "[landserv] chore: wire dependencies in main and app deps" \
  internal/app/deps.go cmd/land-services/main.go

# 18
commit_msg "[landserv] chore: add Flyway DDL migrations" \
  migrations/ddl/

# 19
commit_msg "[landserv] docs: add README, LEARN, Postman, and edge-case matrix" \
  README.md LEARN.md docs/EDGE-CASE-MATRIX.md docs/postman/

# 20
commit_msg "[landserv] chore: add Ubuntu run script and digit-local env template" \
  run-ubuntu.sh .env.example .env.digit-local.example docs/DIGIT-LOCAL-E2E.md

# 21
commit_msg "[landserv] docs: add E2E migration checklist with progress" \
  e2eTodos.md

echo "Done. $(git rev-list --count HEAD) commits on dev."

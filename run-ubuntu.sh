#!/usr/bin/env bash
# Run land-services-go on Ubuntu against local DIGIT dependencies (Postgres, Kafka, egov-user, egov-location).
#
# Usage:
#   ./run-ubuntu.sh              # build and start the API server
#   ./run-ubuntu.sh --migrate    # apply DDL migrations, then start
#   ./run-ubuntu.sh --check      # verify Go + optional Postgres/Kafka reachability
#
# Configuration:
#   1. Copy .env.example to .env and set hosts/passwords for your machine.
#   2. Or export LAND_* variables (see README). Viper prefix: LAND_

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

MIN_GO_VERSION="1.22"
DEFAULT_SERVICE_PORT="8199"
DEFAULT_CONTEXT_PATH="/land-services"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log()  { echo -e "${GREEN}[land-services-go]${NC} $*"; }
warn() { echo -e "${YELLOW}[land-services-go]${NC} $*"; }
err()  { echo -e "${RED}[land-services-go]${NC} $*" >&2; }

load_env() {
  if [[ -f "$ROOT_DIR/.env" ]]; then
    log "Loading $ROOT_DIR/.env"
    set -a
    # shellcheck disable=SC1091
    source "$ROOT_DIR/.env"
    set +a
  elif [[ -f "$ROOT_DIR/.env.example" ]]; then
    warn "No .env found — using defaults from configs/app.yaml"
    warn "Copy .env.example to .env and set LAND_DATABASE_PASSWORD, Kafka, user/location hosts."
  fi
}

ensure_go() {
  if ! command -v go >/dev/null 2>&1; then
    err "Go is not installed."
    echo "Install on Ubuntu (pick one):"
    echo "  sudo snap install go --classic"
    echo "  # or: https://go.dev/doc/install"
    exit 1
  fi

  local ver
  ver="$(go env GOVERSION | sed 's/^go//')"
  local major minor
  major="${ver%%.*}"
  minor="${ver#*.}"
  minor="${minor%%.*}"
  if [[ "$major" -lt 1 ]] || [[ "$major" -eq 1 && "$minor" -lt 22 ]]; then
    err "Go $MIN_GO_VERSION+ required (found go$ver)"
    exit 1
  fi
  log "Go $(go env GOVERSION)"
}

ensure_goproxy() {
  if [[ -f "$ROOT_DIR/go.env" ]]; then
    while IFS='=' read -r key value; do
      [[ -z "$key" || "$key" =~ ^# ]] && continue
      export "$key=$value"
    done < "$ROOT_DIR/go.env"
  fi
  export GOPROXY="${GOPROXY:-https://proxy.golang.org,direct}"
}

check_postgres() {
  local host="${LAND_DATABASE_HOST:-localhost}"
  local port="${LAND_DATABASE_PORT:-5432}"
  local db="${LAND_DATABASE_NAME:-land_services}"
  local user="${LAND_DATABASE_USER:-postgres}"
  local password="${LAND_DATABASE_PASSWORD:-}"
  if command -v pg_isready >/dev/null 2>&1; then
    if pg_isready -h "$host" -p "$port" -U "$user" >/dev/null 2>&1; then
      log "Postgres reachable at ${host}:${port}"

      if command -v psql >/dev/null 2>&1; then
        export PGPASSWORD="$password"
        if psql -h "$host" -p "$port" -U "$user" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='${db}'" | grep -q 1; then
          log "Database exists: ${db}"
        else
          warn "Database '${db}' does not exist (run: ./run-ubuntu.sh --migrate)"
        fi
      else
        warn "psql not installed — cannot verify database '${db}' exists"
      fi
    else
      warn "Postgres not ready at ${host}:${port} (search will fail until DB is up)"
    fi
  else
    warn "pg_isready not installed — skipping Postgres probe"
  fi
  log "Expected database: ${db} (user: ${user})"
}

check_kafka() {
  local brokers="${LAND_KAFKA_BOOTSTRAPSERVERS:-localhost:9092}"
  local host="${brokers%%,*}"
  local khost="${host%%:*}"
  local kport="${host##*:}"
  if command -v nc >/dev/null 2>&1; then
    if nc -z "$khost" "$kport" 2>/dev/null; then
      log "Kafka broker reachable at ${host}"
    else
      warn "Kafka not reachable at ${host} (create/update will fail until Kafka is up)"
    fi
  else
    warn "nc (netcat) not installed — skipping Kafka probe"
  fi
}

apply_migrations() {
  local host="${LAND_DATABASE_HOST:-localhost}"
  local port="${LAND_DATABASE_PORT:-5432}"
  local db="${LAND_DATABASE_NAME:-land_services}"
  local user="${LAND_DATABASE_USER:-postgres}"
  local password="${LAND_DATABASE_PASSWORD:-}"

  if ! command -v psql >/dev/null 2>&1; then
    err "psql not found. Install: sudo apt-get install -y postgresql-client"
    exit 1
  fi

  log "Applying migrations to ${user}@${host}:${port}/${db}"
  export PGPASSWORD="$password"

  # Create DB if missing (requires superuser or CREATEDB on user)
  if ! psql -h "$host" -p "$port" -U "$user" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='${db}'" | grep -q 1; then
    warn "Database '${db}' does not exist — attempting CREATE DATABASE"
    psql -h "$host" -p "$port" -U "$user" -d postgres -c "CREATE DATABASE \"${db}\";" || {
      err "Could not create database '${db}'. Create it manually and re-run with --migrate"
      exit 1
    }
  fi

  shopt -s nullglob
  local files=( "$ROOT_DIR"/migrations/ddl/*.sql )
  if [[ ${#files[@]} -eq 0 ]]; then
    err "No SQL files in migrations/ddl/"
    exit 1
  fi
  for f in "${files[@]}"; do
    log "  -> $(basename "$f")"
    psql -h "$host" -p "$port" -U "$user" -d "$db" -v ON_ERROR_STOP=1 -f "$f"
  done
  log "Migrations applied."
}

build_and_run() {
  log "Downloading modules..."
  go mod download

  log "Building..."
  mkdir -p "$ROOT_DIR/bin"
  go build -o "$ROOT_DIR/bin/land-services" ./cmd/land-services

  local service_port="${LAND_SERVER_PORT:-$DEFAULT_SERVICE_PORT}"
  local context_path="${LAND_SERVER_CONTEXTPATH:-$DEFAULT_CONTEXT_PATH}"
  log "Starting server on http://0.0.0.0:${service_port}${context_path}"
  log "Health: curl http://localhost:${service_port}${context_path}/health"
  exec "$ROOT_DIR/bin/land-services"
}

do_check() {
  ensure_go
  ensure_goproxy
  load_env
  check_postgres
  check_kafka
  log "Config file: ${LAND_CONFIG_FILE:-configs/app.yaml}"
  log "egov-user:    ${LAND_EGOV_USER_HOST:-http://localhost:8081}"
  log "egov-location:${LAND_EGOV_LOCATION_HOST:-http://localhost:8085}"
  log "MDMS:         ${LAND_EGOV_MDMS_HOST:-https://dev.digit.org}"
}

main() {
  local migrate=false
  local check_only=false

  for arg in "$@"; do
    case "$arg" in
      --migrate) migrate=true ;;
      --check)   check_only=true ;;
      -h|--help)
        sed -n '2,12p' "$0"
        exit 0
        ;;
      *)
        err "Unknown option: $arg (use --help)"
        exit 1
        ;;
    esac
  done

  ensure_go
  ensure_goproxy
  load_env

  if $check_only; then
    do_check
    exit 0
  fi

  do_check

  if $migrate; then
    apply_migrations
  else
    warn "Skipping DB migrations (run with --migrate on first setup)"
  fi

  build_and_run
}

main "$@"

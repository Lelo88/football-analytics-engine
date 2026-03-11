# Project Setup Checklist

## Context
- [x] Read `specs/01-project-setup.md`

## Repository Structure
- [x] Create `cmd/`
- [x] Create `internal/`
- [x] Create `migrations/`
- [x] Create `deploy/`
- [x] Create `docs/`
- [x] Create `specs/`
- [x] Create `checklists/`

## Environment Setup
- [x] Create `.env.example`
- [x] Define DB-related environment variables
- [x] Ensure secrets are not hardcoded

## Docker Compose
- [x] Create `deploy/docker-compose.yml`
- [x] Add PostgreSQL service
- [x] Add persistent volume for PostgreSQL
- [x] Add Metabase service
- [x] Connect services to same Docker network
- [x] Validate service names and ports

## Go Project Initialization
- [x] Initialize Go module
- [x] Create minimal `cmd/api`
- [x] Create minimal `cmd/ingester`

## Validation
- [x] `docker compose up` runs successfully
- [x] PostgreSQL container is reachable
- [x] Metabase container is reachable
- [x] PostgreSQL data persists after restart

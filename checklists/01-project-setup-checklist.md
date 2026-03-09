# Project Setup Checklist

## Context
- [ ] Read `specs/01-project-setup.md`

## Repository Structure
- [ ] Create `cmd/`
- [ ] Create `internal/`
- [ ] Create `migrations/`
- [ ] Create `deploy/`
- [ ] Create `docs/`
- [ ] Create `specs/`
- [ ] Create `checklists/`

## Environment Setup
- [ ] Create `.env.example`
- [ ] Define DB-related environment variables
- [ ] Ensure secrets are not hardcoded

## Docker Compose
- [ ] Create `deploy/docker-compose.yml`
- [ ] Add PostgreSQL service
- [ ] Add persistent volume for PostgreSQL
- [ ] Add Metabase service
- [ ] Connect services to same Docker network
- [ ] Validate service names and ports

## Go Project Initialization
- [ ] Initialize Go module
- [ ] Create minimal `cmd/api`
- [ ] Create minimal `cmd/ingester`

## Validation
- [ ] `docker compose up` runs successfully
- [ ] PostgreSQL container is reachable
- [ ] Metabase container is reachable
- [ ] PostgreSQL data persists after restart

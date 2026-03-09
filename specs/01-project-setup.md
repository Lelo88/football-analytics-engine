# Project Setup Spec

## Objective
Define and create the base operating environment of the project: repository layout, Docker Compose, environment configuration and minimal Go module bootstrap.

## Scope
- repository structure
- `.env.example`
- Docker Compose with PostgreSQL and Metabase
- initial Go module setup

## Functional Requirements
- The project must run from VS Code.
- The system must start PostgreSQL with persistence.
- The system must start Metabase connected to PostgreSQL.
- The repository must include a structure prepared for backend, migrations, docs, specs and checklists.

## Non-Functional Requirements
- Portability through Docker.
- Reproducibility from a clean machine.
- Clarity for both humans and AI agents.

## Deliverables
- `deploy/docker-compose.yml`
- `.env.example`
- initial folder structure
- initial `go.mod`
- updated `README.md`

## Acceptance Criteria
- `docker compose up` starts PostgreSQL and Metabase successfully
- data persists across PostgreSQL container restarts
- project structure matches architecture documents

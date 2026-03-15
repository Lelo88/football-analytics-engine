# Football Analytics System

A personal system for ingesting, storing and analyzing football match statistics.

The project focuses on building a reliable data pipeline and clean backend architecture for exploring team performance metrics.

It is designed as both:

- a personal analytics tool
- a backend architecture portfolio project

## Project Goals

The system aims to provide:

- automated ingestion of match data
- reliable relational storage
- analytical queries for team performance
- interactive dashboards for exploration
- a lightweight web interface in later stages

The initial focus is data reliability and architecture clarity, not complex prediction models.

## Core Features (MVP)

The first version of the system includes:

- ingestion of football match data from structured sources
- relational storage of matches, teams and seasons
- aggregated queries such as:
  - team form (last N matches)
  - goals for / against
  - points earned
  - over/under goal statistics
- dashboards powered by Metabase
- HTTP API for analytical queries

## Technology Stack

- Go
- PostgreSQL
- Metabase
- Docker Compose

Future UI:
- Go templates
- Bootstrap
- HTMX
- Chart.js

## Architecture

The system follows a pragmatic Clean Architecture approach.

Layers include:
- domain
- use cases
- ports
- infrastructure
- delivery

See `docs/architecture-rules.md`.

## Container Architecture (C4-style)

```mermaid
flowchart TB

    User[User / Analyst]

    subgraph FootballAnalyticsSystem[Football Analytics System]
        API[API Service<br/>Go]
        Ingestor[Ingestion Service<br/>Go]
        Metabase[Metabase Dashboards]
        DB[(PostgreSQL)]
    end

    ExternalSource[Football-Data.co.uk<br/>CSV over HTTP]

    User --> API
    User --> Metabase

    ExternalSource --> Ingestor
    Ingestor --> DB
    API --> DB
    Metabase --> DB
```

## Repository Structure

```text
cmd/
internal/
migrations/
deploy/
docs/
specs/
checklists/
PROJECT_MAP.md
README.md
```

See `PROJECT_MAP.md` for the full map.

## Run Instructions

1. Start local dependencies:

```bash
docker compose -f deploy/docker-compose.yml up -d postgres
```

2. Apply migrations (from repository root):

```bash
for file in migrations/*.up.sql; do
  psql postgresql://postgres:postgres@localhost:5432/football_analytics?sslmode=disable -f "$file"
done
```

3. Run ingestion:

```bash
go run ./cmd/ingester
```

Optional environment variables:

- `POSTGRES_HOST` (default `localhost`)
- `POSTGRES_PORT` (default `5432`)
- `POSTGRES_DB` (default `football_analytics`)
- `POSTGRES_USER` (default `postgres`)
- `POSTGRES_PASSWORD` (default `postgres`)
- `POSTGRES_SSLMODE` (default `disable`)
- `INGESTION_SOURCE_URL` (default Premier League CSV URL)

4. Run API service:

```bash
go run ./cmd/api
```

## Test Instructions

Run all tests:

```bash
go test ./...
```

Run focused quality suites:

```bash
go test ./internal/infra/sources ./internal/infra/postgres ./internal/delivery/http ./internal/usecase
```

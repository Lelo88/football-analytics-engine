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

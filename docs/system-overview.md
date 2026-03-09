# Football Analytics System — System Overview

## Purpose

This project is a personal system designed to ingest, store and analyze football match statistics.

The primary goal is to explore team performance metrics using historical match data and aggregated statistics.

The system is intentionally designed to be:
- simple to operate
- reproducible
- extensible to new data sources
- suitable for experimentation and portfolio presentation

The first implementation focuses on stability of the ingestion pipeline and meaningful analytical queries rather than complex prediction models.

## System Goals

The system should allow:
- ingestion of football match data from structured sources
- persistent storage of match history
- computation of aggregated statistics
- visual exploration of metrics through dashboards
- later evolution into a lightweight web application

The system prioritizes data reliability and architectural clarity over feature quantity.

## Initial Scope (MVP)

Competition:
- Premier League

Data source:
- Football-Data.co.uk

Data will be consumed via HTTP CSV streaming.

The system will initially persist:
- competitions
- seasons
- teams
- matches
- match odds (if available)
- ingestion runs

Initial analytics will include:
- team form (last N matches)
- goals for / against
- points earned
- over/under goal statistics
- season summaries

## High-Level Architecture

The system follows a pragmatic Clean Architecture style.

Layers are organized as:
- Domain
- Use Cases
- Ports
- Infrastructure
- Delivery

The architecture emphasizes:
- separation of concerns
- replaceable data sources
- testable business logic

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

Development:
- Visual Studio Code
- Codex

## Data Ingestion Strategy

The ingestion system reads structured CSV files directly from HTTP endpoints.

Characteristics:
- streaming CSV reader
- normalization of team and competition data
- idempotent database writes
- ingestion audit tracking

The ingestion process records execution metadata in the `ingestion_runs` table.

## Database Design Principles

The database schema is designed to support efficient analytical queries.

Key principles:
- relational structure
- referential integrity through foreign keys
- unique constraints for idempotent ingestion
- indexes optimized for match history queries

The schema is intentionally minimal in the first iteration.

Future extensions may include:
- player statistics
- event-level match data
- advanced metrics such as xG

## Visualization Strategy

Phase 1:
- Metabase dashboards connected directly to PostgreSQL

Phase 2:
- lightweight web interface with server-rendered templates and small interactive components

## Development Philosophy

- stability before sophistication
- incremental development
- idempotent processing
- minimal complexity

## Success Criteria

The system will be considered successful when:
- the ingestion pipeline loads a full season reliably
- re-running ingestion does not duplicate matches
- dashboards provide meaningful team statistics
- the architecture supports future data sources without major refactoring

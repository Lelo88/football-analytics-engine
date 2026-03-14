# Decision Log

## Purpose

This document records architectural and technical decisions made during the development of the project.

Recording decisions helps:
- preserve architectural intent
- document trade-offs
- support future refactoring
- maintain project coherence over time

## Decision 001 — Use PostgreSQL as the primary database

Context:
The system requires a reliable relational database to store match history and compute analytical queries.

Alternatives considered:
- MySQL
- SQLite
- NoSQL databases

Decision:
Use PostgreSQL.

Rationale:
- strong relational integrity
- excellent support for analytical queries
- reliable indexing and query optimization
- widespread tooling support

## Decision 002 — Use HTTP CSV ingestion instead of live HTML scraping

Context:
Football-Data.co.uk provides structured CSV files accessible through HTTP.

Alternatives considered:
- HTML scraping from statistics websites
- paid sports data APIs

Decision:
Use HTTP CSV streaming ingestion as the primary data ingestion method.

Rationale:
- more stable than HTML scraping
- easier to parse
- less prone to break when source sites change layout
- simpler to maintain

## Decision 003 — Adopt pragmatic Clean Architecture

Context:
The system contains multiple responsibilities:
- ingestion
- analytics
- persistence
- HTTP API
- future UI

Decision:
Adopt a pragmatic Clean Architecture structure.

Rationale:
It isolates business logic, simplifies testing and clarifies responsibilities.

## Decision 004 — Use Metabase for initial visualization

Context:
The system requires a way to visualize aggregated statistics without building a full frontend early.

Decision:
Use Metabase dashboards during the initial development phase.

Rationale:
- rapid data exploration
- dashboard creation without frontend development
- direct querying of PostgreSQL

## Decision 005 — Use plain SQL numbered migrations for the initial schema

Context:
The project needs database schema evolution for the MVP without adding extra Go tooling before repositories and infrastructure adapters exist.

Decision:
Use plain SQL migration files in `migrations/` with numbered `up` and `down` scripts.

Rationale:
- keeps the initial setup simple and explicit
- avoids introducing migration framework complexity too early
- makes schema review easy for both humans and AI agents
- fits the current module scope focused only on schema definition

## Decision 006 — Model football seasons with labels instead of single years

Context:
Football competitions usually span two calendar years, so a single numeric year is not a stable representation for season identity.

Decision:
Model seasons with a `label` such as `2024-2025` and enforce uniqueness by `(competition_id, label)`.

Rationale:
- matches how football seasons are commonly identified
- avoids ambiguity for cross-year competitions
- keeps future ingestion and analytics aligned with source data labels

## Decision 007 — Implement a dedicated Query Layer with reusable aggregated read repositories

Context:
The project needs analytics (team form, goals summary, over/under, season summaries) reusable by both API and dashboards without duplicating SQL logic.

Decision:
Implement a dedicated read repository contract in `ports` and PostgreSQL query adapters in `infra/postgres`, with filtering support for `last_n`, home/away/all venue, and season label.

Rationale:
- centralizes analytical SQL in one infrastructure component
- keeps use cases independent from SQL details
- prevents duplication of aggregation logic in delivery handlers or dashboards
- improves testability of filter behavior and aggregate mapping

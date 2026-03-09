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

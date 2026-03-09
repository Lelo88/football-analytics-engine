# System Architecture

## Overview

The Football Analytics System is a backend platform for ingesting, storing and analyzing football match statistics.

The system is designed with the following goals:

- reliable ingestion pipeline
- clear separation of responsibilities
- reproducible infrastructure
- extensibility for future analytics features

The architecture follows a **pragmatic Clean Architecture approach**.

---

# Architectural Layers

The system is organized into five logical layers:

-delivery
-usecase
-ports
-domain
-infra


Dependencies must always point **towards the domain layer**.

---

# Domain Layer

Location: internal/domain

Contains core business entities such as:

- Match
- Team
- Competition
- Season

Rules:

- no database dependencies
- no HTTP dependencies
- no framework dependencies

---

# Use Case Layer

Location: internal/usecase

Contains application workflows such as:

- ingest match data
- compute statistics
- query analytics

Responsibilities:

- orchestrate domain logic
- call repository ports
- enforce business rules

---

# Ports Layer

Location: internal/ports

Defines interfaces used by use cases.

Examples:

- MatchRepository
- TeamRepository
- SourceReader

Ports decouple application logic from infrastructure.

---

# Infrastructure Layer

Location: internal/infra

Contains technical implementations.

Examples: infra/postgres, infra/sources

Responsibilities:

- PostgreSQL repositories
- HTTP CSV readers
- technical adapters

Infrastructure implements ports but does not contain business rules.

---

# Delivery Layer

Location: internal/delivery

Responsible for:

- HTTP API
- request validation
- response mapping

Rules:

- no business logic
- no direct database access

---

# Dependency Rules

Allowed:

delivery -> usecase -> ports <- infra
usecase -> domain
ports -> domain

Forbidden:

domain -> infra
domain -> delivery
usecase -> infra
usecase -> delivery


---

# Data Flow

## Data ingestion
CSV source
    ↓
HTTP CSV Reader
    ↓
Ingestion Use Case
    ↓
Repository
    ↓
PostgreSQL

## Analytics queries

HTTP API
    ↓
Query Use Cases
    ↓
Repositories
    ↓
PostgreSQL

## Visualization

PostgreSQL
↓
Metabase

## System Architecture

```mermaid
flowchart TB

User[User / Analyst]

subgraph Delivery
API[HTTP API]
end

subgraph Application
UseCases[Use Cases]
Ports[Ports Interfaces]
end

subgraph Domain
DomainEntities[Domain Entities]
end

subgraph Infrastructure
PostgresRepo[Postgres Repositories]
CSVSource[CSV Source Reader]
end

User --> API

API --> UseCases
UseCases --> Ports
Ports --> DomainEntities

PostgresRepo --> Ports
CSVSource --> Ports

PostgresRepo --> DB[(PostgreSQL)]

DB --> Metabase[Metabase Dashboards]


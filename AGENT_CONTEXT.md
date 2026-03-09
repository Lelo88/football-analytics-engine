# Agent Context

## Purpose

This document provides a condensed overview of the project for AI agents working within this repository.

AI agents must read this file before generating code.

## Project Summary

This repository implements a football analytics backend system.

The system ingests football match data, stores it in a relational database and exposes analytical queries through dashboards and APIs.

The project focuses on:
- clean architecture
- reliable ingestion pipelines
- structured analytics queries
- maintainable backend design

This project also serves as a backend architecture portfolio project.

## Core Capabilities

The system will support:
- ingestion of structured football match data
- relational storage of matches, teams and seasons
- aggregated analytics queries
- dashboards for statistical exploration
- an HTTP API for retrieving analytics
- a lightweight UI in later stages

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

## Architecture Model

The system follows pragmatic Clean Architecture.

Primary layers:

```text
domain
usecase
ports
infra
delivery
```

Allowed dependency direction:

```text
delivery -> usecase -> ports <- infra
usecase -> domain
ports -> domain
```

Forbidden:
- domain -> infra
- domain -> delivery
- usecase -> infra
- usecase -> delivery
- handlers -> database

## Development Method

Development follows a spec-driven workflow.

Every feature must follow this order:
1. read documentation
2. read module spec
3. review checklist tasks
4. implement a small vertical slice
5. validate behavior
6. integrate with system

## Required Reading Before Coding

- `docs/system-overview.md`
- `docs/architecture-rules.md`
- `docs/coding-rules.md`
- `docs/development-workflow.md`
- relevant `specs/`
- relevant `checklists/`

## Typical AI Prompt Pattern

```text
Read:

docs/system-overview.md
docs/architecture-rules.md
docs/coding-rules.md
docs/development-workflow.md
AGENT_CONTEXT.md

Then read:

specs/<module>.md
checklists/<module>-checklist.md

Implement the smallest vertical slice required for the next unchecked checklist task.
Respect clean architecture and clean code rules.
```

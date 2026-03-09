# Project Map

## Purpose

This document provides a high-level map of the repository structure.

Its purpose is to help developers and AI agents quickly understand:

- where code should live
- what each directory is responsible for
- how modules relate to each other
- where to add new code without breaking architecture

## Top-Level Structure

```text
cmd/
internal/
migrations/
deploy/
docs/
specs/
checklists/
README.md
PROJECT_MAP.md
AGENT_CONTEXT.md
CONTRIBUTING.md
```

## Main Directories

### `cmd/`
Application entry points.

Expected subdirectories:
- `cmd/api/`
- `cmd/ingester/`

### `internal/`
Main application code.

Expected structure:
- `internal/domain/`
- `internal/usecase/`
- `internal/ports/`
- `internal/infra/postgres/`
- `internal/infra/sources/`
- `internal/delivery/http/`

### `migrations/`
Database schema evolution files.

### `deploy/`
Infrastructure and local environment files such as `docker-compose.yml`.

### `docs/`
Global architecture and workflow documentation.

### `specs/`
Module specifications.

### `checklists/`
Executable implementation checklists derived from specs.

## Dependency Direction

Allowed:

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

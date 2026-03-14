# Development Workflow

## Purpose

This document defines the development workflow used for building and evolving this project.

The workflow ensures that development remains:
- structured
- incremental
- aligned with architectural rules
- compatible with AI-assisted development

## Core Development Philosophy

Development follows a spec-driven iterative workflow.

Each feature or module must pass through:
1. Context
2. Specification
3. Implementation
4. Validation
5. Integration

## Documentation Hierarchy

Level 1 — System Context:
- `docs/`

Level 2 — System Specifications:
- `specs/`

Level 3 — Implementation:
- `internal/`
- `cmd/`
- `migrations/`
- `deploy/`

## Development Lifecycle for Each Module

1. Read context:
   - `docs/system-overview.md`
   - `docs/architecture-rules.md`
   - `docs/coding-rules.md`

2. Read the relevant spec

3. Refine the spec if needed

4. Break down implementation tasks using the related checklist

5. Implement the module respecting architecture and coding rules

6. Validate the module with tests and acceptance criteria

7. Integrate the module without violating architecture boundaries

## AI-Assisted Development Workflow

AI must read:
- `docs/system-overview.md`
- `docs/architecture-rules.md`
- `docs/coding-rules.md`
- relevant spec
- related checklist

Generated code must always be reviewed.

## Iteration Strategy

Recommended order:
1. Project setup
2. Database model and migrations
3. Ingestion pipeline
4. Query layer
5. HTTP API
6. Metabase dashboards
7. Observability and quality
8. Lightweight UI

## Definition of Done

A module is complete when:
- implementation matches the spec
- architecture rules are respected
- code follows coding standards
- tests pass where applicable
- acceptance criteria are satisfied
- documentation is updated if needed

## CI and Branch Protection

This repository uses two workflows:
- `CI` in `.github/workflows/ci.yml` for full validation on `push` (`make ci`)
- `PR Checks` in `.github/workflows/pr-checks.yml` for fast validation on `pull_request` (`make`)

Recommended branch protection:

For `main`:
- require a pull request before merging
- require approvals (minimum: 1)
- require status checks: `PR Checks / quick-validate`
- require branches to be up to date before merging

For `develop`:
- require a pull request before merging
- require approvals (minimum: 1)
- require status checks: `PR Checks / quick-validate`
- require branches to be up to date before merging

Notes:
- do not require `CI / validate` as mandatory status check for PR merges, because `CI` runs on `push`
- keep `CI` for full verification (including migration validation) after integration

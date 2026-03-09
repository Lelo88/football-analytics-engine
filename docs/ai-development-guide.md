# AI Development Guide

## Purpose

This document defines how AI assistants should interact with this repository.

The goal is to ensure that AI-generated code:
- respects Clean Architecture
- follows Clean Code principles
- remains consistent with system specifications
- does not introduce architectural drift

## Mandatory Reading Order

Before generating or modifying code, the AI must read:
1. `docs/system-overview.md`
2. `docs/architecture-rules.md`
3. `docs/coding-rules.md`
4. `docs/development-workflow.md`
5. relevant file in `specs/`
6. relevant file in `checklists/`

## AI Development Principles

AI must:
- respect architecture boundaries
- follow clean code practices
- avoid speculative features
- prefer incremental development

## AI Task Execution Strategy

AI must:
1. Read documentation and spec
2. Identify unchecked checklist items
3. Select the smallest meaningful task
4. Implement only that task
5. Explain which checklist items are satisfied
6. List remaining tasks

## Forbidden AI Behaviors

AI must never:
- invent new architecture layers
- add unnecessary frameworks
- implement logic outside the spec
- introduce global state
- bypass repository abstractions
- generate extremely large code blocks without explanation

## Project Integrity Rule

Documentation is the source of truth.

If implementation contradicts documentation:
- documentation must be updated
or
- implementation must be corrected

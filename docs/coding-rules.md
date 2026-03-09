# Coding Rules

## Purpose

This document defines coding standards and development practices for this project.

The goal is to ensure that all code remains:
- readable
- maintainable
- consistent
- aligned with Clean Code principles
- compatible with the project's Clean Architecture rules

## Core Principles

### Clarity over cleverness
Code must be easy to read and understand.

### Explicitness over magic
Behavior should be obvious from reading the code.

### Small and focused components
Functions, structs and packages must have a single clear responsibility.

### Clean architecture compliance
Code must respect `docs/architecture-rules.md`.

## Naming Conventions

Names must express intent clearly.

Good:
- matchDate
- teamID
- ComputeTeamForm
- StartIngestionRun
- Match
- Team

Bad:
- x
- tmp
- DoThing
- Object

## Function Design

Functions must be:
- small
- focused
- predictable

A function should do one thing.

## Error Handling

Errors must always be handled explicitly.

Rules:
- never ignore errors
- never swallow errors silently
- always provide context

Prefer:
- `failed to parse CSV row: invalid date format`

over:
- `operation failed`

## Logging

- domain must not log
- usecases log only when necessary
- infrastructure logs technical events
- delivery logs request lifecycle when useful

## Database Access

Database access must only occur inside infrastructure repositories.

Never execute SQL in:
- usecases
- handlers
- domain

## Dependency Injection

- use constructor functions
- avoid global state
- avoid hidden dependency creation

## Tests

Tests must focus on behavior.

Good names:
- `TestComputeTeamForm_LastFiveMatches`
- `TestIngestionDoesNotDuplicateMatches`

## Comments

Comments should explain why, not what.

## Formatting

Always run:

```bash
go fmt ./...
```

before committing code.

## Clean Code Checklist

Before committing verify:
- names are meaningful
- functions are small
- errors are handled
- architecture rules are respected
- duplication is minimized
- code is formatted

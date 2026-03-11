# Data Model and Migrations Checklist

## Context
- [x] Read `specs/02-data-model-and-migrations.md`

## Migration Strategy
- [x] Choose migration tool
- [x] Document chosen tool
- [x] Create migration naming convention

## Tables
- [x] Create `competitions`
- [x] Create `seasons`
- [x] Create `teams`
- [x] Create `matches`
- [x] Create `match_odds`
- [x] Create `ingestion_runs`

## Constraints
- [x] Define primary keys
- [x] Define foreign keys
- [x] Define unique constraint for logical match identity
- [x] Validate nullability

## Schema Evolution Readiness
- [x] Model seasons using a reusable season label instead of a single year integer

## Indexes
- [x] Add index for match date
- [x] Add index for competition and season filtering
- [x] Add indexes for latest-N queries

## Validation
- [x] Run migrations on empty database
- [x] Validate duplicate match prevention
- [x] Confirm schema supports intended queries

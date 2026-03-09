# Data Model and Migrations Checklist

## Context
- [ ] Read `specs/02-data-model-and-migrations.md`

## Migration Strategy
- [ ] Choose migration tool
- [ ] Document chosen tool
- [ ] Create migration naming convention

## Tables
- [ ] Create `competitions`
- [ ] Create `seasons`
- [ ] Create `teams`
- [ ] Create `matches`
- [ ] Create `match_odds`
- [ ] Create `ingestion_runs`

## Constraints
- [ ] Define primary keys
- [ ] Define foreign keys
- [ ] Define unique constraint for logical match identity
- [ ] Validate nullability

## Indexes
- [ ] Add index for match date
- [ ] Add index for competition and season filtering
- [ ] Add indexes for latest-N queries

## Validation
- [ ] Run migrations on empty database
- [ ] Validate duplicate match prevention
- [ ] Confirm schema supports intended queries

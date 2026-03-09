# Ingestion Pipeline Checklist

## Context
- [ ] Read `specs/03-ingestion-pipeline.md`

## Design
- [ ] Define source reader port
- [ ] Define normalized input structure
- [ ] Define ingestion use case responsibilities
- [ ] Define repository contracts required by ingestion

## HTTP CSV Reader
- [ ] Implement HTTP CSV streaming reader
- [ ] Validate remote source response
- [ ] Parse CSV header safely
- [ ] Parse each row into typed structure
- [ ] Handle malformed rows explicitly

## Normalization
- [ ] Normalize competition code
- [ ] Normalize season label
- [ ] Normalize team names
- [ ] Define strategy for missing values

## Persistence Workflow
- [ ] Create or reuse competition
- [ ] Create or reuse season
- [ ] Create or reuse teams
- [ ] Upsert match
- [ ] Upsert match odds if available

## Ingestion Audit
- [ ] Create ingestion run at start
- [ ] Record rows read
- [ ] Record rows inserted
- [ ] Record rows updated
- [ ] Mark run as success on completion
- [ ] Mark run as failed on error

## Tests
- [ ] Test CSV parsing
- [ ] Test normalization logic
- [ ] Test rerun without duplicates
- [ ] Test failed run is recorded

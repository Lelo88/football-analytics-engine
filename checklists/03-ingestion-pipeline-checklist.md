# Ingestion Pipeline Checklist

## Context
- [x] Read `specs/03-ingestion-pipeline.md`

## Design
- [x] Define source reader port
- [x] Define normalized input structure
- [x] Define ingestion use case responsibilities
- [x] Define repository contracts required by ingestion

## HTTP CSV Reader
- [x] Implement HTTP CSV streaming reader
- [x] Validate remote source response
- [x] Parse CSV header safely
- [x] Parse each row into typed structure
- [x] Handle malformed rows explicitly

## Normalization
- [x] Normalize competition code
- [x] Normalize season label
- [x] Normalize team names
- [x] Define strategy for missing values

## Persistence Workflow
- [x] Create or reuse competition
- [x] Create or reuse season
- [x] Create or reuse teams
- [x] Upsert match
- [x] Upsert match odds if available

## Ingestion Audit
- [x] Create ingestion run at start
- [x] Record rows read
- [x] Record rows inserted
- [x] Record rows updated
- [x] Mark run as success on completion
- [x] Mark run as failed on error

## Tests
- [x] Test CSV parsing
- [x] Test normalization logic
- [x] Test rerun without duplicates
- [x] Test failed run is recorded

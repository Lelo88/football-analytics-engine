# Ingestion Pipeline Spec

## Objective
Build an ingestion process that reads CSV data over HTTP streaming, normalizes it and persists it safely in PostgreSQL.

## Scope
- HTTP CSV reader
- row parsing
- normalization
- upsert persistence
- ingestion audit

## Functional Requirements
- read CSV from remote HTTP URL
- parse rows into internal structures
- create or reuse competition, season and teams
- upsert matches
- upsert match odds when available
- register run metadata in `ingestion_runs`
- support reruns without duplicates

## Non-Functional Requirements
- idempotency
- observability through logs and counters
- robustness against network and parsing failures
- traceability through audit records

## Acceptance Criteria
- successful ingestion persists consistent data
- rerunning the same source does not duplicate matches
- failed runs are recorded as failed in audit table

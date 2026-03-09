# Data Model and Migrations Spec

## Objective
Define the relational schema required to support the MVP ingestion and analytics workflow.

## Scope
- tables
- keys and constraints
- indexes
- migration strategy

## MVP Tables
- competitions
- seasons
- teams
- matches
- match_odds
- ingestion_runs

## Functional Requirements
- store competitions
- store seasons
- store teams
- store matches with home, away, date and result
- store match odds when available
- record ingestion audit data

## Non-Functional Requirements
- referential integrity
- idempotency support through unique constraints
- query efficiency for latest-N match lookups
- readiness for future evolution

## Acceptance Criteria
- schema migrates from empty database
- duplicate matches are prevented by database design
- indexes support intended queries

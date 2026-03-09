# Master Spec

## Purpose
Define the system as a whole, the implementation order and the global criteria for success.

## Product Objective
Build a personal football analytics system capable of ingesting structured match data, storing it in PostgreSQL and exposing useful aggregated metrics through Metabase and a future HTTP API/UI.

## Initial Scope
- Competition: Premier League
- Source: Football-Data.co.uk
- Storage: PostgreSQL
- Visualization: Metabase
- Backend: Go
- Local orchestration: Docker Compose

## Out of Scope (Initial Phase)
- live HTML scraping as primary source
- player-level statistics
- xG
- authentication
- predictive models

## Global Success Criteria
- system starts with Docker Compose
- schema can be recreated from zero
- ingestion is idempotent
- useful dashboards exist
- architecture supports future evolution

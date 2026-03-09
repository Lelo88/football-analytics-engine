# HTTP API Spec

## Objective
Expose a minimal JSON API for analytics queries.

## Scope
- list teams
- team form endpoint
- over/under endpoint
- season summary endpoint

## Functional Requirements
- `GET /teams`
- `GET /teams/{id}/form`
- `GET /teams/{id}/overunder`
- `GET /teams/{id}/season-summary`
- validate query parameters
- return meaningful HTTP errors

## Non-Functional Requirements
- stable JSON contracts
- thin handlers
- no direct DB access from delivery
- extensibility for future UI

## Acceptance Criteria
- endpoints respond with consistent JSON
- invalid input returns 400
- not found returns 404
- unexpected failures return 500

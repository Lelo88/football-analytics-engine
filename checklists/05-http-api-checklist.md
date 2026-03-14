# HTTP API Checklist

## Context
- [x] Read `specs/05-http-api.md`

## API Design
- [x] Choose router
- [x] Define route structure
- [x] Define JSON response format
- [x] Define error response format

## Endpoints
- [x] Implement `GET /teams`
- [x] Implement `GET /teams/{id}/form`
- [x] Implement `GET /teams/{id}/overunder`
- [x] Implement `GET /teams/{id}/season-summary`

## Validation
- [x] Validate team id
- [x] Validate `last_n`
- [x] Validate filter values
- [x] Validate threshold parameter
- [x] Return 400 for invalid input
- [x] Return 404 when team is not found
- [x] Return 500 for unexpected failures

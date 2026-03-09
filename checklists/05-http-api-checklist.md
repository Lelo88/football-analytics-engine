# HTTP API Checklist

## Context
- [ ] Read `specs/05-http-api.md`

## API Design
- [ ] Choose router
- [ ] Define route structure
- [ ] Define JSON response format
- [ ] Define error response format

## Endpoints
- [ ] Implement `GET /teams`
- [ ] Implement `GET /teams/{id}/form`
- [ ] Implement `GET /teams/{id}/overunder`
- [ ] Implement `GET /teams/{id}/season-summary`

## Validation
- [ ] Validate team id
- [ ] Validate `last_n`
- [ ] Validate filter values
- [ ] Validate threshold parameter
- [ ] Return 400 for invalid input
- [ ] Return 404 when team is not found
- [ ] Return 500 for unexpected failures

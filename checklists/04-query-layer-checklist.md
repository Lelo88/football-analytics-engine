# Query Layer Checklist

## Context
- [x] Read `specs/04-query-layer.md`

## Query Definitions
- [x] Define output contract for team form
- [x] Define output contract for goals summary
- [x] Define output contract for over/under
- [x] Define output contract for season summary

## Filters
- [x] Support `last_n`
- [x] Support home-only filter
- [x] Support away-only filter
- [x] Support all matches filter
- [x] Support season filtering

## Repository Design
- [x] Create read repository interfaces
- [x] Implement SQL queries in infrastructure layer
- [x] Ensure queries return domain-friendly structures

## Tests
- [x] Test team form correctness
- [x] Test home/away filtering
- [x] Test over/under correctness
- [x] Test season summary correctness

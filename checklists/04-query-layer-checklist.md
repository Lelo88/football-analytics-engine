# Query Layer Checklist

## Context
- [ ] Read `specs/04-query-layer.md`

## Query Definitions
- [ ] Define output contract for team form
- [ ] Define output contract for goals summary
- [ ] Define output contract for over/under
- [ ] Define output contract for season summary

## Filters
- [ ] Support `last_n`
- [ ] Support home-only filter
- [ ] Support away-only filter
- [ ] Support all matches filter
- [ ] Support season filtering

## Repository Design
- [ ] Create read repository interfaces
- [ ] Implement SQL queries in infrastructure layer
- [ ] Ensure queries return domain-friendly structures

## Tests
- [ ] Test team form correctness
- [ ] Test home/away filtering
- [ ] Test over/under correctness
- [ ] Test season summary correctness

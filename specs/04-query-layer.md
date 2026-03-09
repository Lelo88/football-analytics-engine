# Query Layer Spec

## Objective
Build reusable aggregated queries for team analytics.

## Scope
- team form
- goals for/against
- over/under
- season summary
- home/away filtering
- latest-N filtering

## Functional Requirements
- compute team form for latest N matches
- calculate points, goals for and goals against
- filter by home, away or all
- compute over/under statistics for a configurable threshold
- return season summaries per team

## Non-Functional Requirements
- correctness
- efficiency
- reuse by both API and dashboards
- testability

## Acceptance Criteria
- results match known fixture data
- filters behave correctly
- logic is not duplicated in handlers

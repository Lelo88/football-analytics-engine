# Master Checklist

## Global Execution Rules
- [ ] Read `docs/system-overview.md`
- [ ] Read `docs/architecture-rules.md`
- [ ] Read `docs/coding-rules.md`
- [ ] Read `docs/development-workflow.md`
- [ ] Read `specs/00-master-spec.md`

## Module Order
- [ ] Complete project setup
- [ ] Complete data model and migrations
- [ ] Complete ingestion pipeline
- [ ] Complete query layer
- [ ] Complete HTTP API
- [ ] Complete Metabase dashboards
- [ ] Complete observability and quality baseline
- [ ] Complete lightweight UI only after previous modules are stable

## Global Validation
- [ ] Project runs with Docker Compose
- [ ] Database schema can be recreated from zero
- [ ] Ingestion is idempotent
- [ ] Queries return meaningful results
- [ ] API exposes initial analytics
- [ ] Metabase shows useful dashboards
- [ ] Documentation is updated to match implementation

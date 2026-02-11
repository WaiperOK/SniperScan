# Changelog

## v0.2.0 - 2026-02-11
- Replaced legacy Windows-only scanner with a cross-platform Go architecture.
- Added concurrent TCP scanner with configurable timeout/workers and banner grabbing.
- Added HTTP API mode (`/healthz`, `/v1/scan`, `/metrics`) and Prometheus counters.
- Added unit/integration tests, CI, Docker stack, and architecture/security documentation.

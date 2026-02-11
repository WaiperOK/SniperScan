# Threat Model

## Assets
- Scan reports and target metadata.
- Scanner runtime resources in CI/self-hosted environments.

## Threats
- Misuse for unauthorized scanning.
- Resource exhaustion via oversized port specs.
- Fingerprinting data leakage in shared logs.

## Mitigations
- Explicit target+port requirements and bounded timeout settings.
- Configurable worker concurrency.
- Recommendation: enforce allowlisted target domains/IP ranges in production.

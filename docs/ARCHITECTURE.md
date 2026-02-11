# SniperScan Architecture

## Components
- Port parser: deterministic expansion of user port specs.
- Scan engine: high-concurrency TCP dial workers with timeout control.
- Fingerprint layer: lightweight service hints from port + banner.
- API service: endpoint for remote orchestration and report export.
- Metrics: Prometheus counters for scans and discovered open ports.

## Data flow
```mermaid
flowchart LR
  A["CLI scan command"] --> B["Scan Engine"]
  C["API /v1/scan"] --> B
  B --> D["Port Results"]
  B --> E["Summary JSON"]
  B --> F["Metrics /metrics"]
```

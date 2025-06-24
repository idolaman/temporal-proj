# Crawler Demo

This repository contains a two-service Go application orchestrated with Temporal and deployed via Helm to the built-in Kubernetes cluster that ships with Docker Desktop.

## Services

| Service | Description |
|---------|-------------|
| **service1** | Temporal **worker** that crawls a single URL, extracts all outbound links that belong to the same host, and returns them as the workflow result. |
| **service2** | REST **API** that triggers the scan workflow and stores results in Postgres (`scans` / `links` tables). |

A lightweight Postgres Pod and a Temporal Server Pod (using `temporalio/auto-setup`) are included in the Helm chart for local development.

## Quick start

### Prerequisites
* Docker Desktop with Kubernetes enabled (Mac/Win)
* Helm 3.x
* Go 1.23+

### Build & deploy

Windows PowerShell:
```powershell
./deploy.ps1
```

macOS/Linux:
```bash
./deploy.sh
```

Both scripts:
1. Switch the `kubectl` context to `docker-desktop`.
2. Build `crawler/service1` and `crawler/service2` images.
3. `helm upgrade --install crawler ./helm-chart`.

### Test the API

```
curl -X POST \
  --json '{"url":"https://en.wikipedia.org/wiki/Iran"}' \
  http://localhost:30080/scan
```

*The NodePort (30080) is defined in `helm-chart/values.yaml`; change it there if needed.*

## Database schema
The API auto-migrates two tables on startup:

* `scans`  – stores one row per requested URL.
* `links`  – stores each link found (`scan_id` → `scans.id`).

## Repository structure

```
nomaproj/
  ├─ service1/           # Temporal worker
  ├─ service2/           # REST API & Temporal client
  ├─ pkg/
  │   ├─ models/         # Shared structs (ScanTask, ScanResult)
  │   └─ utils/          # Small helper functions (env vars)
  ├─ helm-chart/         # Helm chart that deploys Postgres, Temporal, service pods
  ├─ deploy.ps1          # Windows deployment helper
  └─ deploy.sh           # macOS/Linux deployment helper
```
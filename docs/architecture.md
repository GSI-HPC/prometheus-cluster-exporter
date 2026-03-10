# Internal Architecture Overview

This exporter is a **meta-exporter**: it acts as a bridge between raw Lustre filesystem metrics already stored in an upstream Prometheus instance and HPC cluster context (SLURM jobs, OS users/groups). The result is a set of enriched Prometheus metrics that answer questions such as *"how much Lustre I/O is SLURM account `projectX` consuming right now?"*

---

## Key Components

| Component | File | Role |
|---|---|---|
| Entry point & PromQL queries | `main.go` | Parses flags, registers the collector, defines the three PromQL query templates |
| HTTP client | `client_prom_http.go` | Issues GET requests to the upstream Prometheus HTTP API |
| SLURM client | `client_slurm_squeue.go` | Runs `squeue` to list running jobs |
| User/group client | `client_getent.go` | Runs `getent passwd` / `getent group` to build UID→user and GID→group maps |
| Collector / correlator | `exporter.go` | Implements `prometheus.Collector`; fetches, parses, correlates, and emits all metrics |

---

## PromQL Queries

Three queries are hardcoded in `main.go` with a configurable `__TIME_RANGE__` placeholder (default `1m`):

| Purpose | Decoded PromQL |
|---|---|
| Metadata operations | `round(sum by(target,jobid)(irate(lustre_job_stats_total[1m])>=1))` |
| Read throughput | `sum by(jobid)(irate(lustre_job_read_bytes_total[1m])!=0)` |
| Write throughput | `sum by(jobid)(irate(lustre_job_write_bytes_total[1m])!=0)` |

These are sent as URL-encoded query strings to the upstream Prometheus `/api/v1/query` endpoint via `httpRequest()` in `client_prom_http.go`.

---

## Data Flow

```
Prometheus scrapes :9846/metrics
        │
        ▼
  exporter.Collect()
        │
        ├──[goroutine]──► squeue -ah -o "%A %a %u" ──► jobID→{account, user}
        ├──[goroutine]──► getent passwd ─────────────► UID→{username, GID}
        └──[goroutine]──► getent group ──────────────► GID→groupname
        │
        ▼  (wait for all 3 results on channels)
        │
        ├──► HTTP GET upstream Prometheus → parse JSON → metadataInfo[]
        ├──► HTTP GET upstream Prometheus → parse JSON → throughputInfo[] (read)
        └──► HTTP GET upstream Prometheus → parse JSON → throughputInfo[] (write)
        │
        ▼
   For each jobid in Lustre results:
     if numeric ──► match SLURM job ──► emit cluster_job_* {account, user}
     else        ──► split "procname.uid" ──► lookup getent ──► emit cluster_proc_* {proc_name, group_name, user_name}
        │
        ▼
   Push GaugeVec metrics into Prometheus channel
```

---

## Parsing and Correlation

The Lustre `jobid` label produced by the Lustre exporter takes two forms:

- **Plain integer** (e.g. `"12345"`) — a SLURM job ID. The exporter looks it up in the `squeue` result and emits `cluster_job_*` metrics labelled with `account` and `user`.
- **`procname.uid`** (e.g. `"mpirun.1001"`) — a non-SLURM process. The exporter splits on `.`, resolves the UID via the `getent` map, and emits `cluster_proc_*` metrics labelled with `proc_name`, `group_name`, and `user_name`.

For metadata metrics only MDT targets matching the pattern `^.*-MDT[[:xdigit:]]{4}$` (e.g. `lustre-MDT0000`) are kept; OST targets are skipped.

---

## Concurrency and Scrape Guard

`exporter.go` uses a `sync.Mutex` and a `scrapeActive bool` flag to prevent overlapping scrapes. If a scrape is still in progress when Prometheus polls again, the new request is skipped immediately and `cluster_exporter_scrape_ok` is set to `0`. A warning is logged:

> *"Collect is still active... - Skipping now"*

The three data-gathering operations (SLURM + two `getent` calls) run as concurrent goroutines and communicate results back over buffered channels. The three metric-building stages then execute sequentially; each stage's wall-clock time is recorded in `cluster_exporter_stage_execution_seconds`.

---

## Metrics Summary

All cluster metrics are prefixed with `cluster_`.

| Metric | Labels | Description |
|---|---|---|
| `cluster_exporter_scrape_ok` | — | `1` if scrape succeeded, `0` if skipped or failed |
| `cluster_exporter_stage_execution_seconds` | `name` | Wall-clock time per metric-building stage |
| `cluster_job_metadata_operations` | `account`, `user`, `target` | Metadata ops for SLURM jobs per MDT |
| `cluster_proc_metadata_operations` | `proc_name`, `group_name`, `user_name`, `target` | Metadata ops for non-SLURM processes per MDT |
| `cluster_job_read_throughput_bytes` | `account`, `user` | Read throughput for SLURM jobs (bytes/s) |
| `cluster_job_write_throughput_bytes` | `account`, `user` | Write throughput for SLURM jobs (bytes/s) |
| `cluster_proc_read_throughput_bytes` | `proc_name`, `group_name`, `user_name` | Read throughput for non-SLURM processes (bytes/s) |
| `cluster_proc_write_throughput_bytes` | `proc_name`, `group_name`, `user_name` | Write throughput for non-SLURM processes (bytes/s) |

---

## AI Disclosure & Authorship

- The original code analysis that this document is based on was performed using **Claude Sonnet 4.5** as a starting point.
- This document was assembled and written by the **GitHub Copilot assistant** (`@copilot`) based on that analysis and direct inspection of the source code.

The content reflects the code at the time of writing; if the implementation changes, this document should be updated accordingly.

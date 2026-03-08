# OpenTelemetry End-to-End Lab (Flask + Spring Boot)

This project is a practical observability playground for real services, not just hello-world demos.

It includes:
- Automatic instrumentation for Python (Flask) and Java (Spring Boot)
- Manual instrumentation for trace enrichment and business-level spans
- OpenTelemetry Collector pipelines
- Jaeger for trace exploration
- Prometheus for metrics
- LGTM stack (Loki, Grafana, Tempo, and metrics backend) for full observability workflows

The goal is simple: run the app, generate traffic, and debug behavior with traces, metrics, and logs from the same environment.

## Architecture

- `todobackend-springboot` handles Todo APIs
- `todoui-flask` consumes backend APIs and serves the UI
- Apps export telemetry via OTLP
- OpenTelemetry Collector receives/processes/exports telemetry
- Backends store and visualize telemetry (Jaeger + Prometheus + LGTM)

## Project Layout

- `automatic-instrumentation/`: auto-instrumented app versions
- `manual-instrumentation/`: manually instrumented examples
- `custom-automatic-instrumentation/`: profile-driven variants (initial/manual-tracing)
- `otel-collector-config.yml`: main Collector config
- `otel-collector-config-save-costs.yaml`: filtering/sampling-oriented config
- `prometheus.yml`: metrics scraping config
- `docker-compose.yaml`: main stack and automatic instrumentation profile
- `custom-automatic-instrumentation/docker-compose.yaml`: initial/manual-tracing profiles
- `Makefile`: unified commands for compose + profiles

## Prerequisites

- Docker + Docker Compose (plugin or standalone)
- GNU Make

## Quick Start

From project root:

```bash
make help
make profiles
```

Run automatic instrumentation stack:

```bash
make up PROFILE=automatic-instrumentation
```

Run manual tracing profile (custom compose):

```bash
make up PROFILE=manual-tracing COMPOSE_FILE=custom-automatic-instrumentation/docker-compose.yaml
```

Run initial baseline profile (custom compose):

```bash
make up PROFILE=initial COMPOSE_FILE=custom-automatic-instrumentation/docker-compose.yaml
```

Stop a running profile:

```bash
make down PROFILE=automatic-instrumentation
```

## LGTM Stack Workflow

The LGTM environment is integrated for production-like observability workflows:
- **Loki** for logs
- **Grafana** for dashboards and exploration
- **Tempo** for distributed tracing
- **Metrics backend** for PromQL-based analysis

Typical flow:
1. Start LGTM profile
2. Start instrumented services
3. Generate app traffic
4. Correlate logs, traces, and metrics in Grafana

## Useful Endpoints

- Todo UI: `http://localhost:7000` (or `http://localhost:5000` depending on profile)
- Spring Boot backend: `http://localhost:8080`
- Jaeger UI: `http://localhost:16686`
- Prometheus UI: `http://localhost:9090`
- Grafana UI: `http://localhost:3000`

## Observability Modes

- **Automatic Instrumentation**
  Fastest way to get telemetry from existing services with minimal code changes.

- **Manual Instrumentation**
  Best for domain-specific spans, custom attributes, and richer trace context.

- **Cost-aware Collection**
  Use the save-costs collector config to drop noisy metrics/spans/logs and reduce telemetry volume.

## Security and Container Hardening

Docker images are hardened with practical defaults:
- Multi-stage builds
- Non-root runtime users
- Minimal runtime dependencies
- Cleaner package installation (`--no-cache-dir`, reduced layer leftovers)
- Safer artifact handling for OpenTelemetry agent setup

## Common Commands

```bash
# list profiles available in the selected compose file
make profiles
make profiles COMPOSE_FILE=custom-automatic-instrumentation/docker-compose.yaml

# build and run
make build PROFILE=automatic-instrumentation
make up PROFILE=automatic-instrumentation

# inspect
make ps PROFILE=automatic-instrumentation
make logs PROFILE=automatic-instrumentation
make config PROFILE=automatic-instrumentation

# shutdown
make down PROFILE=automatic-instrumentation
```

## Why This Repo Exists

This repository is focused on hands-on OpenTelemetry learning for real troubleshooting scenarios:
- understand service latency using traces
- validate throughput and saturation with metrics
- investigate incidents through logs correlated with trace IDs

If you are preparing for OpenTelemetry certification or building your first production-ready observability baseline, this project is designed for that path.

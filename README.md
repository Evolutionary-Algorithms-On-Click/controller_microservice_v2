# Controller Microservice V2

Source code of controller microservice that handles data to and from frontend and jupyter kernel gateway for evolutionary algorithms on click v2.

## Overview

This service acts as the central orchestrator for the platform, managing:
- [User Sessions & Authentication](https://github.com/Evolutionary-Algorithms-On-Click/auth_microservice)
- Jupyter Kernel Lifecycle (via `jupyter_gateway`)
- Database Interactions (CockroachDB)
- Inter-service communication (gRPC, HTTP)
- [LLM microservice ](https://github.com/Evolutionary-Algorithms-On-Click/evocv2_llm_microservice)
- [Volpe Integration service](https://github.com/Evolutionary-Algorithms-On-Click/volpe-integration)

## Prerequisites

- **Go**: 1.24.1 or higher
- **Docker**: For containerization and dependencies (CockroachDB, MinIO, )
- **Docker Compose**: For orchestration
- **Make**: For running standard commands
- **Lefthook**: For git hooks (optional but recommended)

## Getting Started

### 1. Environment Setup

1.  Clone the repository.
2.  Ensure you have the necessary environment variables set. Refer to `docker-compose.yaml` for required keys (e.g., `DATABASE_URL`, `JUPYTER_GATEWAY_URL`).
    *Note: The project uses `godotenv` to load environment variables from a `.env` file in development.*

### 2. Running Dependencies

Start the supporting services (Database, Object Storage, Jupyter Gateway, etc.):

```bash
make docker-up
```

or 

```bash
docker-compose up --build
```


This will spin up CockroachDB, MinIO, Jupyter Gateway, and Python Runner as defined in `docker-compose.yaml`.

### 3. Local Development

To run the controller service locally:

```bash
# Install tools and git hooks
make setup

# Build and run the application
make run
```

The service will start on port `8080` (default).

## API Documentation

- **HTTP API**: Versioned under `/api/v1/`. Defined in `routes/api.go`.
- **gRPC**: Defined in `proto/authenticate.proto`.

# Job Crawler

Job Crawler is a full-stack job aggregation platform designed to help users discover, triage, and track relevant employment opportunities from multiple job boards in a single interface. The system includes a web application, API service, and background worker for scheduled scraping, enrichment, and notification delivery.

## Authors

- Tyler Mestery  
- Mason Hart

## Project Purpose

This project centralizes job discovery by:

- collecting and normalizing postings from major job sources,
- filtering results against user preferences,
- supporting profile-based workflows (privacy, resume, history),
- tracking application decisions over time, and
- delivering periodic email digests based on user-selected notification cadence.

## System Architecture

The repository is organized as a multi-service system:

- `apps/web` — React + TypeScript frontend
- `apps/api` — Go API service
- `apps/worker` — Go background worker (scraping + digest orchestration)
- `docker-compose.yml` — local orchestration for core services

Persistent data is stored in MySQL for authentication, profile data, feed decisions, and notification metadata.

## Core Features

- User signup and login
- Profile management (privacy, resume metadata, notification preferences)
- Job feed with filtering and sorting
- Apply/reject decision workflow
- Application history by user profile
- Worker-driven scraping from prioritized job sources
- Notification digest generation and SMTP email sending

## Technology Stack

- **Frontend:** React, TypeScript, Vite
- **Backend:** Go (API + Worker)
- **Database:** MySQL
- **Containerization:** Docker, Docker Compose

## Local Development

### Prerequisites

- Docker Desktop (or equivalent Docker runtime)
- Git

### Start the Full Stack

```bash
docker compose up --build
```

### Service Endpoints

- Web: `http://localhost:5173`
- API: `http://localhost:8080`
- MySQL: `localhost:3306`

### Stop Services

```bash
docker compose down
```

## Configuration Notes

The worker supports SMTP configuration for notification email delivery via environment variables:

- `WORKER_SMTP_HOST`
- `WORKER_SMTP_PORT`
- `WORKER_SMTP_USERNAME`
- `WORKER_SMTP_PASSWORD`
- `WORKER_SMTP_FROM`

If SMTP settings are not configured, scraping and ingestion continue to function, while email sending is skipped safely.

## Quality and Validation

The project includes automated test suites for API, worker, and web modules. Recommended validation commands:

```bash
cd apps/api && go test ./...
cd apps/worker && go test ./...
npm --prefix apps/web test -- --run
```

## License

This repository includes a `LICENSE` file at the project root. Refer to it for usage and distribution terms.

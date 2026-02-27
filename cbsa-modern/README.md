# CBSA Modern — Quick Reference

The modernized Go + PostgreSQL + React/TypeScript banking application.
For the full project overview, migration details, and performance comparisons, see the [root README](../README.md).

## Quick Start

### Docker (recommended — one command)

```bash
docker-compose up --build
```

This starts:
- **PostgreSQL** on port 5432 (schema auto-migrates on startup)
- **Go API** on port 8080
- **React frontend** on port 3000 (Nginx proxies `/api` to the Go backend)

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Local Development (with hot reload)

```bash
# Terminal 1 — database only
docker-compose up postgres

# Terminal 2 — Go API
go run ./cmd/server

# Terminal 3 — React frontend
cd frontend
npm install --legacy-peer-deps   # or: bun install
npm start                        # or: bun start
```

- API: [http://localhost:8080](http://localhost:8080)
- Frontend: [http://localhost:3000](http://localhost:3000)

### Seed Test Data

```bash
go run ./cmd/seed
```

Creates 100 customers with up to 5 accounts each.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Go API listen port |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `cbsa` | PostgreSQL user |
| `DB_PASSWORD` | `cbsa` | PostgreSQL password |
| `DB_NAME` | `cbsa` | PostgreSQL database name |
| `COMPANY_NAME` | `CBSA Modern Bank` | Bank name shown in the UI |
| `SORT_CODE` | `987654` | Bank sort code |

The frontend uses `REACT_APP_API_BASE_URL` (defaults to `http://localhost:8080/api/v1`), configured in `frontend/.env`.

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/customers` | List customers |
| POST | `/api/v1/customers` | Create customer |
| GET | `/api/v1/customers/:id` | Get customer |
| PUT | `/api/v1/customers/:id` | Update customer |
| DELETE | `/api/v1/customers/:id` | Delete customer |
| GET | `/api/v1/accounts` | List accounts |
| POST | `/api/v1/accounts` | Create account |
| GET | `/api/v1/accounts/:id` | Get account |
| PUT | `/api/v1/accounts/:id` | Update account |
| DELETE | `/api/v1/accounts/:id` | Delete account |
| GET | `/api/v1/accounts/customer/:id` | Accounts by customer |
| GET | `/api/v1/transactions` | List transactions |
| PUT | `/api/v1/transactions/debit` | Debit account |
| PUT | `/api/v1/transactions/credit` | Credit account |
| PUT | `/api/v1/transactions/transfer` | Transfer between accounts |

## Running Tests

```bash
# Unit tests
go test ./... -v

# Benchmarks
go test ./... -bench=. -benchmem -run=^$
```

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.21, chi router, sqlx, shopspring/decimal |
| Database | PostgreSQL 16 |
| Frontend | React 18, TypeScript, Carbon Design System |
| Containerization | Docker, Docker Compose, Nginx |

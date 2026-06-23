# monyLonger

A personal finance tracker built with Go and vanilla JavaScript. Track your spending across budget categories, monitor income vs. expenses, and get a monthly overview — all in a dark-themed UI without any external frameworks.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat&logo=go) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=flat&logo=postgresql)

---

## Features

- **Dashboard** — monthly KPIs (balance, spent, income, over-budget count), budget progress bars, recent transactions, income vs. expense bar chart
- **Budgets** — create, edit, and delete budget envelopes with emoji icons and custom colors; real-time remaining balance
- **Transactions** — full CRUD, search, filter by budget/month/type (income/expense), pagination, CSV export
- **Event-driven balance** — creating a transaction automatically adjusts the linked budget's remaining balance asynchronously
- **Dark UI** — pure HTML + CSS + vanilla JS, no build step required

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.25, `net/http` (stdlib router) |
| Database | PostgreSQL 16 |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| DB driver | [lib/pq](https://github.com/lib/pq) |
| Frontend | Vanilla JS, Google Fonts (DM Serif Display, Instrument Sans, DM Mono) |
| Tests | [testcontainers-go](https://testcontainers.com/guides/getting-started-with-testcontainers-for-go/) |

---

## Project Structure

```
monyLonger/
├── cmd/
│   └── main.go                  # Entry point, route registration, server start
├── internal/
│   ├── domain/
│   │   ├── transaction.go       # Transaction model + repository interface + event
│   │   └── vault.go             # Vault (budget) model + repository interface
│   ├── handler/
│   │   ├── dashboard.go         # GET /api/dashboard
│   │   ├── transaction.go       # CRUD handlers for /api/transactions
│   │   └── vault.go             # CRUD handlers for /api/budgets
│   └── storage/
│       ├── db.go                # DB connection, migrations, connection pool
│       ├── transaction_storage.go
│       └── vault_storage.go
├── migrations/
│   ├── 000001_init.up.sql
│   └── 000001_init.down.sql
├── public/
│   ├── index.html               # Dashboard page
│   ├── budgets.html             # Budgets page
│   ├── transactions.html        # Transactions list
│   ├── transactions-new.html    # New transaction form
│   └── transactions-edit.html   # Edit transaction form
└── tests/
    └── transaction_func_test.go # Integration tests (testcontainers)
```

---

## Getting Started

### Prerequisites

- Go 1.22+
- Docker (for the database)

### 1. Start PostgreSQL

```bash
docker run -d \
  --name monyLonger-db \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=monyLonger \
  -p 5433:5432 \
  postgres:16-alpine
```

> The app binds to port **5433** by default to avoid conflicts with a local PostgreSQL instance.

### 2. Run the server

```bash
go run ./cmd/main.go
```

The server starts at [http://localhost:8090](http://localhost:8090).

Migrations run automatically on startup — no manual setup needed.

---

## Configuration

All configuration is via environment variables.

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `host=localhost port=5433 user=admin password=secret dbname=monyLonger sslmode=disable` | PostgreSQL connection string |

Example with a custom DSN:

```bash
DATABASE_URL="host=db port=5432 user=app password=pass dbname=finance sslmode=require" \
  go run ./cmd/main.go
```
---

## Useful Docker Commands

```bash
# Connect to the database
docker exec -it monyLonger-db psql -U admin -d monyLonger

# Stop / start the container
docker stop monyLonger-db
docker start monyLonger-db
```

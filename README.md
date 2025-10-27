# Treasury Management Application

ModernFi take-home assignment - A full-stack application for managing treasury securities and liquidity.

## Overview

This application provides real-time treasury yield curve visualization and allows users to manage treasury security portfolios with buy/sell functionality and comprehensive transaction tracking.

## Features

- Live treasury yield data from Treasury.gov API (1-hour cache)
- Interactive yield curve visualization
- User account management with fund/withdraw operations
- Buy treasury securities (bills, notes, bonds)
- Sell holdings with yield calculations
- Complete transaction history
- Mobile-responsive interface

## Tech Stack

**Backend:** Go 1.21+, Chi router, PostgreSQL driver (pgx), sqlc for type-safe SQL
**Frontend:** React 18, TypeScript, Vite, Tailwind CSS, Tremor charts
**Database:** PostgreSQL 16
**Deployment:** Docker Compose

## Prerequisites

- Docker Desktop 20.10+ with Docker Compose v2
- (For local development) Go 1.21+, Node.js 18+, npm

## Quick Start (Docker - Recommended)

Run the entire application with a single command:

```bash
docker compose up -d
```

Access the application:
- Frontend: http://localhost
- Backend API: http://localhost:8080
- Health check: http://localhost:8080/health

Stop the application:
```bash
docker compose down
```

## Development Setup

For faster development with hot reloading:

### 1. Start PostgreSQL
```bash
docker compose up -d postgres
```

### 2. Deploy Database Schema
```bash
# Deploy schema with demo data
cd backend/db
./deploy-fresh-schema.sh local

# Verify deployment
docker exec treasury_postgres psql -U postgres -d treasury_db -c "SELECT name, balance FROM users;"
```

### 3. Start Backend
```bash
cd backend
go run ./cmd/server/main.go
```

Backend runs on http://localhost:8080

### 4. Start Frontend
```bash
cd frontend
npm install
npm run dev
```

Frontend runs on http://localhost:5173

## Building Docker Images

To build images manually:

```bash
# Build all services
docker compose build

# Build specific service
docker compose build backend
docker compose build frontend
```

## API Endpoints

- `GET /api/yields` - Current treasury yield curve data
- `GET /api/yields/historical` - Historical yield data for charting
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/{userId}/transactions` - User transaction history
- `GET /api/v1/users/{userId}/holdings` - User active holdings
- `POST /api/v1/fund` - Add funds to account
- `POST /api/v1/withdraw` - Withdraw funds from account
- `POST /api/v1/buy` - Purchase treasury security
- `POST /api/v1/sell` - Sell treasury holding
- `GET /health` - Backend health check

## Database Schema

The application uses PostgreSQL with the following main tables:
- `users` - User accounts with balances
- `transactions` - All financial transactions (fund, withdraw, buy, sell)
- `holdings` - Treasury security holdings with remaining amounts

Schema is in `backend/db/schema.sql` and is automatically applied via Docker.

## Project Structure

```
.
├── backend/
│   ├── cmd/server/          # Main application entry point
│   ├── internal/
│   │   ├── database/        # sqlc generated code
│   │   ├── handlers/        # HTTP request handlers
│   │   ├── services/        # Business logic
│   │   └── utils/           # Utilities (yield calculations)
│   └── db/
│       ├── queries/         # SQL queries for sqlc
│       └── schema.sql       # Database schema
├── frontend/
│   └── src/
│       ├── components/      # React components
│       ├── services/        # API client
│       ├── contexts/        # React contexts
│       └── types/           # TypeScript types
└── docker-compose.yml       # Docker orchestration
```

## Testing

```bash
# Backend tests
cd backend
go test -v ./...

# Frontend tests
cd frontend
npm test
```

## Troubleshooting

### Port 80 already in use
```bash
# Check what's using port 80
lsof -i :80

# Or change the port in docker-compose.yml:
# ports:
#   - "8080:80"  # Access on localhost:8080 instead
```

### Reset database
```bash
docker compose down -v
docker compose up -d
```

### View logs
```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f backend
docker compose logs -f frontend
```

## Video Demo

A 30-second screen recording demonstrating the application functionality is included with this submission as required.

## Notes

- Initial user accounts are created via seed data with demo balances
- Treasury yield data is cached for 1 hour from Treasury.gov
- Buy orders for T-Bills use discount pricing (pay less than face value)
- Sell operations calculate accrued yield based on time held and current rates

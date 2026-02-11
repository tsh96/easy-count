# Easy Count Server - Golang + Gin Backend

Production-ready backend API server for Easy Count application using Golang, Gin framework, and PostgreSQL (Neon) for data storage.

## Features

✅ **Production-Ready Security**
- JWT-based authentication with refresh tokens
- Bcrypt password hashing (cost factor 10)
- Rate limiting (10 req/s per IP with burst of 20)
- Security headers (XSS, clickjacking, MIME sniffing protection)
- Strict CORS configuration
- Input validation and sanitization
- SQL injection prevention (parameterized queries)

✅ **High Performance**
- Golang's native concurrency
- Connection pooling with pgx
- Database transactions for bulk operations

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database (Neon recommended)

## Setup

1. Install dependencies:
```bash
go mod download
```

2. Create `.env` from `.env.example` and configure:
   - DATABASE_URL: Your Neon PostgreSQL connection string
   - JWT_SECRET: Strong random string (generate with `openssl rand -base64 32`)
   - PORT, ENVIRONMENT, ALLOW_ORIGINS

3. Run:
```bash
go run cmd/api/main.go
```

## API Endpoints

See full documentation in the README for all endpoints including authentication, transactions, customer records, and backup/restore.


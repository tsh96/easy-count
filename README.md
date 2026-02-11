# Easy Count

Production-ready accounting application for tracking personal transactions and customer records (invoices and cheques).

## Architecture

Modern client-server architecture with production-grade security:
- **Frontend**: Vue 3 application with TypeScript
- **Backend**: Golang + Gin framework with JWT authentication
- **Database**: PostgreSQL (Neon) with connection pooling
- **Deployment**: Docker support for easy containerized deployment

## Features

✅ **Security**
- JWT-based authentication with refresh tokens
- Bcrypt password hashing
- Rate limiting and security headers
- CORS protection
- SQL injection prevention

✅ **Functionality**
- Personal transaction tracking
- Customer records (Private & Government)
- Backup and restore
- Multi-user support with data isolation

✅ **Deployment**
- Docker support (single image with frontend + backend)
- Multiple cloud platform options
- Easy local development setup

## Quick Start

### Option 1: Using Docker (Recommended for Production)

```bash
# Clone repository
git clone https://github.com/tsh96/easy-count.git
cd easy-count

# Set up environment
cp .env.docker.example .env
# Edit .env with your DATABASE_URL and JWT_SECRET

# Build and run
docker-compose up --build

# Access at http://localhost:8080
```

See [DOCKER.md](DOCKER.md) for complete Docker documentation.

### Option 2: Local Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL database (Neon recommended)

### Backend Setup

1. Navigate to server directory:
```bash
cd server
```

2. Install dependencies:
```bash
go mod download
```

3. Create `.env` file:
```bash
cp .env.example .env
```

4. Configure environment variables:
   - **DATABASE_URL**: Your Neon PostgreSQL connection string
   - **JWT_SECRET**: Generate with `openssl rand -base64 32`
   - **ALLOW_ORIGINS**: Your frontend URL

5. Start server:
```bash
go run cmd/api/main.go
```

### Option 2: Local Development

**Prerequisites**
- Go 1.21+
- Node.js 18+
- PostgreSQL database (Neon recommended)

**Backend Setup**

1. Navigate to server directory:
```bash
cd server
```

2. Install dependencies:
```bash
go mod download
```

3. Create `.env` file:
```bash
cp .env.example .env
```

4. Configure environment variables:
   - **DATABASE_URL**: Your Neon PostgreSQL connection string
   - **JWT_SECRET**: Generate with `openssl rand -base64 32`
   - **ALLOW_ORIGINS**: Your frontend URL

5. Start server:
```bash
go run cmd/api/main.go
```

Server runs on `http://localhost:3001`

**Frontend Setup**

1. Return to root directory:
```bash
cd ..
```

2. Install dependencies:
```bash
pnpm install
```

3. Create `.env`:
```bash
cp .env.example .env
```

4. Start development server:
```bash
pnpm dev
```

Application available at `http://localhost:5173`

## Deployment

### Docker (Recommended)

Single command deployment with Docker:
```bash
docker build -t easy-count .
docker run -p 8080:8080 -e DATABASE_URL="..." -e JWT_SECRET="..." easy-count
```

See [DOCKER.md](DOCKER.md) for:
- Docker Compose setup
- Cloud deployment (Fly.io, Cloud Run, ECS)
- Production configuration
- Troubleshooting

### Other Options

See [DEPLOYMENT.md](DEPLOYMENT.md) for:
- Fly.io deployment
- Railway deployment
- Traditional VPS
- Vercel/Netlify (frontend)

## Documentation

- **[DOCKER.md](DOCKER.md)** - Docker deployment guide
- **[DEPLOYMENT.md](DEPLOYMENT.md)** - Cloud deployment options
- **[MIGRATION.md](MIGRATION.md)** - Migration from IndexedDB
- **[server/README.md](server/README.md)** - API documentation

## First Use

1. Open `http://localhost:5173`
2. Create an account (register)
3. Login with your credentials
4. Start tracking transactions and customer records

## API Documentation

See [server/README.md](server/README.md) for complete API documentation including:
- Authentication endpoints
- Protected endpoints
- Security features
- Request/response formats

## Deployment

See [DEPLOYMENT.md](DEPLOYMENT.md) for deployment options:
- Fly.io / Railway (Backend)
- Vercel / Netlify (Frontend)
- Traditional VPS
- Docker containers (see DOCKER.md)

## Migration from IndexedDB

If migrating from the old IndexedDB version:
1. Export data using Backup button in old version
2. Register new account in new version
3. Import data using Restore button

See [MIGRATION.md](MIGRATION.md) for details.

## Security

This version includes production-ready security features:
- **Authentication**: JWT tokens with 24h access and 7d refresh
- **Password**: Bcrypt hashing with strength validation
- **Rate Limiting**: 10 req/s per IP
- **Headers**: XSS, clickjacking, MIME sniffing protection
- **CORS**: Strict origin checking
- **Database**: Parameterized queries, user isolation

**Important**: Set a strong JWT_SECRET before production deployment.

## Development

### Backend
```bash
cd server
go run cmd/api/main.go
```

### Frontend
```bash
pnpm dev
```

### Build
```bash
# Backend
cd server
go build -o bin/server cmd/api/main.go

# Frontend
pnpm build
```

## Project Structure

```
.
├── server/              # Golang backend
│   ├── cmd/api/        # Application entry point
│   ├── internal/       # Internal packages
│   │   ├── config/    # Configuration
│   │   ├── db/        # Database layer
│   │   ├── handlers/  # HTTP handlers
│   │   ├── middleware/# Middleware (auth, CORS, etc.)
│   │   ├── models/    # Data models
│   │   └── utils/     # Utilities (JWT, bcrypt)
│   └── go.mod
├── src/                # Vue 3 frontend
│   ├── components/
│   ├── composables/   # Vue composables
│   ├── pages/         # Page components
│   └── main.ts
└── README.md
```

## Special Thanks

<a href="https://www.flaticon.com/free-icons/accounting" title="accounting icons">Accounting icons created by Freepik - Flaticon</a>

## License

MIT

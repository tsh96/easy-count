# Docker Deployment Guide

This guide explains how to build and run Easy Count using Docker.

## Overview

The Dockerfile creates a single image containing both the frontend (Vue 3) and backend (Golang + Gin). The backend serves the API endpoints and also serves the frontend static files.

## Quick Start

### Using Docker Compose (Recommended)

1. **Set up environment variables:**
   ```bash
   cp .env.docker.example .env
   # Edit .env and set your DATABASE_URL and JWT_SECRET
   ```

2. **Build and run:**
   ```bash
   docker-compose up --build
   ```

3. **Access the application:**
   - Open http://localhost:8080
   - Register a new account
   - Start using the app!

### Using Docker CLI

1. **Build the image:**
   ```bash
   docker build -t easy-count:latest .
   ```

2. **Run the container:**
   ```bash
   docker run -p 8080:8080 \
     -e DATABASE_URL="your-neon-db-url" \
     -e JWT_SECRET="your-secret-key" \
     -e ENVIRONMENT="production" \
     -e ALLOW_ORIGINS="http://localhost:8080" \
     easy-count:latest
   ```

3. **Access the application:**
   - Open http://localhost:8080

## Image Details

### Multi-Stage Build

The Dockerfile uses a multi-stage build process:

1. **Stage 1: Frontend Builder** (node:18-alpine)
   - Installs pnpm and frontend dependencies
   - Builds Vue 3 application → `dist/` directory

2. **Stage 2: Backend Builder** (golang:1.21-alpine)
   - Downloads Go dependencies
   - Compiles Golang binary

3. **Stage 3: Runtime** (alpine:latest)
   - Minimal Alpine Linux base (~5MB)
   - Contains only the compiled backend binary
   - Contains frontend static files
   - Final image size: ~50-60MB

### What's Included

- ✅ Backend API (Golang + Gin)
- ✅ Frontend static files (Vue 3)
- ✅ All security features (JWT, rate limiting, etc.)
- ✅ Static file serving from backend
- ✅ SPA fallback routing

### What's NOT Included

- ❌ PostgreSQL database (use Neon or external PostgreSQL)
- ❌ Development tools
- ❌ Source code (only compiled binaries)

## Environment Variables

Required:
- `DATABASE_URL` - PostgreSQL connection string (Neon recommended)
- `JWT_SECRET` - Secret key for JWT signing (32+ characters)

Optional:
- `PORT` - Server port (default: 8080)
- `ENVIRONMENT` - Set to "production" for production mode
- `ALLOW_ORIGINS` - CORS allowed origins (comma-separated)

## Configuration

### Using Neon Database

```bash
DATABASE_URL=postgres://user:pass@host.neon.tech/dbname?sslmode=require
JWT_SECRET=$(openssl rand -base64 32)
ENVIRONMENT=production
PORT=8080
ALLOW_ORIGINS=http://localhost:8080
```

### Using Local PostgreSQL

If you want to use PostgreSQL in Docker:

1. Uncomment the `postgres` service in `docker-compose.yml`
2. Set `DATABASE_URL=postgres://easycount:easycount@postgres:5432/easycount?sslmode=disable`

## Production Deployment

### Docker Hub

1. **Build and tag:**
   ```bash
   docker build -t yourusername/easy-count:latest .
   docker tag yourusername/easy-count:latest yourusername/easy-count:v1.0.0
   ```

2. **Push to Docker Hub:**
   ```bash
   docker login
   docker push yourusername/easy-count:latest
   docker push yourusername/easy-count:v1.0.0
   ```

3. **Pull and run on server:**
   ```bash
   docker pull yourusername/easy-count:latest
   docker run -p 80:8080 \
     -e DATABASE_URL="$DATABASE_URL" \
     -e JWT_SECRET="$JWT_SECRET" \
     -e ENVIRONMENT="production" \
     -e ALLOW_ORIGINS="https://yourdomain.com" \
     -d yourusername/easy-count:latest
   ```

### Cloud Platforms

#### Fly.io

```bash
# Install flyctl
curl -L https://fly.io/install.sh | sh

# Launch app
fly launch

# Set secrets
fly secrets set DATABASE_URL="your-db-url"
fly secrets set JWT_SECRET="your-secret"

# Deploy
fly deploy
```

#### Google Cloud Run

```bash
# Build and push to GCR
gcloud builds submit --tag gcr.io/PROJECT-ID/easy-count

# Deploy to Cloud Run
gcloud run deploy easy-count \
  --image gcr.io/PROJECT-ID/easy-count \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars DATABASE_URL="$DATABASE_URL",JWT_SECRET="$JWT_SECRET"
```

#### AWS ECS/Fargate

1. Push image to ECR:
   ```bash
   aws ecr create-repository --repository-name easy-count
   docker tag easy-count:latest AWS_ACCOUNT.dkr.ecr.REGION.amazonaws.com/easy-count:latest
   docker push AWS_ACCOUNT.dkr.ecr.REGION.amazonaws.com/easy-count:latest
   ```

2. Create ECS task definition with environment variables
3. Create ECS service using the task definition

## Development

### Local Development with Docker

```bash
# Build for development
docker build -t easy-count:dev .

# Run with volume mount for hot reload (not recommended, use native dev instead)
# Better to use: pnpm dev (frontend) + go run (backend) for development
```

**Note:** For development, it's better to run frontend and backend separately:
- Frontend: `pnpm dev`
- Backend: `cd server && go run cmd/api/main.go`

### Debugging

**View logs:**
```bash
docker-compose logs -f app
```

**Execute commands in container:**
```bash
docker-compose exec app sh
```

**Check health:**
```bash
curl http://localhost:8080/health
```

## Troubleshooting

### Build Issues

**Problem:** Frontend build fails
- **Solution:** Ensure you have enough memory (min 2GB) for Node.js build

**Problem:** Backend build fails
- **Solution:** Check Go version compatibility (requires Go 1.21+)

### Runtime Issues

**Problem:** Cannot connect to database
- **Solution:** Verify DATABASE_URL is correct and database is accessible
- Check that SSL mode is correct for your database

**Problem:** CORS errors
- **Solution:** Set ALLOW_ORIGINS to include your frontend URL
- Format: `ALLOW_ORIGINS=http://localhost:8080,https://yourdomain.com`

**Problem:** 404 for frontend routes
- **Solution:** This shouldn't happen with the SPA fallback, check logs

### Performance Issues

**Problem:** Slow response times
- **Solution:** 
  - Check database connection (use connection pooling)
  - Ensure database is in same region as container
  - Monitor rate limiting settings

## Image Size Optimization

The Docker image is optimized for size:
- Multi-stage build (only runtime artifacts in final image)
- Alpine Linux base (minimal footprint)
- No development dependencies
- Compiled binaries (no runtime compilation)

**Typical sizes:**
- Frontend build artifacts: ~10-15MB
- Backend binary: ~20-25MB
- Alpine + dependencies: ~10-15MB
- **Total: ~50-60MB**

## Security Considerations

1. **Environment Variables:** Never commit `.env` files to Git
2. **JWT Secret:** Use a strong random secret (32+ characters)
3. **Database:** Use SSL/TLS for database connections
4. **HTTPS:** Use reverse proxy (Nginx, Caddy) for HTTPS in production
5. **Updates:** Regularly rebuild image with updated dependencies

## Health Checks

The application includes a health check endpoint:

```bash
curl http://localhost:8080/health
# Expected: {"status":"ok"}
```

Use this for container health checks:

```yaml
# docker-compose.yml
services:
  app:
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

## Backup and Data

The Docker container is **stateless**. All data is stored in the PostgreSQL database.

**To backup:**
1. Backup your PostgreSQL database (Neon provides automatic backups)
2. Or use the built-in backup feature in the app

**To restore:**
1. Restore PostgreSQL database
2. Or use the built-in restore feature in the app

## Support

For issues related to:
- **Docker build/run:** Check this guide and Docker documentation
- **Application features:** See main README.md
- **Deployment:** See DEPLOYMENT.md
- **Database migration:** See MIGRATION.md

## License

Same as main project.

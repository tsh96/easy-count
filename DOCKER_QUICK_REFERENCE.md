# Docker Quick Reference

## Build and Run

### Using Docker Compose (Easiest)
```bash
# Start (builds if needed)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down

# Rebuild
docker-compose up --build
```

### Using Docker CLI
```bash
# Build
docker build -t easy-count:latest .

# Run
docker run -d \
  -p 8080:8080 \
  -e DATABASE_URL="your-neon-url" \
  -e JWT_SECRET="your-secret" \
  -e ENVIRONMENT="production" \
  --name easy-count \
  easy-count:latest

# View logs
docker logs -f easy-count

# Stop
docker stop easy-count

# Remove
docker rm easy-count
```

## Common Commands

```bash
# Check running containers
docker ps

# Check all containers
docker ps -a

# View logs
docker logs easy-count
docker logs -f easy-count  # follow

# Execute command in container
docker exec -it easy-count sh

# Check health
curl http://localhost:8080/health

# Restart container
docker restart easy-count

# View resource usage
docker stats easy-count
```

## Troubleshooting

```bash
# View build output
docker build -t easy-count:latest . --progress=plain

# Run with environment file
docker run -p 8080:8080 --env-file .env easy-count:latest

# Check container details
docker inspect easy-count

# Check networks
docker network ls
docker network inspect bridge
```

## Cleanup

```bash
# Remove stopped containers
docker container prune

# Remove unused images
docker image prune

# Remove all unused data
docker system prune

# Remove all (including volumes)
docker system prune -a --volumes
```

## Production

```bash
# Tag for registry
docker tag easy-count:latest username/easy-count:v1.0.0

# Push to Docker Hub
docker push username/easy-count:v1.0.0

# Pull on server
docker pull username/easy-count:v1.0.0

# Run in production
docker run -d \
  -p 80:8080 \
  -e DATABASE_URL="$DATABASE_URL" \
  -e JWT_SECRET="$JWT_SECRET" \
  -e ENVIRONMENT="production" \
  -e ALLOW_ORIGINS="https://yourdomain.com" \
  --restart unless-stopped \
  --name easy-count \
  username/easy-count:v1.0.0
```

## Environment Variables

```bash
# List environment variables in container
docker exec easy-count env

# Run with specific variables
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  -e JWT_SECRET="secret123" \
  -e PORT="8080" \
  easy-count:latest
```

## Health Checks

```bash
# Manual check
curl http://localhost:8080/health

# Check from inside container
docker exec easy-count wget -qO- http://localhost:8080/health

# View health status
docker inspect --format='{{.State.Health.Status}}' easy-count
```

For more details, see [DOCKER.md](DOCKER.md)

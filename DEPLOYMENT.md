# Deployment Guide - Golang + Gin Backend

## Prerequisites

- Neon database account (free tier: [neon.tech](https://neon.tech))
- Go 1.21+ (backend)
- Node.js 18+ (frontend build)

## Environment Variables

### Backend (.env in server/)
```bash
DATABASE_URL=postgres://user:password@host/database?sslmode=require
JWT_SECRET=<generate-with-openssl-rand-base64-32>
ENVIRONMENT=production
PORT=3001
ALLOW_ORIGINS=https://your-frontend-domain.com
```

### Frontend (.env in root)
```bash
VITE_API_BASE_URL=https://your-backend-domain.com
```

## Backend Deployment Options

### Option 1: Fly.io (Recommended)

1. Install Fly CLI:
```bash
curl -L https://fly.io/install.sh | sh
```

2. Login and create app:
```bash
fly auth login
cd server
fly launch
```

3. Set secrets:
```bash
fly secrets set DATABASE_URL="your-neon-connection-string"
fly secrets set JWT_SECRET="your-jwt-secret"
fly secrets set ALLOW_ORIGINS="https://your-frontend.com"
```

4. Create `fly.toml`:
```toml
app = "your-app-name"

[build]
  [build.args]
    GO_VERSION = "1.21"

[env]
  ENVIRONMENT = "production"
  PORT = "8080"

[[services]]
  internal_port = 8080
  protocol = "tcp"

  [[services.ports]]
    handlers = ["http"]
    port = 80

  [[services.ports]]
    handlers = ["tls", "http"]
    port = 443
```

5. Deploy:
```bash
fly deploy
```

### Option 2: Railway

1. Connect GitHub repository to Railway
2. Create new project → Deploy from GitHub
3. Set root directory to `server/`
4. Add environment variables:
   - DATABASE_URL
   - JWT_SECRET
   - ENVIRONMENT=production
   - ALLOW_ORIGINS
5. Railway auto-detects Go and deploys

### Option 3: Google Cloud Run

1. Build Docker image:
```dockerfile
# server/Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o server cmd/api/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

2. Deploy:
```bash
gcloud builds submit --tag gcr.io/PROJECT-ID/easy-count
gcloud run deploy easy-count \
  --image gcr.io/PROJECT-ID/easy-count \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

3. Set environment variables in Cloud Run console

### Option 4: Traditional VPS

1. Install Go on server:
```bash
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

2. Clone and build:
```bash
git clone your-repo
cd server
go build -o bin/server cmd/api/main.go
```

3. Create systemd service `/etc/systemd/system/easy-count.service`:
```ini
[Unit]
Description=Easy Count API Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/server
Environment="DATABASE_URL=your-connection-string"
Environment="JWT_SECRET=your-secret"
Environment="ENVIRONMENT=production"
Environment="PORT=3001"
Environment="ALLOW_ORIGINS=https://your-frontend.com"
ExecStart=/path/to/server/bin/server
Restart=always

[Install]
WantedBy=multi-user.target
```

4. Start service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable easy-count
sudo systemctl start easy-count
```

5. Configure Nginx reverse proxy:
```nginx
server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://localhost:3001;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Frontend Deployment

### Option 1: Vercel (Recommended)

1. Build frontend:
```bash
pnpm build
```

2. Deploy:
```bash
vercel
```

3. Set environment variable in Vercel dashboard:
   - VITE_API_BASE_URL: Your backend URL

### Option 2: Netlify

1. Build:
```bash
pnpm build
```

2. Deploy `dist/` directory to Netlify

3. Set environment variable:
   - VITE_API_BASE_URL: Your backend URL

### Option 3: GitHub Pages

1. Update `vite.config.ts`:
```typescript
export default defineConfig({
  base: '/easy-count/',
  // ...
})
```

2. Build and deploy:
```bash
pnpm build
# Deploy dist/ to gh-pages branch
```

## Database Setup (Neon)

1. Sign up at [neon.tech](https://neon.tech)
2. Create new project
3. Copy connection string
4. Database tables are auto-created on first backend start

## Security Checklist

Before production:
- [ ] Set strong JWT_SECRET (32+ characters, random)
- [ ] Configure ALLOW_ORIGINS to only your frontend domain
- [ ] Use HTTPS for both frontend and backend
- [ ] Enable Neon's IP allowlist (optional)
- [ ] Set ENVIRONMENT=production
- [ ] Test rate limiting
- [ ] Verify CORS configuration
- [ ] Review security headers
- [ ] Set up monitoring/logging

## Health Checks

Backend health endpoint:
```bash
curl https://your-backend.com/health
# Expected: {"status":"ok"}
```

Database connection test:
```bash
psql $DATABASE_URL -c "SELECT 1"
```

## Monitoring

### Recommended Services
- **Backend**: Fly.io metrics, Railway logs
- **Frontend**: Vercel Analytics
- **Errors**: Sentry (optional)
- **Uptime**: UptimeRobot

### Important Metrics
- API response times
- Error rates
- Database connection pool usage
- Authentication success/failure rates
- Rate limit triggers

## Scaling

### Backend
- Fly.io: Auto-scaling available
- Railway: Horizontal scaling
- VPS: Use load balancer (Nginx) with multiple instances

### Database
- Neon: Upgrade to higher tier for more compute/storage
- Connection pooling already implemented via pgx

## Backup Strategy

1. **Database**: Neon provides automatic backups
2. **Application**: Regular data exports via Backup button
3. **Code**: Git version control

## Cost Estimates

### Free Tier
- Neon: 0.5GB storage, 10 compute hours/month
- Fly.io: 3 shared-cpu instances, 3GB storage
- Vercel: Unlimited deployments

### Paid (estimated)
- Neon Pro: $19/month (10GB storage)
- Fly.io: ~$5-10/month (small instance)
- Vercel Pro: $20/month

Total: ~$25-50/month for small-medium usage

## Troubleshooting

### Backend not starting
- Check DATABASE_URL format
- Verify JWT_SECRET is set
- Check port availability
- Review logs for errors

### CORS errors
- Verify ALLOW_ORIGINS includes frontend URL
- Check protocol (http vs https)
- Ensure no trailing slashes

### Database connection issues
- Test connection: `psql $DATABASE_URL`
- Check Neon database is active
- Verify SSL mode in connection string

### Authentication not working
- Verify JWT_SECRET matches between deploys
- Check token expiration times
- Clear browser localStorage and retry

## Updates

To update deployed version:
```bash
# Backend
cd server
git pull
go build -o bin/server cmd/api/main.go
sudo systemctl restart easy-count

# Frontend
git pull
pnpm build
# Deploy new build
```

## Rollback

If deployment fails:
```bash
# Fly.io
fly releases
fly deploy --image <previous-image>

# Railway
# Use web UI to rollback

# VPS
git checkout <previous-commit>
# Rebuild and restart
```


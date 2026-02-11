# Deployment Guide

## Prerequisites

1. A Neon account (free tier available at [neon.tech](https://neon.tech))
2. Node.js 18+ installed
3. pnpm installed (`npm install -g pnpm`)

## Backend Deployment

### Option 1: Deploy to Vercel (Recommended)

The backend can be deployed as a serverless function on Vercel:

1. Install Vercel CLI:
```bash
npm install -g vercel
```

2. Create a `vercel.json` in the server directory:
```json
{
  "version": 2,
  "builds": [
    {
      "src": "src/index.ts",
      "use": "@vercel/node"
    }
  ],
  "routes": [
    {
      "src": "/(.*)",
      "dest": "src/index.ts"
    }
  ]
}
```

3. Deploy:
```bash
cd server
vercel
```

4. Set environment variables in Vercel dashboard:
   - `DATABASE_URL`: Your Neon connection string

### Option 2: Deploy to Railway

1. Create a new project on [Railway](https://railway.app)
2. Connect your GitHub repository
3. Set root directory to `server/`
4. Add environment variables:
   - `DATABASE_URL`: Your Neon connection string
5. Deploy

### Option 3: Traditional VPS (DigitalOcean, AWS, etc.)

1. SSH into your server
2. Clone the repository
3. Install dependencies:
```bash
cd server
npm install
npm run build
```

4. Create `.env` file with your configuration
5. Use PM2 to run the server:
```bash
npm install -g pm2
pm2 start dist/index.js --name easy-count-server
pm2 save
```

## Frontend Deployment

### Deploy to Vercel (Recommended)

1. Update `.env` with your production backend URL:
```bash
VITE_API_BASE_URL=https://your-backend-url.vercel.app
```

2. Deploy:
```bash
vercel
```

### Deploy to Netlify

1. Build the frontend:
```bash
pnpm build
```

2. Deploy the `dist/` directory to Netlify
3. Set environment variable `VITE_API_BASE_URL` in Netlify dashboard

### Deploy to GitHub Pages

1. Update `vite.config.ts` base path if needed
2. Build:
```bash
pnpm build
```

3. Deploy the `dist/` directory to GitHub Pages

## Environment Variables

### Backend (.env)
```
DATABASE_URL=postgresql://user:password@host/database?sslmode=require
PORT=3001
```

### Frontend (.env)
```
VITE_API_BASE_URL=https://your-backend-url.com
```

## Database Setup

The database tables are automatically created when the backend starts for the first time. No manual migration is needed.

## Monitoring

For production deployments, consider:
1. Setting up error logging (Sentry, LogRocket, etc.)
2. Monitoring database performance in Neon dashboard
3. Setting up uptime monitoring
4. Regular database backups (Neon provides automatic backups)

## Scaling Considerations

- Neon's free tier includes 0.5GB storage and 10 compute hours per month
- For higher traffic, upgrade to Neon's paid plans
- Consider implementing caching for frequently accessed data
- Add rate limiting to API endpoints for production use

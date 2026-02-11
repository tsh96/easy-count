# Migration Guide: IndexedDB to PostgreSQL with Authentication

This guide helps you migrate from the old IndexedDB version to the new PostgreSQL + Authentication version.

## What's New

The new version includes:
- ✅ Cloud-based PostgreSQL storage (via Neon)
- ✅ User authentication with JWT
- ✅ Production-ready security features
- ✅ Multi-user support with data isolation
- ✅ Golang + Gin backend for high performance

## Migration Steps

### Step 1: Export Your Old Data

Before migrating, backup your data from the old version:

1. Open the old version of Easy Count
2. **Personal Records**: Click Backup → Save JSON file
3. **Customer Records**: Click Backup → Save JSON file

**Important**: Keep these backup files safe!

### Step 2: Set Up New Backend

1. **Get Neon Database**:
   - Sign up at [neon.tech](https://neon.tech) (free tier available)
   - Create new project
   - Copy connection string

2. **Configure Backend**:
   ```bash
   cd server
   cp .env.example .env
   # Edit .env:
   # - Add your DATABASE_URL
   # - Generate JWT_SECRET: openssl rand -base64 32
   # - Set ALLOW_ORIGINS to your frontend URL
   ```

3. **Start Backend**:
   ```bash
   go mod download
   go run cmd/api/main.go
   ```

### Step 3: Set Up Frontend

1. **Configure Frontend**:
   ```bash
   cd ..
   cp .env.example .env
   # VITE_API_BASE_URL should be http://localhost:3001
   ```

2. **Start Frontend**:
   ```bash
   pnpm install
   pnpm dev
   ```

### Step 4: Create Account

1. Open `http://localhost:5173`
2. Click "Create Account"
3. Enter email and password (minimum 8 characters)
4. Register

### Step 5: Restore Your Data

1. After login, navigate to Personal Records
2. Click Restore
3. Select your personal records JSON backup
4. Wait for completion

5. Navigate to Customer Records
6. Click Restore
7. Select your customer records JSON backup
8. Wait for completion

### Step 6: Verify

- Check all personal transactions are present
- Check all customer records (Private and Government)
- Verify calculations and totals

## Key Differences

### Authentication Required
- Old: No login required (local storage)
- New: Email/password authentication required

### Data Storage
- Old: Browser IndexedDB (local only)
- New: PostgreSQL cloud database (accessible anywhere)

### Multi-User Support
- Old: Single user per browser
- New: Multiple users with isolated data

### Security
- Old: No authentication
- New: JWT tokens, password hashing, rate limiting, etc.

## Troubleshooting

### "Failed to fetch" Error
- Ensure backend server is running on port 3001
- Check VITE_API_BASE_URL in frontend .env

### Authentication Errors
- Verify JWT_SECRET is set in backend .env
- Try logging out and back in
- Clear browser localStorage

### Data Not Appearing
- Check browser console for errors (F12 → Console)
- Verify backup JSON files are not corrupted
- Ensure you're logged in

### Database Connection Error
- Verify DATABASE_URL in server/.env is correct
- Check Neon database is active
- Test connection: `psql $DATABASE_URL`

## Automatic Migration

The app still includes automatic migration for very old data:
- Old localStorage-based transactions
- Old localStorage-based customer records

These migrate automatically on first page visit after login.

## Backup/Restore Format

The backup JSON format is compatible between versions:
- Export from old version → Import to new version ✅
- Export from new version → Keep as backup ✅

## Multiple Devices

With the new version, you can:
1. Register once on any device
2. Login from any device with same credentials
3. Access your data from anywhere

All data is synced to PostgreSQL in real-time.

## Security Best Practices

For production deployment:
1. Use strong, unique password
2. Set strong JWT_SECRET (32+ characters)
3. Enable HTTPS
4. Use secure DATABASE_URL
5. Regularly backup your data

## Need Help?

- Check [GitHub Issues](https://github.com/tsh96/easy-count/issues)
- Review server logs for errors
- Verify environment variables are correct
- Test backend health: `curl http://localhost:3001/health`


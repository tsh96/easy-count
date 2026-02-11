# Migration Guide: IndexedDB to PostgreSQL

This guide helps you migrate your existing data from the old IndexedDB version to the new PostgreSQL-based version.

## Overview

The application has been upgraded from using IndexedDB (local browser storage) to PostgreSQL (Neon cloud database). This provides:
- ✅ Data persistence across devices
- ✅ Better performance for large datasets
- ✅ Cloud backup and recovery
- ✅ Multi-device access (future feature)

## Migration Steps

### Step 1: Backup Your Existing Data

Before migrating, backup your existing data from the old version:

1. Open the old version of Easy Count
2. Navigate to Personal Records page
3. Click the **Backup** button
4. Save the JSON file (e.g., `personal_records_backup.json`)
5. Navigate to Customer Records page
6. Click the **Backup** button
7. Save the JSON file (e.g., `customer_records_backup.json`)

**Important**: Keep these backup files safe! They contain all your data.

### Step 2: Set Up the New Version

1. **Set up Neon Database:**
   - Sign up for a free account at [neon.tech](https://neon.tech)
   - Create a new project
   - Copy the connection string (it looks like: `postgresql://user:password@host/database?sslmode=require`)

2. **Configure Backend:**
   ```bash
   cd server
   cp .env.example .env
   # Edit .env and paste your Neon connection string
   npm install
   npm run dev
   ```

   The server should start and automatically create the database tables.

3. **Configure Frontend:**
   ```bash
   cd ..
   cp .env.example .env
   # The default API URL (http://localhost:3001) should work for local development
   pnpm install
   pnpm dev
   ```

### Step 3: Restore Your Data

1. Open the new version at `http://localhost:5173`
2. Navigate to Personal Records page
3. Click the **Restore** button
4. Select your `personal_records_backup.json` file
5. Wait for the restoration to complete
6. Navigate to Customer Records page
7. Click the **Restore** button
8. Select your `customer_records_backup.json` file
9. Wait for the restoration to complete

### Step 4: Verify Your Data

1. Check that all your personal transactions are displayed correctly
2. Check that all your customer records (both Private and Government) are displayed correctly
3. Verify the totals and calculations are correct

## Troubleshooting

### Issue: "Failed to fetch" error

**Solution**: Make sure the backend server is running on port 3001.

```bash
cd server
npm run dev
```

### Issue: Data not appearing after restore

**Solution**: 
1. Check browser console for errors (F12 → Console tab)
2. Verify the backup JSON file is not corrupted
3. Restart both backend and frontend servers

### Issue: Database connection error

**Solution**:
1. Verify your `DATABASE_URL` in `server/.env` is correct
2. Check that your Neon database is active
3. Ensure your IP is allowed in Neon's settings (Neon allows all IPs by default)

## What Happens to Old Data?

After successful migration:
- Your old IndexedDB data remains in your browser but is no longer used
- The new version reads from PostgreSQL database
- You can clear browser data safely after verifying migration was successful

## Automatic Old Data Migration

The application includes automatic migration for:
- Old localStorage-based transactions
- Old localStorage-based customer records

These migrations run once when you first open the new version. If you had data in the very old format, it will be automatically migrated to the database when you open each page for the first time.

## Data Format Compatibility

The backup/restore JSON format remains compatible between versions, so you can:
- Restore old backups into the new version
- Create backups in the new version
- Share backups between users

## Need Help?

If you encounter issues during migration:
1. Check the [GitHub Issues](https://github.com/tsh96/easy-count/issues)
2. Create a new issue with details about your problem
3. Include any error messages from the browser console

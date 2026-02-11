# Easy Count

An accounting application for tracking personal transactions and customer records (invoices and cheques).

## Architecture

This application now uses a client-server architecture:
- **Frontend**: Vue 3 application (in root directory)
- **Backend**: Express.js server with PostgreSQL/Neon database (in `server/` directory)

## Setup

### Backend Setup

1. Navigate to the server directory:
```bash
cd server
```

2. Install dependencies:
```bash
npm install
```

3. Create a `.env` file:
```bash
cp .env.example .env
```

4. Configure your Neon database:
   - Sign up at [Neon](https://neon.tech)
   - Create a new project
   - Copy the connection string to `.env` as `DATABASE_URL`

5. Start the backend server:
```bash
npm run dev
```

The server will run on `http://localhost:3001` by default.

### Frontend Setup

1. Return to the root directory:
```bash
cd ..
```

2. Install dependencies:
```bash
pnpm install
```

3. Create a `.env` file:
```bash
cp .env.example .env
```

4. Start the development server:
```bash
pnpm dev
```

The application will be available at `http://localhost:5173`.

## Features

- **Personal Records**: Track income and expenses with automatic balance calculation
- **Customer Records**: Manage invoice and cheque tracking for private and government customers
- **Backup/Restore**: Export and import your data as JSON files
- **Data Persistence**: All data is stored in PostgreSQL database via Neon

## Migration from IndexedDB

If you're migrating from the old IndexedDB version:
1. Export your data using the Backup button in the old version
2. Set up the new backend server with Neon
3. Import your data using the Restore button in the new version

## Special Thanks
<a href="https://www.flaticon.com/free-icons/accounting" title="accounting icons">Accounting icons created by Freepik - Flaticon</a>
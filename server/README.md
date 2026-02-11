# Easy Count Server

Backend API server for Easy Count application using PostgreSQL (Neon) for data storage.

## Setup

1. Install dependencies:
```bash
npm install
```

2. Create a `.env` file based on `.env.example`:
```bash
cp .env.example .env
```

3. Configure your Neon database connection:
   - Sign up for a free account at [Neon](https://neon.tech)
   - Create a new project
   - Copy the connection string to your `.env` file as `DATABASE_URL`

4. Start the development server:
```bash
npm run dev
```

The server will automatically create the necessary database tables on first run.

## API Endpoints

### Personal Records (Transactions)

- `GET /api/transactions` - Get all transactions
- `POST /api/transactions` - Create a new transaction
- `PUT /api/transactions/:id` - Update a transaction
- `DELETE /api/transactions/:id` - Delete a transaction
- `GET /api/backup/transactions` - Backup all transactions
- `POST /api/restore/transactions` - Restore transactions from backup

### Customer Records

- `GET /api/customer-records/:type` - Get customer records (type: Private or Government)
- `POST /api/customer-records/:type` - Create a new customer record
- `POST /api/customer-records/:type/bulk` - Create multiple customer records
- `PUT /api/customer-records/:type/:id` - Update a customer record
- `DELETE /api/customer-records/:type/:id` - Delete a customer record
- `POST /api/customer-records/:type/replace-name` - Replace customer name
- `GET /api/backup/customer-records` - Backup all customer records
- `POST /api/restore/customer-records` - Restore customer records from backup

## Production Build

Build the server for production:
```bash
npm run build
```

Run the production build:
```bash
npm start
```

## Environment Variables

- `DATABASE_URL` - Neon PostgreSQL connection string (required)
- `PORT` - Server port (default: 3001)

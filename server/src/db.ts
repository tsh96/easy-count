import { neon } from '@neondatabase/serverless';
import dotenv from 'dotenv';

dotenv.config();

export const sql = neon(process.env.DATABASE_URL!);

export async function initDatabase() {
  // Create transactions table for personal records
  await sql`
    CREATE TABLE IF NOT EXISTS transactions (
      id SERIAL PRIMARY KEY,
      date BIGINT NOT NULL,
      description TEXT NOT NULL DEFAULT '',
      credit DECIMAL(10, 2) NOT NULL DEFAULT 0,
      debit DECIMAL(10, 2) NOT NULL DEFAULT 0,
      user_id TEXT NOT NULL DEFAULT 'default',
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
  `;

  // Create index for date ordering
  await sql`
    CREATE INDEX IF NOT EXISTS idx_transactions_date_id ON transactions(date, id)
  `;

  // Create index for user_id
  await sql`
    CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id)
  `;

  // Create customer_records table (both private and government)
  await sql`
    CREATE TABLE IF NOT EXISTS customer_records (
      id SERIAL PRIMARY KEY,
      record_type TEXT NOT NULL CHECK (record_type IN ('Private', 'Government')),
      invoice_date BIGINT,
      invoice_no TEXT NOT NULL DEFAULT '',
      customer_name TEXT NOT NULL DEFAULT '',
      invoice_amount DECIMAL(10, 2),
      cheque_date BIGINT,
      cheque_no TEXT NOT NULL DEFAULT '',
      cheque_amount DECIMAL(10, 2),
      remark TEXT NOT NULL DEFAULT '',
      user_id TEXT NOT NULL DEFAULT 'default',
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
  `;

  // Create indexes for customer_records
  await sql`
    CREATE INDEX IF NOT EXISTS idx_customer_records_invoice_date_no ON customer_records(invoice_date, invoice_no)
  `;

  await sql`
    CREATE INDEX IF NOT EXISTS idx_customer_records_type_user ON customer_records(record_type, user_id)
  `;

  await sql`
    CREATE INDEX IF NOT EXISTS idx_customer_records_invoice_date ON customer_records(invoice_date)
  `;

  console.log('Database initialized successfully');
}

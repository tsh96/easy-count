package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Database wraps the database connection
type Database struct {
	Pool *pgxpool.Pool
}

// NewDatabase creates a new database connection
func NewDatabase(databaseURL string) (*Database, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &Database{Pool: pool}, nil
}

// Close closes the database connection
func (db *Database) Close() {
	db.Pool.Close()
}

// InitSchema initializes the database schema
func (db *Database) InitSchema() error {
	ctx := context.Background()

	// Create users table
	_, err := db.Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create index on email
	_, err = db.Pool.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users email index: %w", err)
	}

	// Create transactions table
	_, err = db.Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			date BIGINT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			credit DECIMAL(10, 2) NOT NULL DEFAULT 0,
			debit DECIMAL(10, 2) NOT NULL DEFAULT 0,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}

	// Create indexes for transactions
	_, err = db.Pool.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_transactions_date_id ON transactions(date, id)
	`)
	if err != nil {
		return fmt.Errorf("failed to create transactions date index: %w", err)
	}

	_, err = db.Pool.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id)
	`)
	if err != nil {
		return fmt.Errorf("failed to create transactions user_id index: %w", err)
	}

	// Create customer_records table
	_, err = db.Pool.Exec(ctx, `
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
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create customer_records table: %w", err)
	}

	// Create indexes for customer_records
	_, err = db.Pool.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_customer_records_invoice_date_no ON customer_records(invoice_date, invoice_no)
	`)
	if err != nil {
		return fmt.Errorf("failed to create customer_records invoice index: %w", err)
	}

	_, err = db.Pool.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_customer_records_type_user ON customer_records(record_type, user_id)
	`)
	if err != nil {
		return fmt.Errorf("failed to create customer_records type_user index: %w", err)
	}

	_, err = db.Pool.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_customer_records_invoice_date ON customer_records(invoice_date)
	`)
	if err != nil {
		return fmt.Errorf("failed to create customer_records invoice_date index: %w", err)
	}

	_, err = db.Pool.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_customer_records_user_id ON customer_records(user_id)
	`)
	if err != nil {
		return fmt.Errorf("failed to create customer_records user_id index: %w", err)
	}

	return nil
}

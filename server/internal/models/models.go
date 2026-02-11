package models

import "time"

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never send password hash to client
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Transaction represents a financial transaction
type Transaction struct {
	ID          int     `json:"id,omitempty"`
	Date        int64   `json:"date"`
	Description string  `json:"description"`
	Credit      float64 `json:"credit"`
	Debit       float64 `json:"debit"`
	UserID      int     `json:"-"`
}

// CustomerRecord represents a customer invoice/cheque record
type CustomerRecord struct {
	ID            int      `json:"id,omitempty"`
	InvoiceDate   *int64   `json:"invoiceDate,omitempty"`
	InvoiceNo     string   `json:"invoiceNo"`
	CustomerName  string   `json:"customerName"`
	InvoiceAmount *float64 `json:"invoiceAmount,omitempty"`
	ChequeDate    *int64   `json:"chequeDate,omitempty"`
	ChequeNo      string   `json:"chequeNo"`
	ChequeAmount  *float64 `json:"chequeAmount,omitempty"`
	Remark        string   `json:"remark"`
	UserID        int      `json:"-"`
}

// AuthRequest represents a login/register request
type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	User         User   `json:"user"`
}

// RefreshRequest represents a refresh token request
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// BackupTransactions represents backup data for transactions
type BackupTransactions struct {
	Transactions []Transaction `json:"transactions"`
}

// BackupCustomerRecords represents backup data for customer records
type BackupCustomerRecords struct {
	PrivateRecords    []CustomerRecord `json:"privateCustomerRecords"`
	GovernmentRecords []CustomerRecord `json:"governmentCustomerRecords"`
}

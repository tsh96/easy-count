package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tsh96/easy-count/server/internal/middleware"
	"github.com/tsh96/easy-count/server/internal/models"
)

// BackupTransactions backs up all transactions for a user
func (h *Handlers) BackupTransactions(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	rows, err := h.db.Pool.Query(ctx, `
		SELECT date, description, credit, debit
		FROM transactions
		WHERE user_id = $1
		ORDER BY date, id
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}
	defer rows.Close()

	transactions := []models.Transaction{}
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.Date, &t.Description, &t.Credit, &t.Debit); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan transaction"})
			return
		}
		transactions = append(transactions, t)
	}

	c.JSON(http.StatusOK, models.BackupTransactions{Transactions: transactions})
}

// RestoreTransactions restores transactions from backup
func (h *Handlers) RestoreTransactions(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var backup models.BackupTransactions
	if err := c.ShouldBindJSON(&backup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	tx, err := h.db.Pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	// Delete existing transactions
	_, err = tx.Exec(ctx, `DELETE FROM transactions WHERE user_id = $1`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear transactions"})
		return
	}

	// Insert new transactions
	for _, t := range backup.Transactions {
		_, err := tx.Exec(ctx, `
			INSERT INTO transactions (date, description, credit, debit, user_id)
			VALUES ($1, $2, $3, $4, $5)
		`, t.Date, t.Description, t.Credit, t.Debit, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore transaction"})
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// BackupCustomerRecords backs up all customer records for a user
func (h *Handlers) BackupCustomerRecords(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()

	// Get private records
	privateRows, err := h.db.Pool.Query(ctx, `
		SELECT invoice_date, invoice_no, customer_name, invoice_amount,
			   cheque_date, cheque_no, cheque_amount, remark
		FROM customer_records
		WHERE record_type = 'Private' AND user_id = $1
		ORDER BY invoice_date, invoice_no
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch private records"})
		return
	}
	defer privateRows.Close()

	privateRecords := []models.CustomerRecord{}
	for privateRows.Next() {
		var r models.CustomerRecord
		if err := privateRows.Scan(&r.InvoiceDate, &r.InvoiceNo, &r.CustomerName,
			&r.InvoiceAmount, &r.ChequeDate, &r.ChequeNo, &r.ChequeAmount, &r.Remark); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan private record"})
			return
		}
		privateRecords = append(privateRecords, r)
	}

	// Get government records
	govRows, err := h.db.Pool.Query(ctx, `
		SELECT invoice_date, invoice_no, customer_name, invoice_amount,
			   cheque_date, cheque_no, cheque_amount, remark
		FROM customer_records
		WHERE record_type = 'Government' AND user_id = $1
		ORDER BY invoice_date, invoice_no
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch government records"})
		return
	}
	defer govRows.Close()

	govRecords := []models.CustomerRecord{}
	for govRows.Next() {
		var r models.CustomerRecord
		if err := govRows.Scan(&r.InvoiceDate, &r.InvoiceNo, &r.CustomerName,
			&r.InvoiceAmount, &r.ChequeDate, &r.ChequeNo, &r.ChequeAmount, &r.Remark); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan government record"})
			return
		}
		govRecords = append(govRecords, r)
	}

	c.JSON(http.StatusOK, models.BackupCustomerRecords{
		PrivateRecords:    privateRecords,
		GovernmentRecords: govRecords,
	})
}

// RestoreCustomerRecords restores customer records from backup
func (h *Handlers) RestoreCustomerRecords(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var backup models.BackupCustomerRecords
	if err := c.ShouldBindJSON(&backup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	tx, err := h.db.Pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	// Delete existing records
	_, err = tx.Exec(ctx, `DELETE FROM customer_records WHERE user_id = $1`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear customer records"})
		return
	}

	// Insert private records
	for _, r := range backup.PrivateRecords {
		_, err := tx.Exec(ctx, `
			INSERT INTO customer_records 
				(record_type, invoice_date, invoice_no, customer_name, invoice_amount,
				 cheque_date, cheque_no, cheque_amount, remark, user_id)
			VALUES ('Private', $1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, r.InvoiceDate, r.InvoiceNo, r.CustomerName, r.InvoiceAmount,
			r.ChequeDate, r.ChequeNo, r.ChequeAmount, r.Remark, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore private record"})
			return
		}
	}

	// Insert government records
	for _, r := range backup.GovernmentRecords {
		_, err := tx.Exec(ctx, `
			INSERT INTO customer_records 
				(record_type, invoice_date, invoice_no, customer_name, invoice_amount,
				 cheque_date, cheque_no, cheque_amount, remark, user_id)
			VALUES ('Government', $1, $2, $3, $4, $5, $6, $7, $8, $9)
		`, r.InvoiceDate, r.InvoiceNo, r.CustomerName, r.InvoiceAmount,
			r.ChequeDate, r.ChequeNo, r.ChequeAmount, r.Remark, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore government record"})
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

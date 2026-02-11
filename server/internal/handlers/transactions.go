package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/tsh96/easy-count/server/internal/middleware"
	"github.com/tsh96/easy-count/server/internal/models"
)

// GetTransactions retrieves all transactions for a user
func (h *Handlers) GetTransactions(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	rows, err := h.db.Pool.Query(ctx, `
		SELECT id, date, description, credit, debit
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
		if err := rows.Scan(&t.ID, &t.Date, &t.Description, &t.Credit, &t.Debit); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan transaction"})
			return
		}
		transactions = append(transactions, t)
	}

	c.JSON(http.StatusOK, transactions)
}

// CreateTransaction creates a new transaction
func (h *Handlers) CreateTransaction(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var t models.Transaction
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	err := h.db.Pool.QueryRow(ctx, `
		INSERT INTO transactions (date, description, credit, debit, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, date, description, credit, debit
	`, t.Date, t.Description, t.Credit, t.Debit, userID).Scan(
		&t.ID, &t.Date, &t.Description, &t.Credit, &t.Debit,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, t)
}

// UpdateTransaction updates an existing transaction
func (h *Handlers) UpdateTransaction(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var t models.Transaction
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	err = h.db.Pool.QueryRow(ctx, `
		UPDATE transactions
		SET date = $1, description = $2, credit = $3, debit = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5 AND user_id = $6
		RETURNING id, date, description, credit, debit
	`, t.Date, t.Description, t.Credit, t.Debit, id, userID).Scan(
		&t.ID, &t.Date, &t.Description, &t.Credit, &t.Debit,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	c.JSON(http.StatusOK, t)
}

// DeleteTransaction deletes a transaction
func (h *Handlers) DeleteTransaction(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	ctx := context.Background()
	result, err := h.db.Pool.Exec(ctx, `
		DELETE FROM transactions
		WHERE id = $1 AND user_id = $2
	`, id, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// BulkCreateTransactions creates multiple transactions
func (h *Handlers) BulkCreateTransactions(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Transactions []models.Transaction `json:"transactions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
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

	results := []models.Transaction{}
	for _, t := range req.Transactions {
		var result models.Transaction
		err := tx.QueryRow(ctx, `
			INSERT INTO transactions (date, description, credit, debit, user_id)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, t.Date, t.Description, t.Credit, t.Debit, userID).Scan(&result.ID)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
			return
		}
		results = append(results, result)
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, results)
}

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

// GetCustomerRecords retrieves customer records by type
func (h *Handlers) GetCustomerRecords(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	recordType := c.Param("type")
	if recordType != "Private" && recordType != "Government" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record type"})
		return
	}

	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	ctx := context.Background()
	var rows pgx.Rows
	var err error

	if startDate != "" && endDate != "" {
		rows, err = h.db.Pool.Query(ctx, `
			SELECT id, invoice_date, invoice_no, customer_name, invoice_amount,
				   cheque_date, cheque_no, cheque_amount, remark
			FROM customer_records
			WHERE record_type = $1 AND user_id = $2
			  AND invoice_date >= $3 AND invoice_date < $4
			ORDER BY invoice_date, invoice_no
		`, recordType, userID, startDate, endDate)
	} else {
		rows, err = h.db.Pool.Query(ctx, `
			SELECT id, invoice_date, invoice_no, customer_name, invoice_amount,
				   cheque_date, cheque_no, cheque_amount, remark
			FROM customer_records
			WHERE record_type = $1 AND user_id = $2
			ORDER BY invoice_date, invoice_no
		`, recordType, userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customer records"})
		return
	}
	defer rows.Close()

	records := []models.CustomerRecord{}
	for rows.Next() {
		var r models.CustomerRecord
		if err := rows.Scan(&r.ID, &r.InvoiceDate, &r.InvoiceNo, &r.CustomerName,
			&r.InvoiceAmount, &r.ChequeDate, &r.ChequeNo, &r.ChequeAmount, &r.Remark); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan customer record"})
			return
		}
		records = append(records, r)
	}

	c.JSON(http.StatusOK, records)
}

// CreateCustomerRecord creates a new customer record
func (h *Handlers) CreateCustomerRecord(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	recordType := c.Param("type")
	if recordType != "Private" && recordType != "Government" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record type"})
		return
	}

	var r models.CustomerRecord
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	err := h.db.Pool.QueryRow(ctx, `
		INSERT INTO customer_records 
			(record_type, invoice_date, invoice_no, customer_name, invoice_amount,
			 cheque_date, cheque_no, cheque_amount, remark, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, invoice_date, invoice_no, customer_name, invoice_amount,
				  cheque_date, cheque_no, cheque_amount, remark
	`, recordType, r.InvoiceDate, r.InvoiceNo, r.CustomerName, r.InvoiceAmount,
		r.ChequeDate, r.ChequeNo, r.ChequeAmount, r.Remark, userID).Scan(
		&r.ID, &r.InvoiceDate, &r.InvoiceNo, &r.CustomerName, &r.InvoiceAmount,
		&r.ChequeDate, &r.ChequeNo, &r.ChequeAmount, &r.Remark,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer record"})
		return
	}

	c.JSON(http.StatusCreated, r)
}

// UpdateCustomerRecord updates a customer record
func (h *Handlers) UpdateCustomerRecord(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	recordType := c.Param("type")
	if recordType != "Private" && recordType != "Government" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record type"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	var r models.CustomerRecord
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	err = h.db.Pool.QueryRow(ctx, `
		UPDATE customer_records
		SET invoice_date = $1, invoice_no = $2, customer_name = $3, invoice_amount = $4,
			cheque_date = $5, cheque_no = $6, cheque_amount = $7, remark = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $9 AND record_type = $10 AND user_id = $11
		RETURNING id, invoice_date, invoice_no, customer_name, invoice_amount,
				  cheque_date, cheque_no, cheque_amount, remark
	`, r.InvoiceDate, r.InvoiceNo, r.CustomerName, r.InvoiceAmount,
		r.ChequeDate, r.ChequeNo, r.ChequeAmount, r.Remark, id, recordType, userID).Scan(
		&r.ID, &r.InvoiceDate, &r.InvoiceNo, &r.CustomerName, &r.InvoiceAmount,
		&r.ChequeDate, &r.ChequeNo, &r.ChequeAmount, &r.Remark,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer record"})
		return
	}

	c.JSON(http.StatusOK, r)
}

// DeleteCustomerRecord deletes a customer record
func (h *Handlers) DeleteCustomerRecord(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	recordType := c.Param("type")
	if recordType != "Private" && recordType != "Government" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record type"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record ID"})
		return
	}

	ctx := context.Background()
	result, err := h.db.Pool.Exec(ctx, `
		DELETE FROM customer_records
		WHERE id = $1 AND record_type = $2 AND user_id = $3
	`, id, recordType, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer record"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer record not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// BulkCreateCustomerRecords creates multiple customer records
func (h *Handlers) BulkCreateCustomerRecords(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	recordType := c.Param("type")
	if recordType != "Private" && recordType != "Government" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record type"})
		return
	}

	var req struct {
		Records []models.CustomerRecord `json:"records"`
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

	results := []models.CustomerRecord{}
	for _, r := range req.Records {
		var result models.CustomerRecord
		err := tx.QueryRow(ctx, `
			INSERT INTO customer_records 
				(record_type, invoice_date, invoice_no, customer_name, invoice_amount,
				 cheque_date, cheque_no, cheque_amount, remark, user_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id
		`, recordType, r.InvoiceDate, r.InvoiceNo, r.CustomerName, r.InvoiceAmount,
			r.ChequeDate, r.ChequeNo, r.ChequeAmount, r.Remark, userID).Scan(&result.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer record"})
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

// ReplaceCustomerName replaces customer name in records
func (h *Handlers) ReplaceCustomerName(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	recordType := c.Param("type")
	if recordType != "Private" && recordType != "Government" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid record type"})
		return
	}

	var req struct {
		OldName string `json:"oldName" binding:"required"`
		NewName string `json:"newName" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	_, err := h.db.Pool.Exec(ctx, `
		UPDATE customer_records
		SET customer_name = $1, updated_at = CURRENT_TIMESTAMP
		WHERE customer_name = $2 AND record_type = $3 AND user_id = $4
	`, req.NewName, req.OldName, recordType, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to replace customer name"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

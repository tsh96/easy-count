package handlers

import (
	"github.com/tsh96/easy-count/server/internal/config"
	"github.com/tsh96/easy-count/server/internal/db"
)

// Handlers holds all handler dependencies
type Handlers struct {
	db  *db.Database
	cfg *config.Config
}

// NewHandlers creates a new Handlers instance
func NewHandlers(database *db.Database, cfg *config.Config) *Handlers {
	return &Handlers{
		db:  database,
		cfg: cfg,
	}
}

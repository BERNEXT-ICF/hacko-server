package repository

import (
	"hacko-app/internal/module/z_v2/ports"

	"github.com/jmoiron/sqlx"
)

var _ ports.XxxRepository = &xxxRepository{}

type xxxRepository struct {
	db *sqlx.DB
}

func NewXxxRepository(db *sqlx.DB) *xxxRepository {
	return &xxxRepository{
		db: db,
	}
}

package model

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	Id        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Number    string    `db:"number"`
	Status    string    `db:"status"`
	Accrual   float64   `db:"accrual"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

const ORDER_STATUS_NEW = "NEW"

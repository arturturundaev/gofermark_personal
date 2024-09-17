package model

import (
	"github.com/google/uuid"
	"time"
)

type UserWithdrawals struct {
	Id        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Sum       float64   `db:"sum"`
	Number    string    `db:"number"`
	CreatedAt time.Time `db:"created_at"`
}

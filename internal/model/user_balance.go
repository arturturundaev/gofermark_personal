package model

import "github.com/google/uuid"

type UserBalance struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Sum       float64   `db:"sum" json:"Current"`
	WithDrawn float64   `db:"with_drawn"`
}

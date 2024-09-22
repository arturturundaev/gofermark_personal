package model

import "github.com/google/uuid"

type UserBalance struct {
	Id        uuid.UUID `db:"id"`
	UserId    uuid.UUID `db:"user_id"`
	Sum       float64   `db:"sum" json:"Current"`
	WithDrawn float64   `db:"with_drawn"`
}

package model

import "github.com/google/uuid"

type UserBalance struct {
	id        uuid.UUID `db:"id"`
	UserId    uuid.UUID `db:"user_id"`
	sum       float64   `db:"sum"`
	withDrawn float64   `db:"with_drawn"`
}

package model

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `json:"id" db:"id,omitempty"`
	Login    string    `json:"login,omitempty" db:"login,omitempty"`
	Password string    `json:"password,omitempty" db:"password,omitempty"`
}

func NewUserFromGin(ctx *gin.Context) (*User, error) {
	var user User

	if err := ctx.BindJSON(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

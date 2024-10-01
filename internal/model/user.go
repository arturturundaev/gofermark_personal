package model

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id" db:"id,omitempty"`
	Login    string    `json:"login,omitempty" validate:"required,min=4" db:"login,omitempty"`
	Password string    `json:"password,omitempty" validate:"required,min=4" db:"password,omitempty"`
}

func NewUserFromGin(ctx *gin.Context) (*User, error) {
	var user User

	if err := ctx.ShouldBindWith(&user, binding.JSON); err != nil {
		return nil, err
	}

	validate := validator.New()
	err := validate.Struct(user)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return nil, validationErrors
	}

	return &user, nil
}

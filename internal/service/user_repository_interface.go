package service

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gofermark_personal/internal/model"
)

type IUserRepository interface {
	Save(id uuid.UUID, login string, password string) error
	GetDB() *sqlx.DB
	UserExistsByLogin(login string) (bool, error)
	UserExistsByID(id uuid.UUID) (bool, error)
	GetByLogin(login string) (*model.User, error)

	GetBalance(userID uuid.UUID) (*model.UserBalance, error)
	Withdraw(userID uuid.UUID, number string, sum float64) error
	GetWithdrawals(userID uuid.UUID) ([]model.UserWithdrawals, error)
}

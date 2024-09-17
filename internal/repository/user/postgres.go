package user

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"gofermark_personal/internal/model"
)

const tableName = "users"

type UserRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewUserRepository(dns string, logger *zap.Logger) (*UserRepository, error) {
	database, err := sqlx.Open("postgres", dns)
	if err != nil {
		return nil, err
	}

	return &UserRepository{db: database, logger: logger}, nil
}

func (repository *UserRepository) Save(id uuid.UUID, login string, password string) error {
	_, err := repository.db.Exec(
		fmt.Sprintf(`INSERT INTO %s (id, login, password) VALUES ($1, $2, $3)`, tableName), id, login, password,
	)

	return err
}

func (repository *UserRepository) UserExistsByLogin(login string) (bool, error) {
	return repository.userExistsByProperty("login", login)
}

func (repository *UserRepository) UserExistsByID(id uuid.UUID) (bool, error) {
	return repository.userExistsByProperty("id", id)
}

func (repository *UserRepository) userExistsByProperty(propertyName string, value interface{}) (bool, error) {
	var exists bool

	err := repository.db.Get(
		&exists,
		fmt.Sprintf(`SELECT true FROM %s WHERE %s = $1`, tableName, propertyName), value,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (repository *UserRepository) GetByLogin(login string) (*model.User, error) {
	var user model.User

	err := repository.db.Get(
		&user,
		fmt.Sprintf(`SELECT id, login, password FROM %s WHERE login = $1`, tableName), login,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepository) GetDB() *sqlx.DB {
	return repository.db
}

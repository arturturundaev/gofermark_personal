package user

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	errors2 "gofermark_personal/internal/errors"
	"gofermark_personal/internal/model"
	"time"
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

func (repository *UserRepository) GetBalance(userID uuid.UUID) (*model.UserBalance, error) {
	var result model.UserBalance
	err := repository.db.Get(
		&result,
		"SELECT id, user_id, sum, with_drawn FROM user_balance WHERE user_id=$1", userID)
	if err != nil {
		repository.logger.Error("Failed to get balance", zap.String("error", err.Error()))
		return &result, err
	}
	return &result, nil
}

func (repository *UserRepository) Withdraw(userID uuid.UUID, number string, sum float64) error {
	tx, err := repository.db.Begin()
	if err != nil {
		repository.logger.Error("Failed to create transaction", zap.String("error", err.Error()))
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE user_balance SET sum=sum-$1, with_drawn=with_drawn+$1 WHERE user_id=$2", sum, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return errors2.ErrNotEnoughtMoney
		}
		repository.logger.Error("Failed to update balance", zap.String("error", err.Error()))
		return err
	}

	_, err = tx.Exec("INSERT INTO user_withdrawals (id, user_id, number, sum, created_at) Values ($1, $2, $3, $4, $5)", uuid.New(), userID, number, sum, time.Now())
	if err != nil {
		repository.logger.Error("Failed to insert into withdrawals", zap.String("error", err.Error()))
		return err
	}

	return tx.Commit()
}

func (repository *UserRepository) GetWithdrawals(userID uuid.UUID) ([]model.UserWithdrawals, error) {
	var result []model.UserWithdrawals

	err := repository.db.Select(
		&result,
		`SELECT id, user_id, sum, number, created_at FROM user_withdrawals WHERE user_id=$1`,
		userID)
	if err != nil {
		repository.logger.Error("Failed to get withdrawals", zap.String("error", err.Error()))
		return result, err
	}

	return result, nil
}

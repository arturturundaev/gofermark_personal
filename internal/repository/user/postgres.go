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
	"gofermark_personal/internal/handler/user"
	"gofermark_personal/internal/model"
	"time"
)

const tableName = "users"

type UserRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewUserRepository(db *sqlx.DB, logger *zap.Logger) *UserRepository {
	return &UserRepository{db: db, logger: logger}
}

var createUserSql = `INSERT INTO %s (id, login, password) VALUES (:id, :login, :password)`

func (repository *UserRepository) Save(id uuid.UUID, login string, password string) error {
	_, err := repository.db.NamedExec(fmt.Sprintf(createUserSql, tableName), map[string]interface{}{"id": id, "login": login, "password": password})

	return err
}

func (repository *UserRepository) UserExistsByLogin(login string) (bool, error) {
	return repository.userExistsByProperty("login", login)
}

func (repository *UserRepository) UserExistsByID(id uuid.UUID) (bool, error) {
	return repository.userExistsByProperty("id", id)
}

var checkUserSql = `SELECT true FROM %s WHERE %s = $1`

func (repository *UserRepository) userExistsByProperty(propertyName string, value interface{}) (bool, error) {
	var exists bool

	err := repository.db.Get(&exists, fmt.Sprintf(checkUserSql, tableName, propertyName), value)

	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return exists, nil
}

var getUserSql = `SELECT id, login, password FROM %s WHERE login = $1`

func (repository *UserRepository) GetByLogin(login string) (*model.User, error) {
	var user model.User

	err := repository.db.Get(&user, fmt.Sprintf(getUserSql, tableName), login)

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

var getUserBalanceSql = "SELECT id, user_id, sum, with_drawn FROM user_balance WHERE user_id=$1"

func (repository *UserRepository) GetBalance(userID uuid.UUID) (*model.UserBalance, error) {
	var result model.UserBalance
	err := repository.db.Get(&result, getUserBalanceSql, userID)
	if err != nil {
		repository.logger.Error("Failed to get balance", zap.String("error", err.Error()))
		return &result, err
	}
	return &result, nil
}

var setWithdrawUpdateSql = "UPDATE user_balance SET sum=sum-$1, with_drawn=with_drawn+$1 WHERE user_id=$2"
var setWithdrawCreateSql = "INSERT INTO user_withdrawals (id, user_id, number, sum, created_at) Values ($1, $2, $3, $4, $5)"

func (repository *UserRepository) Withdraw(userID uuid.UUID, number string, sum float64) error {
	tx, err := repository.db.Begin()
	if err != nil {
		repository.logger.Error("Failed to create transaction", zap.String("error", err.Error()))
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(setWithdrawUpdateSql, sum, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return user.ErrNotEnoughtMoney
		}
		repository.logger.Error("Failed to update balance", zap.String("error", err.Error()))
		return err
	}

	_, err = tx.Exec(setWithdrawCreateSql, uuid.New(), userID, number, sum, time.Now())
	if err != nil {
		repository.logger.Error("Failed to insert into withdrawals", zap.String("error", err.Error()))
		return err
	}

	return tx.Commit()
}

var getWithdrawalsSql = `SELECT id, user_id, sum, number, created_at FROM user_withdrawals WHERE user_id=$1`

func (repository *UserRepository) GetWithdrawals(userID uuid.UUID) ([]model.UserWithdrawals, error) {
	var result []model.UserWithdrawals

	err := repository.db.Select(&result, getWithdrawalsSql, userID)
	if err != nil {
		repository.logger.Error("Failed to get withdrawals", zap.String("error", err.Error()))
		return result, err
	}

	return result, nil
}

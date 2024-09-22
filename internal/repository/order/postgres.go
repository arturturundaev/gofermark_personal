package order

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	errors2 "gofermark_personal/internal/errors"
	"gofermark_personal/internal/model"
	"time"
)

const tableName = "orders"

type OrderRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewOrderRepository(dns string, logger *zap.Logger) (*OrderRepository, error) {
	db, err := sqlx.Open("postgres", dns)
	if err != nil {
		return nil, err
	}
	return &OrderRepository{db: db, logger: logger}, nil
}

func (repository *OrderRepository) CreateOrder(id uuid.UUID, userID uuid.UUID, number string, status string, createdAt time.Time, updatedAt time.Time, accrual float64) error {
	tx, err := repository.db.Beginx()
	if err != nil {
		repository.logger.Error("Failed open transaction", zap.String("error", err.Error()))
		return err
	}
	defer tx.Rollback()

	var someUserID *uuid.UUID

	err = tx.Get(
		&someUserID,
		fmt.Sprintf(`SELECT user_id FROM %s WHERE number=$1;`, tableName),
		number,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		repository.logger.Error("failed get user from order", zap.String("error", err.Error()))
		return err
	}

	if someUserID != nil && *someUserID != userID {
		return errors2.ErrConflict
	}

	if someUserID != nil {
		return errors2.ErrExists
	}

	_, err = tx.Exec(fmt.Sprintf(`INSERT INTO %s (id, user_id, number, status, accrual, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`, tableName),
		id,
		userID,
		number,
		status,
		accrual,
		createdAt,
		updatedAt)
	if err != nil {
		repository.logger.Error("Failed insert order to Database", zap.String("error", err.Error()))
		return err
	}
	repository.logger.Info("order success creates", zap.String("number", number))
	_ = tx.Commit()
	return nil
}

func (repository *OrderRepository) GetOrder(number string) (*model.Order, error) {
	repository.logger.Info("get order", zap.String("number", number))
	var result model.Order
	err := repository.db.Get(&result, fmt.Sprintf(`SELECT number, status, accrual, date FROM %s WHERE number=$1;`, tableName), number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors2.ErrInvalid
		}
		repository.logger.Error("Failed to get order from db", zap.String("error", err.Error()))
		return nil, err
	}

	return &result, nil
}

func (repository *OrderRepository) GetOrders(userID uuid.UUID) ([]model.Order, error) {
	var result []model.Order
	err := repository.db.Select(&result, fmt.Sprintf(`SELECT  id, user_id, number, status, accrual, created_at, updated_at FROM %s WHERE user_id=$1;`, tableName), userID)
	if err != nil {
		repository.logger.Error("Failed get orders from Database", zap.String("error", err.Error()))
		return nil, err
	}

	return result, nil
}

func (repository *OrderRepository) UpdateOrder(userID uuid.UUID, number string, status string, accrual float32) error {
	tx, err := repository.db.Begin()
	if err != nil {
		repository.logger.Error("Failed to create transaction", zap.String("error", err.Error()))
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(fmt.Sprintf(`UPDATE %s SET status=$1, accrual=$2 WHERE number=$3`, tableName), status, accrual, number)
	if err != nil {
		repository.logger.Error("Failed to delete urls from db", zap.String("error", err.Error()))
		return err
	}

	query := `INSERT INTO user_balance (id, user_id, sum, with_drawn) Values ($1, $2, $3, $4) ON CONFLICT (user_id) DO UPDATE SET sum=user_balance.sum + EXCLUDED.sum`
	fmt.Println("!!!!!!!!!!!!!" + query + "!!!!!!!!!!!!!")
	_, err = tx.Exec(query, uuid.New().String(), userID, accrual, 0)
	if err != nil {
		repository.logger.Error("Failed to update balance", zap.String("error", err.Error()))
		return err
	}

	return tx.Commit()
}

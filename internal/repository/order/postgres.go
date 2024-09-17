package order

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	errors2 "gofermark_personal/internal/errors"
	"gofermark_personal/internal/model"
	"log"
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

func (repository *OrderRepository) CreateOrder(id uuid.UUID, userId uuid.UUID, number string, status string, createdAt time.Time, updatedAt time.Time, accrual float64) error {
	tx, err := repository.db.Begin()
	if err != nil {
		repository.logger.Error("Failed open transaction", zap.String("error", err.Error()))
		return err
	}
	defer tx.Rollback()

	var someUserId uuid.UUID
	// TODO КАК ЧКРЕЗ ТРАНЗАЦИЮ СРАЗУ НАПОЛНИТЬ СТРУКТУРУ ДАННЫМИ ?!
	rows, err := tx.Query(fmt.Sprintf(`SELECT user_id FROM %s WHERE number=$1;`, tableName), number)
	if err != nil {
		repository.logger.Error("failed get user from order", zap.String("error", err.Error()))
		return err
	}
	if rows.Next() {
		err = rows.Scan(someUserId)
		if err != nil {
			log.Fatalln(err)
			return err
		}
		if someUserId != userId {
			return errors2.ErrConflict
		}

		return errors2.ErrExists
	}
	_, err = tx.Exec(fmt.Sprintf(`INSERT INTO %s (id, user_id, number, status, accrual, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`, tableName),
		id,
		userId,
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

func (repository *OrderRepository) GetOrders(userID uuid.UUID) ([]model.Order, error) {
	var result []model.Order
	err := repository.db.Get(&result, fmt.Sprintf(`SELECT number, status, accrual, date FROM %s WHERE userID=$1;`, tableName), userID)
	if err != nil {
		repository.logger.Error("Failed get orders from Database", zap.String("error", err.Error()))
		return nil, err
	}

	return result, nil
}

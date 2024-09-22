package service

import (
	"github.com/google/uuid"
	"gofermark_personal/internal/model"
	"time"
)

type IOrderRepository interface {
	CreateOrder(id uuid.UUID, userId uuid.UUID, number string, status string, createdAt time.Time, updatedAt time.Time, accrual float64) error
	GetOrder(number string) (*model.Order, error)
	GetOrders(userID uuid.UUID) ([]model.Order, error)
	UpdateOrder(userID uuid.UUID, number string, status string, accrual float32) error
}

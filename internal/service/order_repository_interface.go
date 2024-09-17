package service

import (
	"github.com/google/uuid"
	"gofermark_personal/internal/model"
	"time"
)

type IOrderRepository interface {
	Save(id uuid.UUID, userId uuid.UUID, number string, status string, createdAt time.Time, updatedAt time.Time, accrual float64) error
	CreateOrder(id uuid.UUID, userId uuid.UUID, number string, status string, createdAt time.Time, updatedAt time.Time, accrual float64) error
	GetOrders(userID uuid.UUID) ([]model.Order, error)
}

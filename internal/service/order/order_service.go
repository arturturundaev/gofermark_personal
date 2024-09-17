package order

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gofermark_personal/internal/model"
	"gofermark_personal/internal/service"
	"time"
)

type OrderService struct {
	orderRepository service.IOrderRepository
	userRepository  service.IUserRepository
	logger          *zap.Logger
}

func NewOrderService(orderRepository service.IOrderRepository, userRepository service.IUserRepository, logger *zap.Logger) *OrderService {
	return &OrderService{orderRepository: orderRepository, userRepository: userRepository, logger: logger}
}

func (service *OrderService) Create(userId uuid.UUID, number string) error {
	now := time.Now()
	err := service.orderRepository.CreateOrder(uuid.New(), userId, number, model.ORDER_STATUS_NEW, now, now, 0)

	if err != nil {
		return err
	}

	return nil
}

func (service *OrderService) GetOrders(userId uuid.UUID) ([]model.Order, error) {
	orders, err := service.orderRepository.GetOrders(userId)

	if err != nil {
		return orders, err
	}

	return orders, nil
}

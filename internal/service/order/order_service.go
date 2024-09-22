package order

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gofermark_personal/internal/model"
	"gofermark_personal/internal/service"
	"time"
)

type OrderService struct {
	orderRepository    service.IOrderRepository
	userRepository     service.IUserRepository
	loyalityRepository service.LoyalityRepository
	logger             *zap.Logger
}

func NewOrderService(orderRepository service.IOrderRepository, userRepository service.IUserRepository, loyalityRepository service.LoyalityRepository, logger *zap.Logger) *OrderService {
	return &OrderService{orderRepository: orderRepository, userRepository: userRepository, loyalityRepository: loyalityRepository, logger: logger}
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

func (service *OrderService) GetBalance(userID uuid.UUID) (*model.UserBalance, error) {
	return service.userRepository.GetBalance(userID)
}

func (service *OrderService) Withdraw(userID uuid.UUID, number string, sum float64) error {
	err := service.userRepository.Withdraw(userID, number, sum)
	if err != nil {
		service.logger.Error("failed to withdraw", zap.String("error", err.Error()))
	}
	return nil
}

func (service *OrderService) GetWithdrawals(userID uuid.UUID) ([]model.UserWithdrawals, error) {
	res, err := service.userRepository.GetWithdrawals(userID)
	if err != nil {
		service.logger.Info("failed to withdrawals", zap.String("error", err.Error()))
		return nil, err
	}
	return res, nil
}

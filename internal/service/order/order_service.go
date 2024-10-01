package order

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gofermark_personal/internal/model"
	"time"
)

type orderRepository interface {
	CreateOrder(id uuid.UUID, userID uuid.UUID, number string, status string, createdAt time.Time, updatedAt time.Time, accrual float64) error
	GetOrders(userID uuid.UUID) ([]model.Order, error)
	UpdateOrder(userID uuid.UUID, number string, status string, accrual float32) error
}

type userRepository interface {
	GetBalance(userID uuid.UUID) (*model.UserBalance, error)
	Withdraw(userID uuid.UUID, number string, sum float64) error
	GetWithdrawals(userID uuid.UUID) ([]model.UserWithdrawals, error)
}

type loyaltyService interface {
	SendToGetOrderInfo(orderNumber string, userID uuid.UUID)
}

type OrderService struct {
	orderRepository orderRepository
	userRepository  userRepository
	loyaltyService  loyaltyService
	logger          *zap.Logger
}

func NewOrderService(orderRepository orderRepository, userRepository userRepository, loyaltyService loyaltyService, logger *zap.Logger) *OrderService {
	return &OrderService{orderRepository: orderRepository, userRepository: userRepository, loyaltyService: loyaltyService, logger: logger}
}

func (service *OrderService) Create(userID uuid.UUID, number string) error {
	now := time.Now()
	err := service.orderRepository.CreateOrder(uuid.New(), userID, number, model.OrderStatusNew, now, now, 0)

	if err != nil {
		return err
	}

	err = service.orderRepository.UpdateOrder(userID, number, "PROCESSING", 0)
	if err != nil {
		return err
	}

	service.loyaltyService.SendToGetOrderInfo(number, userID)
	return nil
}

func (service *OrderService) GetOrders(userID uuid.UUID) ([]model.Order, error) {
	orders, err := service.orderRepository.GetOrders(userID)

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

package loyalty

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gofermark_personal/internal/model"
	"time"
)

type loyalityRepository interface {
	GetOrderInfo(orderNumber string) (*model.LoyaltyOrderInfo, error)
}

type orderRepository interface {
	UpdateOrder(userID uuid.UUID, number string, status string, accrual float32) error
}
type LoyaltyService struct {
	repository      loyalityRepository
	orderRepository orderRepository
	logger          *zap.Logger
	ch              chan chanelDto
	ctx             context.Context
}

type chanelDto struct {
	orderNumber string
	userID      uuid.UUID
}

func NewLoyaltyService(repository loyalityRepository, orderRepository orderRepository, logger *zap.Logger, ctx context.Context) *LoyaltyService {
	return &LoyaltyService{repository: repository, orderRepository: orderRepository, logger: logger, ctx: ctx, ch: make(chan chanelDto)}
}

func (s LoyaltyService) Run() {
	for {
		select {
		case v := <-s.ch:
			s.process(v.orderNumber, v.userID)
			time.Sleep(100 * time.Microsecond)
		case <-s.ctx.Done():
			s.logger.Info("context is done")
			return
		}
	}
}

func (s LoyaltyService) SendToGetOrderInfo(orderNumber string, userID uuid.UUID) {
	s.ch <- chanelDto{
		orderNumber: orderNumber,
		userID:      userID,
	}
}

func (s LoyaltyService) process(orderNumber string, userID uuid.UUID) {
	orderInfo, err := s.repository.GetOrderInfo(orderNumber)
	if err != nil {
		s.logger.Error("failed to get order info from loyalty system", zap.String("error", err.Error()))
		s.SendToGetOrderInfo(orderNumber, userID)
	}
	err = s.orderRepository.UpdateOrder(userID, orderInfo.Number, orderInfo.Status, orderInfo.Accrual)
	if err != nil {
		s.logger.Error("failed to update order info", zap.String("error", err.Error()))
		s.SendToGetOrderInfo(orderNumber, userID)
	}
}

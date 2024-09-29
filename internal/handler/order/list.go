package order

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gofermark_personal/internal/helper"
	"gofermark_personal/internal/model"
	"net/http"
)

type orderGetter interface {
	GetOrders(userID uuid.UUID) ([]model.Order, error)
}

type OrderListdHandler struct {
	service orderGetter
	logger  *zap.Logger
}

func NewOrderListHandler(service orderGetter, logger *zap.Logger) *OrderListdHandler {
	return &OrderListdHandler{service: service, logger: logger}
}

func (handler *OrderListdHandler) Handle(ctx *gin.Context) {
	userID, err := helper.GetUserIDFromGin(ctx)
	if err != nil {
		handler.logger.Error("failed to read order number in create order request", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	orders, err := handler.service.GetOrders(*userID)

	if err != nil {
		handler.logger.Error("failed to get orders", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

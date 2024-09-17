package order

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gofermark_personal/internal/helper"
	"gofermark_personal/internal/service/order"
	"net/http"
)

type OrderListdHandler struct {
	service *order.OrderService
	logger  *zap.Logger
}

func NewOrderListdHandler(service *order.OrderService, logger *zap.Logger) *OrderListdHandler {
	return &OrderListdHandler{service: service, logger: logger}
}

func (handler *OrderListdHandler) Handle(ctx *gin.Context) {
	userID, err := helper.GetUserIdFromGin(ctx)
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

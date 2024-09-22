package order

import (
	"errors"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	errors2 "gofermark_personal/internal/errors"
	"gofermark_personal/internal/helper"
	"gofermark_personal/internal/service/order"
	"io"
	"net/http"
)

type OrderUploadHandler struct {
	service *order.OrderService
	logger  *zap.Logger
}

func NewOrderUploadHandler(service *order.OrderService, logger *zap.Logger) *OrderUploadHandler {
	return &OrderUploadHandler{service: service, logger: logger}
}

func (handler *OrderUploadHandler) Handler(ctx *gin.Context) {
	orderNumber, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		handler.logger.Error("failed to read order number in create order request", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	strOrderNumber := string(orderNumber)
	err = goluhn.Validate(strOrderNumber)
	if err != nil {
		handler.logger.Error("failed to validate order number", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	userID, err := helper.GetUserIdFromGin(ctx)
	if err != nil {
		handler.logger.Error("failed to read order number in create order request", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = handler.service.Create(*userID, strOrderNumber)
	if err != nil {
		if errors.Is(err, errors2.ErrExists) {
			ctx.AbortWithStatus(http.StatusOK)
			return
		}
		if errors.Is(err, errors2.ErrConflict) {
			ctx.AbortWithStatus(http.StatusConflict)
			return
		}
		handler.logger.Error("failed to create order", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.AbortWithStatus(http.StatusAccepted)
}

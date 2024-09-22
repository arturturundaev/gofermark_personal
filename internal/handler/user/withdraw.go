package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	errors2 "gofermark_personal/internal/errors"
	"gofermark_personal/internal/helper"
	"gofermark_personal/internal/service/user"
	"net/http"
)

type UserWithdrawHandler struct {
	userService *user.UserService
	logger      *zap.Logger
}

func NewUserWithdrawHandler(userService *user.UserService, logger *zap.Logger) *UserWithdrawHandler {
	return &UserWithdrawHandler{userService: userService, logger: logger}
}

func (handler *UserWithdrawHandler) Handler(ctx *gin.Context) {
	var req struct {
		Number string  `json:"order"`
		Sum    float64 `json:"sum"`
	}

	err := ctx.BindJSON(&req)
	if err != nil {
		handler.logger.Error("failed to decode body", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userID, err := helper.GetUserIdFromGin(ctx)
	if err != nil {
		handler.logger.Error("failed to get balance", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = handler.userService.Withdraw(*userID, req.Number, req.Sum)
	if errors.Is(err, errors2.ErrInvalid) {
		handler.logger.Error("failed to get balance", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	if errors.Is(err, errors2.ErrNotEnoughtMoney) {
		handler.logger.Error("failed to get balance", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusPaymentRequired)
		return
	}

	if err != nil {
		handler.logger.Error("failed to get balance", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

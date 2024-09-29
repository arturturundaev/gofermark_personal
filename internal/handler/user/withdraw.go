package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gofermark_personal/internal/helper"
	"net/http"
)

var ErrInvalid = errors.New("invalid")
var ErrNotEnoughtMoney = errors.New("not enought money")

type withdrawService interface {
	Withdraw(userID uuid.UUID, number string, sum float64) error
}

type UserWithdrawHandler struct {
	userService withdrawService
	logger      *zap.Logger
}

func NewUserWithdrawHandler(userService withdrawService, logger *zap.Logger) *UserWithdrawHandler {
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

	userID, err := helper.GetUserIDFromGin(ctx)
	if err != nil {
		handler.logger.Error("failed to get balance", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = handler.userService.Withdraw(*userID, req.Number, req.Sum)
	if errors.Is(err, ErrInvalid) {
		handler.logger.Error("failed to get balance", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	if errors.Is(err, ErrNotEnoughtMoney) {
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

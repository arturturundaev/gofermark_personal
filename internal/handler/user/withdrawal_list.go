package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gofermark_personal/internal/helper"
	"gofermark_personal/internal/model"
	"net/http"
)

type withdrawGetter interface {
	GetWithdrawals(userID uuid.UUID) ([]model.UserWithdrawals, error)
}

type UserWithdrawalList struct {
	userService withdrawGetter
	logger      *zap.Logger
}

func NewUserWithdrawalList(userService withdrawGetter, logger *zap.Logger) *UserWithdrawalList {
	return &UserWithdrawalList{userService: userService, logger: logger}
}

func (handler UserWithdrawalList) Handle(ctx *gin.Context) {

	userID, err := helper.GetUserIDFromGin(ctx)
	if err != nil {
		handler.logger.Error("failed to get balance", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	resp, err := handler.userService.GetWithdrawals(*userID)
	if err != nil {
		handler.logger.Error("failed to get withdrawals", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(resp) == 0 {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

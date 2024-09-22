package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gofermark_personal/internal/helper"
	"gofermark_personal/internal/service/user"
	"net/http"
)

type UserWithdrawalList struct {
	userService *user.UserService
	logger      *zap.Logger
}

func NewUserWithdrawalList(userService *user.UserService, logger *zap.Logger) *UserWithdrawalList {
	return &UserWithdrawalList{userService: userService, logger: logger}
}

func (handler UserWithdrawalList) Handle(ctx *gin.Context) {

	userID, err := helper.GetUserIdFromGin(ctx)
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

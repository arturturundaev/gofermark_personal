package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gofermark_personal/internal/helper"
	"gofermark_personal/internal/service/user"
	"net/http"
)

type UserBalanceHandler struct {
	userService *user.UserService
	logger      *zap.Logger
}

func NewUserBalanceHandler(userService *user.UserService, logger *zap.Logger) *UserBalanceHandler {
	return &UserBalanceHandler{userService: userService, logger: logger}
}

func (handler *UserBalanceHandler) Handle(ctx *gin.Context) {
	userID, err := helper.GetUserIDFromGin(ctx)
	if err != nil {
		handler.logger.Error("failed to get balance", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	balance, err := handler.userService.GetBalance(*userID)
	if err != nil {
		handler.logger.Error("failed to get balance", zap.String("error", err.Error()))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, balance)
}

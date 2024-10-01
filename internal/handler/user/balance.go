package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gofermark_personal/internal/helper"
	"gofermark_personal/internal/model"
	"net/http"
)

type userBalancer interface {
	GetBalance(userID uuid.UUID) (*model.UserBalance, error)
}

type UserBalanceHandler struct {
	userService userBalancer
	logger      *zap.Logger
}

func NewUserBalanceHandler(userService userBalancer, logger *zap.Logger) *UserBalanceHandler {
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

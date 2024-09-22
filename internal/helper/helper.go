package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gofermark_personal/internal/middleware"
)

func GetUserIdFromGin(ctx *gin.Context) (*uuid.UUID, error) {
	userId, exists := ctx.Get(middleware.UserIDProperty)

	if !exists {
		return nil, fmt.Errorf("user does not exsits in context")
	}

	userUUID := userId.(uuid.UUID)

	return &userUUID, nil
}

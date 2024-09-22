package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gofermark_personal/internal/middleware"
)

func GetUserIDFromGin(ctx *gin.Context) (*uuid.UUID, error) {
	userID, exists := ctx.Get(middleware.UserIDProperty)

	if !exists {
		return nil, fmt.Errorf("user does not exsits in context")
	}

	userUUID := userID.(uuid.UUID)

	return &userUUID, nil
}

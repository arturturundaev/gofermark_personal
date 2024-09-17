package user

import (
	"github.com/gin-gonic/gin"
	"gofermark_personal/internal/middleware"
	"gofermark_personal/internal/model"
	"gofermark_personal/internal/service/user"
	"net/http"
)

type UserRegisterHandler struct {
	service      *user.UserService
	JWTValidator *middleware.JWTValidator
}

func NewUserRegisterHandler(service *user.UserService, JWTValidator *middleware.JWTValidator) *UserRegisterHandler {
	return &UserRegisterHandler{service: service, JWTValidator: JWTValidator}
}

func (handler *UserRegisterHandler) Handler(ctx *gin.Context) {
	currentUser, err := model.NewUserFromGin(ctx)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	exists, err := handler.service.UserExists(currentUser.Login)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	if exists {
		ctx.AbortWithStatus(http.StatusConflict)
	}
	userID, errSave := handler.service.Register(currentUser.Login, currentUser.Password)

	if errSave != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	err = handler.JWTValidator.InitToken(ctx, userID)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	ctx.AbortWithStatus(http.StatusOK)
}

package user

import (
	"github.com/gin-gonic/gin"
	"gofermark_personal/internal/middleware"
	"gofermark_personal/internal/model"
	"gofermark_personal/internal/service/user"
	"net/http"
)

type UserLoginHandler struct {
	service      *user.UserService
	JWTValidator *middleware.JWTValidator
}

func NewUserLoginHandler(service *user.UserService, JWTValidator *middleware.JWTValidator) *UserLoginHandler {
	return &UserLoginHandler{service: service, JWTValidator: JWTValidator}
}

func (handler *UserLoginHandler) Handler(ctx *gin.Context) {
	dto, err := model.NewUserFromGin(ctx)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := handler.service.Auth(dto.Login, dto.Password)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if user == nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err = handler.JWTValidator.InitToken(ctx, &user.Id)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.AbortWithStatus(http.StatusOK)
	return
}

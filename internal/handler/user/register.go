package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gofermark_personal/internal/model"
	"net/http"
)

type userRegister interface {
	UserExists(login string) (bool, error)
	Register(login string, password string) (*uuid.UUID, error)
}

type jwtValidator interface {
	InitToken(ctx *gin.Context, userID *uuid.UUID) error
}

type UserRegisterHandler struct {
	service      userRegister
	JWTValidator jwtValidator
}

func NewUserRegisterHandler(service userRegister, JWTValidator jwtValidator) *UserRegisterHandler {
	return &UserRegisterHandler{service: service, JWTValidator: JWTValidator}
}

func (handler *UserRegisterHandler) Handler(ctx *gin.Context) {
	currentUser, err := model.NewUserFromGin(ctx)

	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Abort()
		return
	}

	exists, err := handler.service.UserExists(currentUser.Login)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Abort()
		return
	}

	if exists {
		ctx.Writer.WriteHeader(http.StatusConflict)
		ctx.Abort()
		return
	}
	userID, errSave := handler.service.Register(currentUser.Login, currentUser.Password)

	if errSave != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Abort()
		return
	}

	err = handler.JWTValidator.InitToken(ctx, userID)

	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Abort()
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
}

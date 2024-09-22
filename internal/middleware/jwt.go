package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gofermark_personal/internal/service"
	"net/http"
	"strings"
	"time"
)

type JWTValidator struct {
	Claims         Claims
	userRepository service.IUserRepository
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

func NewJWTValidator(userRepository service.IUserRepository) *JWTValidator {
	return &JWTValidator{userRepository: userRepository}
}

const TokenExp = 3 * time.Hour
const SecretKey = "0N#6Ke|+OR:(`G;"
const UserIDProperty = "UserId"
const HeaderTokenProperty = "Authorization"

// Handle Проверяем токен пользователя
// В случае успеха продлеваем токен
func (JWTValidator *JWTValidator) Handle(ctx *gin.Context) {
	token := ctx.GetHeader(HeaderTokenProperty)

	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "")
		return
	}

	userID, errorValidateToken := JWTValidator.getUserIDFromToken(token)

	if errorValidateToken != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorValidateToken)
		return
	}

	if userID == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "")
		return
	}

	userExists, err := JWTValidator.userRepository.UserExistsByID(*userID)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "")
		return
	}

	if !userExists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "")
		return
	}

	err = JWTValidator.InitToken(ctx, userID)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, "")
		return
	}

	ctx.Set(UserIDProperty, *userID)
}

func (JWTValidator *JWTValidator) getUserIDFromToken(tokenString string) (*uuid.UUID, error) {
	fmt.Println(tokenString)

	claims := &Claims{}
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, _ := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return &claims.UserID, nil
}

func (JWTValidator *JWTValidator) buildJWT(userID *uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: *userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (JWTValidator *JWTValidator) InitToken(ctx *gin.Context, userID *uuid.UUID) error {
	token, err := JWTValidator.buildJWT(userID)

	if err != nil {
		return err
	}

	ctx.Header(HeaderTokenProperty, token)

	return nil
}

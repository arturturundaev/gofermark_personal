package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

const UserIDProperty = "UserId"

type UserChecker interface {
	UserExistsByID(id uuid.UUID) (bool, error)
}

type JWTValidator struct {
	Claims              Claims
	userRepository      UserChecker
	tokenExpire         time.Duration
	secretKey           string
	headerTokenProperty string
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

func NewJWTValidator(userRepository UserChecker, tokenExpire time.Duration, secretKey string, headerTokenProperty string) *JWTValidator {
	return &JWTValidator{
		userRepository:      userRepository,
		tokenExpire:         tokenExpire,
		secretKey:           secretKey,
		headerTokenProperty: headerTokenProperty,
	}
}

// Handle Проверяем токен пользователя
// В случае успеха продлеваем токен
func (JWTValidator *JWTValidator) Handle(ctx *gin.Context) {
	token := ctx.GetHeader(JWTValidator.headerTokenProperty)

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
			return []byte(JWTValidator.secretKey), nil
		})

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return &claims.UserID, nil
}

func (JWTValidator *JWTValidator) buildJWT(userID *uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWTValidator.tokenExpire)),
		},
		UserID: *userID,
	})

	tokenString, err := token.SignedString([]byte(JWTValidator.secretKey))
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

	ctx.Header(JWTValidator.headerTokenProperty, token)

	return nil
}

package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	orderHandler "gofermark_personal/internal/handler/order"
	userHandler "gofermark_personal/internal/handler/user"
	"gofermark_personal/internal/middleware"
	"gofermark_personal/internal/service/order"
	"gofermark_personal/internal/service/user"
	"time"
)

func NewRouter(
	userRepository middleware.UserChecker,
	userService *user.UserService,
	orderService *order.OrderService,
	logger *zap.Logger,
	tokenExp time.Duration,
	secretKey string,
	headerTokenProperty string,

) *gin.Engine {
	router := gin.Default()

	JWTValidator := middleware.NewJWTValidator(userRepository, tokenExp, secretKey, headerTokenProperty)

	userRegisterHandler := userHandler.NewUserRegisterHandler(userService, JWTValidator)
	orderListHandler := orderHandler.NewOrderListHandler(orderService, logger)
	orderCreateHandler := orderHandler.NewOrderUploadHandler(orderService, logger)

	loginHandler := userHandler.NewUserLoginHandler(userService, JWTValidator)
	userBalanceHandler := userHandler.NewUserBalanceHandler(userService, logger)
	userWithdraw := userHandler.NewUserWithdrawHandler(userService, logger)
	userWithdrawalList := userHandler.NewUserWithdrawalList(userService, logger)

	router.POST("/api/user/register", userRegisterHandler.Handler)
	router.POST("/api/user/login", loginHandler.Handler)
	router.GET("/api/user/balance", JWTValidator.Handle, userBalanceHandler.Handle)
	router.POST("/api/user/balance/withdraw", JWTValidator.Handle, userWithdraw.Handler)
	router.GET("/api/user/withdrawals", JWTValidator.Handle, userWithdrawalList.Handle)

	router.POST("/api/user/orders", JWTValidator.Handle, orderCreateHandler.Handler)
	router.GET("/api/user/orders", JWTValidator.Handle, orderListHandler.Handle)

	return router
}

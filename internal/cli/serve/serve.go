package serve

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gofermark_personal/internal/config"
	orderHandler "gofermark_personal/internal/handler/order"
	userHandler "gofermark_personal/internal/handler/user"
	"gofermark_personal/internal/logger"
	"gofermark_personal/internal/middleware"
	"gofermark_personal/internal/repository/loyalty"
	"gofermark_personal/internal/repository/order"
	"gofermark_personal/internal/repository/user"
	orderService "gofermark_personal/internal/service/order"
	userService "gofermark_personal/internal/service/user"
	"log"
	"net/http"
)

func StartServer() {

	cfg := config.NewConfig()

	router := gin.Default()

	logger, err := logger.AddLoggerToGIN(router)
	if err != nil {
		log.Fatal(err)
	}

	postgres, err := sqlx.Open("postgres", cfg.DatabaseDSN)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// "postgres://postgres:postgres@localhost:5432/gofermark?sslmode=disable"
	userRepository := user.NewUserRepository(postgres, logger)

	orderRepository := order.NewOrderRepository(postgres, logger)

	loyalityRepository := loyalty.NewLoyaltyHTTPRepository(cfg.AccrualSystemURL, http.DefaultClient, logger)
	userService := userService.NewUserService(userRepository)
	orderService := orderService.NewOrderService(orderRepository, userRepository, loyalityRepository, logger)
	JWTValidator := middleware.NewJWTValidator(userRepository, cfg.TokenExp, cfg.SecretKey, cfg.HeaderTokenProperty)

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

	errServer := http.ListenAndServe(cfg.ServerAddress, router)
	if errServer != nil {
		logger.Fatal(errServer.Error())
	}
}

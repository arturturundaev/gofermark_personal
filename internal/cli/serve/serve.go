package serve

import (
	"context"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gofermark_personal/internal/config"
	"gofermark_personal/internal/db"
	currentLogger "gofermark_personal/internal/logger"
	"gofermark_personal/internal/repository/loyalty"
	"gofermark_personal/internal/repository/order"
	"gofermark_personal/internal/repository/user"
	"gofermark_personal/internal/router"
	loyaltyService "gofermark_personal/internal/service/loyalty"
	orderService "gofermark_personal/internal/service/order"
	userService "gofermark_personal/internal/service/user"
	"log"
	"net/http"
)

func StartServer() {

	cfg := config.NewConfig()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	postgres, err := db.NewPostgres(cfg.Postgres.DSN, cfg.Postgres.MaxOpen, cfg.Postgres.MaxIdle, cfg.Postgres.MaxLifeDuration)
	if err != nil {
		logger.Fatal(err.Error())
	}

	ctx := context.Background()
	// "postgres://postgres:postgres@localhost:5432/gofermark?sslmode=disable"
	userRepository := user.NewUserRepository(postgres, logger)

	orderRepository := order.NewOrderRepository(postgres, logger)

	loyalityRepository := loyalty.NewLoyaltyHTTPRepository(cfg.AccrualSystemURL, http.DefaultClient, logger)
	userService := userService.NewUserService(userRepository)
	loyaltyService := loyaltyService.NewLoyaltyService(loyalityRepository, orderRepository, logger, ctx)
	orderService := orderService.NewOrderService(orderRepository, userRepository, loyaltyService, logger)

	r := router.NewRouter(
		userRepository,
		userService,
		orderService,
		logger,
		cfg.TokenExp,
		cfg.SecretKey,
		cfg.HeaderTokenProperty,
	)

	err = currentLogger.AddLoggerToGIN(logger, r)
	if err != nil {
		logger.Fatal(err.Error())
	}

	go loyaltyService.Run()

	errServer := http.ListenAndServe(cfg.ServerAddress, r)
	if errServer != nil {
		logger.Fatal(errServer.Error())
	}
}

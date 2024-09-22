package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	config2 "gofermark_personal/internal/config"
	order3 "gofermark_personal/internal/handler/order"
	user3 "gofermark_personal/internal/handler/user"
	logger2 "gofermark_personal/internal/logger"
	"gofermark_personal/internal/middleware"
	"gofermark_personal/internal/repository/loyalty"
	"gofermark_personal/internal/repository/order"
	"gofermark_personal/internal/repository/user"
	order2 "gofermark_personal/internal/service/order"
	user2 "gofermark_personal/internal/service/user"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	config := config2.GetConfig()

	router := gin.Default()

	logger, err := logger2.AddLoggerToGIN(router)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(5)
	}

	// "postgres://postgres:postgres@localhost:5432/gofermark?sslmode=disable"
	userRepository, err := user.NewUserRepository(config.DatabaseDSN, logger)
	if err != nil {
		logger.Error(err.Error())
	}

	orderRepository, err := order.NewOrderRepository(config.DatabaseDSN, logger)
	if err != nil {
		logger.Error(err.Error())
	}

	absPath, errPathMigration := filepath.Abs(".")
	if errPathMigration != nil {
		fmt.Println("ошибка определения директории для миграций!")
	} else {
		initMigrations("file:////"+absPath+"/internal/repository/migration", userRepository.GetDB())
	}

	loyalityRepository := loyalty.NewLoyaltyHTTPRepository(config.AccrualSystemURL, http.DefaultClient, logger)
	userService := user2.NewUserService(userRepository)
	orderService := order2.NewOrderService(orderRepository, userRepository, loyalityRepository, logger)
	JWTValidator := middleware.NewJWTValidator(userRepository)

	userRegisterHandler := user3.NewUserRegisterHandler(userService, JWTValidator)
	orderListHandler := order3.NewOrderListHandler(orderService, logger)
	orderCreateHandler := order3.NewOrderUploadHandler(orderService, logger)

	loginHandler := user3.NewUserLoginHandler(userService, JWTValidator)
	userBalanceHandler := user3.NewUserBalanceHandler(userService, logger)
	userWithdraw := user3.NewUserWithdrawHandler(userService, logger)
	userWithdrawalList := user3.NewUserWithdrawalList(userService, logger)

	router.POST("/api/user/register", userRegisterHandler.Handler)
	router.POST("/api/user/login", loginHandler.Handler)
	router.GET("/api/user/balance", JWTValidator.Handle, userBalanceHandler.Handle)
	router.POST("/api/user/balance/withdraw", JWTValidator.Handle, userWithdraw.Handler)
	router.GET("/api/user/withdrawals", JWTValidator.Handle, userWithdrawalList.Handle)

	router.POST("/api/user/orders", JWTValidator.Handle, orderCreateHandler.Handler)
	router.GET("/api/user/orders", JWTValidator.Handle, orderListHandler.Handle)

	errServer := http.ListenAndServe(config.ServerAddress, router)
	if errServer != nil {
		fmt.Println(errServer.Error())
	}
}

func initMigrations(migrationPath string, DB *sqlx.DB) {
	driver, err := postgres.WithInstance(DB.DB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres", driver)

	if err != nil {
		log.Fatal(err)
	} else {
		errMigrate := m.Up()
		if errMigrate != nil && errMigrate.Error() != "no change" {
			log.Fatal(errMigrate)
		}
	}
}

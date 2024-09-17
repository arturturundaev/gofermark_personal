package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	user3 "gofermark_personal/internal/handler/user"
	logger2 "gofermark_personal/internal/logger"
	"gofermark_personal/internal/middleware"
	"gofermark_personal/internal/repository/order"
	"gofermark_personal/internal/repository/user"
	user2 "gofermark_personal/internal/service/user"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	router := gin.Default()

	logger, err := logger2.AddLoggerToGIN(router)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(5)
	}

	userRepository, err := user.NewUserRepository("postgres://postgres:postgres@localhost:5432/gofermark?sslmode=disable")
	if err != nil {
		logger.Error(err.Error())
	}

	orderRepository, err := order.NewOrderRepository("postgres://postgres:postgres@localhost:5432/gofermark?sslmode=disable", logger)
	if err != nil {
		logger.Error(err.Error())
	}

	absPath, errPathMigration := filepath.Abs(".")
	if errPathMigration != nil {
		fmt.Println("ошибка определения директории для миграций!")
	} else {
		initMigrations("file:////"+absPath+"/internal/repository/migration", userRepository.GetDB())
	}

	userService := user2.NewUserService(userRepository)
	JWTValidator := middleware.NewJWTValidator(userRepository)

	registerHandler := user3.NewUserRegisterHandler(userService, JWTValidator)
	loginHandler := user3.NewUserLoginHandler(userService, JWTValidator)

	router.POST("/api/user/register", registerHandler.Handler)
	router.POST("/api/user/login", loginHandler.Handler)

	errServer := http.ListenAndServe(":8082", router)
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

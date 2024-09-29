package main

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"gofermark_personal/internal/config"
	"log"
	"path/filepath"
)

func main() {
	config := config.NewConfig()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Open("postgres", config.DatabaseDSN)
	if err != nil {
		logger.Fatal(err.Error())
	}
	absPath, errPathMigration := filepath.Abs(".")
	if errPathMigration != nil {
		log.Fatal("ошибка определения директории для миграций!")
	} else {
		driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
		if err != nil {
			log.Fatal(err)
		}
		m, err := migrate.NewWithDatabaseInstance(
			"file:////"+absPath+"/internal/repository/migration",
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

	logger.Info("Migration complete")
}

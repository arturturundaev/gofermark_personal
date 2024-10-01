package config

import (
	"flag"
	"os"
	"time"
)

type PostgresConfig struct {
	DSN             string
	MaxOpen         int
	MaxIdle         int
	MaxLifeDuration time.Duration
}

type Config struct {
	ServerAddress       string
	AccrualSystemURL    string
	LogLevel            string
	TokenExp            time.Duration
	SecretKey           string
	HeaderTokenProperty string
	Postgres            *PostgresConfig
}

func NewConfig() *Config {

	pgConfig := PostgresConfig{
		MaxOpen:         10,
		MaxIdle:         10,
		MaxLifeDuration: 10 * time.Minute,
	}

	config := &Config{
		LogLevel:            "info",
		TokenExp:            3 * time.Hour,
		SecretKey:           "0N#6Ke|+OR:(`G;",
		HeaderTokenProperty: "Authorization",
		Postgres:            &pgConfig,
	}
	flag.StringVar(&config.ServerAddress, "a", "", "run address")
	flag.StringVar(&config.Postgres.DSN, "d", "", "database uri")
	flag.StringVar(&config.AccrualSystemURL, "r", "", "accrual system address")
	flag.Parse()

	if runAddress, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		config.ServerAddress = runAddress
	}

	if databaseURI, ok := os.LookupEnv("DATABASE_URI"); ok {
		config.Postgres.DSN = databaseURI
	}

	if accrualSystemAddress, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); ok {
		config.AccrualSystemURL = accrualSystemAddress
	}

	return config
}

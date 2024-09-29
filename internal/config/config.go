package config

import (
	"flag"
	"os"
	"time"
)

type Config struct {
	ServerAddress       string
	DatabaseDSN         string
	AccrualSystemURL    string
	LogLevel            string
	TokenExp            time.Duration
	SecretKey           string
	HeaderTokenProperty string
}

func NewConfig() *Config {
	config := &Config{
		LogLevel:            "info",
		TokenExp:            3 * time.Hour,
		SecretKey:           "0N#6Ke|+OR:(`G;",
		HeaderTokenProperty: "Authorization",
	}
	flag.StringVar(&config.ServerAddress, "a", "", "run address")
	flag.StringVar(&config.DatabaseDSN, "d", "", "database uri")
	flag.StringVar(&config.AccrualSystemURL, "r", "", "accrual system address")
	flag.Parse()

	if runAddress, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		config.ServerAddress = runAddress
	}

	if databaseURI, ok := os.LookupEnv("DATABASE_URI"); ok {
		config.DatabaseDSN = databaseURI
	}

	if accrualSystemAddress, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); ok {
		config.AccrualSystemURL = accrualSystemAddress
	}

	return config
}

package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress    string
	DatabaseDSN      string
	AccrualSystemURL string
	LogLevel         string
}

func GetConfig() *Config {
	config := &Config{LogLevel: "info"}
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

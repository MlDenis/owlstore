package main

import (
	"flag"
	"os"
	"strconv"
)

type FlagVar struct {
	runAddr             string
	databaseURI         string
	acuralSystemAddress string
	migrationsDir       string
	rateLimit           int
	logLevel            string
}

func NewFlagVarStruct() *FlagVar {
	return &FlagVar{}
}
func (f *FlagVar) parseFlags() error {
	//postgresql://postgres:docker@localhost:5432/asd?sslmode=disable
	// как аргумент -a со значением :8080 по умолчанию
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.StringVar(&f.runAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&f.logLevel, "l", "info", "log level")
	flag.StringVar(&f.databaseURI, "d", "", "database connection address")
	flag.StringVar(&f.acuralSystemAddress, "r", "localhost:8081", "address of the accrual system")
	flag.StringVar(&f.migrationsDir, "m", "", "migrations to db")
	flag.IntVar(&f.rateLimit, "w", 10, "number of source related materials on the server")
	flag.Parse()
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		f.logLevel = envLogLevel
	}
	if envRunAddr, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		f.runAddr = envRunAddr
	}

	if envDatabaseURI, ok := os.LookupEnv("DATABASE_URI"); ok {
		f.databaseURI = envDatabaseURI
	}

	if envAcuralSystemAddress, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); ok {
		f.acuralSystemAddress = envAcuralSystemAddress
	}

	if envMigrationsDir, ok := os.LookupEnv("MIGRATIONS_DIR"); ok {
		f.migrationsDir = envMigrationsDir
	}

	if envRateLimit, ok := os.LookupEnv("RATE_LIMIT"); ok {
		envRateLimitInt, err := strconv.Atoi(envRateLimit)
		if err != nil {
			return err
		}
		f.rateLimit = envRateLimitInt
	}
	return nil
}

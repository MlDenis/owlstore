package main

import (
	"flag"
	"os"
	"strconv"
)

type FlagVar struct {
	runAddr       string
	databaseURI   string
	migrationsDir string
	rateLimit     int
}

func NewFlagVarStruct() *FlagVar {
	return &FlagVar{}
}
func (f *FlagVar) parseFlags() error {
	//postgresql://postgres:docker@localhost:5432/asd?sslmode=disable
	// как аргумент -a со значением :8080 по умолчанию
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.StringVar(&f.runAddr, "a", "localhost:8081", "address and port to run server")
	flag.StringVar(&f.databaseURI, "d", "", "database connection address")
	flag.StringVar(&f.migrationsDir, "m", "migrations", "migrations to db")
	flag.IntVar(&f.rateLimit, "l", 10, "number of source related materials on the server")
	flag.Parse()

	if envRunAddr, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		f.runAddr = envRunAddr
	}

	if envDatabaseURI, ok := os.LookupEnv("DATABASE_URI"); ok {
		f.databaseURI = envDatabaseURI
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

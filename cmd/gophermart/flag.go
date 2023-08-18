package main

import (
	"flag"
	"os"
)

type FlagVar struct {
	runAddr             string
	databaseURI         string
	acuralSystemAddress string
	migrationsDir       string
}

func NewFlagVarStruct() *FlagVar {
	return &FlagVar{}
}
func (f *FlagVar) parseFlags() error {
	//postgresql://postgres:docker@localhost:5432/asd?sslmode=disable
	// как аргумент -a со значением :8080 по умолчанию
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.StringVar(&f.runAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&f.databaseURI, "d", "", "database connection address")
	flag.StringVar(&f.acuralSystemAddress, "r", "", "address of the accrual system")
	flag.StringVar(&f.migrationsDir, "m", "", "migrations to db")
	flag.Parse()

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
	return nil
}

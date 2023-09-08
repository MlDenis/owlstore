package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/MlDenis/internal/accrual/accrualcalculate"
	"github.com/MlDenis/internal/accrual/handlers"
	"github.com/MlDenis/internal/accrual/storage"
	"github.com/MlDenis/logger"
	"go.uber.org/zap"
)

func main() {
	flagStruct := NewFlagVarStruct()
	err := flagStruct.parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	if err := run(flagStruct); err != nil {
		log.Fatalln(err)
	}
}
func run(flagStruct *FlagVar) error {
	log, err := logger.InitializeLogger(flagStruct.logLevel)
	if err != nil {
		return err
	}
	ctx := context.Background()
	memStorageInterface, postgresDB, err := storage.NewStorage(ctx, flagStruct.migrationsDir, flagStruct.databaseURI, log)
	if err != nil {
		log.Fatal("Error in create storage: ", zap.Error(err))
	}
	if postgresDB != nil {
		defer postgresDB.Close()
		log.Error("Error in connect to DB")
		return fmt.Errorf("Error in connect to DB")
	}
	newHandStruct := handlers.HandlerNew(memStorageInterface)

	go accrualcalculate.WorkerPool(ctx, memStorageInterface, flagStruct.rateLimit, log)
	router := handlers.Router(ctx, log, newHandStruct)
	log.Info("Running server on: ", zap.String("", flagStruct.runAddr))
	return http.ListenAndServe(flagStruct.runAddr, router)
}

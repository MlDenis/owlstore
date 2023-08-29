package main

import (
	"context"
	"log"
	"net/http"

	"github.com/MlDenis/internal/accrual/accrualcalculate"
	"github.com/MlDenis/internal/accrual/handlers"
	"github.com/MlDenis/internal/accrual/storage"
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

	ctx := context.Background()
	memStorageInterface, postgresDB, err := storage.NewStorage(ctx, flagStruct.migrationsDir, flagStruct.databaseURI)
	if err != nil {
		log.Fatal("Error in create storage: ", err)
	}
	if postgresDB != nil {
		defer postgresDB.Close()
	}
	newHandStruct := handlers.HandlerNew(memStorageInterface)

	go accrualcalculate.WorkerPool(ctx, memStorageInterface, flagStruct.rateLimit)
	router := handlers.Router(ctx, newHandStruct)
	log.Println("Running server on: ", flagStruct.runAddr)
	return http.ListenAndServe(flagStruct.runAddr, router)
}

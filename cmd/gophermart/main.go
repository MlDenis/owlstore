package main

import (
	"context"
	"log"
	"net/http"

	"github.com/MlDenis/internal/gofermart/auth/cache"
	"github.com/MlDenis/internal/gofermart/handlers"
	"github.com/MlDenis/internal/gofermart/storage"
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
	JWTForSession := cache.NewDataJWT()
	newHandStruct := handlers.HandlerNew(memStorageInterface, JWTForSession)
	router := handlers.Router(ctx, newHandStruct)
	log.Println("Running server on: ", flagStruct.runAddr)
	return http.ListenAndServe(flagStruct.runAddr, router)
}

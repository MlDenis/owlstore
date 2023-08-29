package main

import (
	"context"
	"log"
	"net/http"

	"github.com/MlDenis/internal/gofermart/auth/cache"
	"github.com/MlDenis/internal/gofermart/handlers"
	"github.com/MlDenis/internal/gofermart/interactionwithaccrual"
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
	go interactionwithaccrual.WorkerPool(ctx, memStorageInterface, flagStruct.rateLimit, flagStruct.acuralSystemAddress)
	log.Println("Running server on: ", flagStruct.runAddr)
	return http.ListenAndServe(flagStruct.runAddr, router)
}

// curl -v --header "Content-Type: application/json"   --request POST   --data '{"login":"xasf","password":"xyz"}'   http://localhost:8080/api/user/register
// curl -X POST -H "Content-Type: text/plain" -H "Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTIzNDA1ODksInVzZXJuYW1lIjoieGFzZiJ9.If3cohJbnVaVGbTUzrc5Ni5KR9u64-fXUJqLvpY-Mpo" --data "4561261212345467" http://localhost:8080/api/user/orders

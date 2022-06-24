// @title Payment Service
// @version 1.0
// @description The MelonWallet microservice responsible for dealing with payment and crypto transaction information.

package main

import (
	"log"
	"net/http"

	"github.com/Melon-Network-Inc/payment-service/pkg/db"
	transaction "github.com/Melon-Network-Inc/payment-service/pkg/transaction"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("Payment Service Started")
	db := db.Init()
	r := mux.NewRouter()

	transaction.RegisterHandlers(r, db)

	log.Println("Going to listen on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

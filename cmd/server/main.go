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

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request")
	w.Write([]byte(transaction.RegisterHandlers()))
}

func main() {
	log.Println("Payment Service Started")
	DB := db.Init()
	h := transaction.New(DB)
	r := mux.NewRouter()

	// Bind to a port and pass our router
	r.HandleFunc("/transactions", h.GetAllTransactions).Methods(http.MethodGet)
	r.HandleFunc("/transactions/{id}", h.GetTransaction).Methods(http.MethodGet)
	r.HandleFunc("/transactions", h.AddTransaction).Methods(http.MethodPost)
	r.HandleFunc("/transactions/{id}", h.UpdateTransaction).Methods(http.MethodPut)
	r.HandleFunc("/transactions/{id}", h.DeleteTransaction).Methods(http.MethodDelete)

	log.Println("Going to listen on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

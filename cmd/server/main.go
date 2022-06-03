package main

import (
	"fmt"
	"log"
	"net/http"

	transaction "github.com/Melon-Network-Inc/payment-service/api/transaction"
	"github.com/gorilla/mux"
)

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request")
	w.Write([]byte(transaction.RegisterHandlers()))
}

func main() {
	fmt.Println("Payment Service Started")
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", TransactionHandler)

	// Bind to a port and pass our router in
	log.Println("Going to listen on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

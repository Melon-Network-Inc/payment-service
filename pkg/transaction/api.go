package transaction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/Melon-Network-Inc/payment-service/pkg/entity"
	"github.com/gorilla/mux"
)

func RegisterHandlers() string {
	return "Transaction Received!"
}

func (h handler) AddTransaction(w http.ResponseWriter, r *http.Request) {
	// Read to request body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var transaction entity.Transaction
	json.Unmarshal(body, &transaction)

	// Append to the Transaction table
	if result := h.DB.Create(&transaction); result.Error != nil {
		fmt.Println(result.Error)
	}

	// Send a 201 created response
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Created")
}

func (h handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	// Read the dynamic id parameter
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Find the transaction by Id

	var transaction entity.Transaction

	if result := h.DB.First(&transaction, id); result.Error != nil {
		fmt.Println(result.Error)
	}

	// Delete that transaction
	h.DB.Delete(&transaction)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Deleted")
}

func (h handler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	var transactions []entity.Transaction

	if result := h.DB.Find(&transactions); result.Error != nil {
		fmt.Println(result.Error)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transactions)
}

func (h handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	// Read dynamic id parameter
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Find transaction by Id
	var transaction entity.Transaction

	if result := h.DB.First(&transaction, id); result.Error != nil {
		fmt.Println(result.Error)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}

func (h handler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	// Read dynamic id parameter
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Read request body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var updatedTransaction entity.Transaction
	json.Unmarshal(body, &updatedTransaction)

	var transaction entity.Transaction

	if result := h.DB.First(&transaction, id); result.Error != nil {
		fmt.Println(result.Error)
	}

	transaction.Name = updatedTransaction.Name
	transaction.Status = updatedTransaction.Status
	transaction.Amount = updatedTransaction.Amount

	h.DB.Save(&transaction)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Updated")
}

// handlers.go

type handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) handler {
	return handler{db}
}

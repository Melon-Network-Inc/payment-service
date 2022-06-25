package transaction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Melon-Network-Inc/payment-service/pkg/entity"
	"github.com/Melon-Network-Inc/payment-service/pkg/log"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterHandlers(r *mux.Router, service Service, db *gorm.DB, logger log.Logger) {
	h := handler{DB: db}
	res := resource{service, logger}

	routes := r.PathPrefix("/transactions/").Subrouter()
	routes.HandleFunc("/{id}", res.GetTransaction).Methods(http.MethodGet)
	routes.HandleFunc("/", h.GetAllTransactions).Methods(http.MethodGet)
	routes.HandleFunc("/", h.AddTransaction).Methods(http.MethodPost)
	routes.HandleFunc("/{id}", h.UpdateTransaction).Methods(http.MethodPut)
	routes.HandleFunc("/{id}", h.DeleteTransaction).Methods(http.MethodDelete)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (h handler) AddTransaction(writer http.ResponseWriter, response *http.Request) {
	// Read to request body
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	var transaction entity.Transaction
	json.Unmarshal(body, &transaction)

	// Append to the Transaction table
	if result := h.DB.Create(&transaction); result.Error != nil {
		fmt.Println(result.Error)
	}

	// Send a 201 created response
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode("Created")
}

func (h handler) DeleteTransaction(res http.ResponseWriter, req *http.Request) {
	// Read the dynamic id parameter
	vars := mux.Vars(req)
	id, _ := strconv.Atoi(vars["id"])

	// Find the transaction by Id

	var transaction entity.Transaction

	if result := h.DB.First(&transaction, id); result.Error != nil {
		fmt.Println(result.Error)
	}

	// Delete that transaction
	h.DB.Delete(&transaction)

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode("Deleted")
}

func (h handler) GetAllTransactions(res http.ResponseWriter, req *http.Request) {
	var transactions []entity.Transaction

	if result := h.DB.Find(&transactions); result.Error != nil {
		fmt.Println(result.Error)
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(transactions)
}

func (r resource) GetTransaction(res http.ResponseWriter, req *http.Request) {
	// Read dynamic id parameter
	vars := mux.Vars(req)
	id, _ := strconv.Atoi(vars["id"])

	// Find transaction by Id
	var transaction Transaction

	transaction, err := r.service.Get(req.Context(), id)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(transaction)
}

func (h handler) UpdateTransaction(res http.ResponseWriter, req *http.Request) {
	// Read dynamic id parameter
	vars := mux.Vars(req)
	id, _ := strconv.Atoi(vars["id"])

	// Read request body
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
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

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode("Updated")
}

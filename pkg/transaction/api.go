package transaction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Melon-Network-Inc/payment-service/pkg/log"
	"github.com/gorilla/mux"
)

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger) {
	res := resource{service, logger}

	routes := r.PathPrefix("/transactions/").Subrouter()
	routes.HandleFunc("/{id}", res.GetTransaction).Methods(http.MethodGet)
	routes.HandleFunc("/", res.GetAllTransactions).Methods(http.MethodGet)
	routes.HandleFunc("/", res.AddTransaction).Methods(http.MethodPost)
	routes.HandleFunc("/{id}", res.UpdateTransaction).Methods(http.MethodPut)
	routes.HandleFunc("/{id}", res.DeleteTransaction).Methods(http.MethodDelete)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) AddTransaction(res http.ResponseWriter, req *http.Request) {
	// Read to request body
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}

	var input AddTransactionRequest
	json.Unmarshal(body, &input)

	var transaction Transaction
	transaction, err = r.service.Add(req.Context(), input)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Send a 201 created response
	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(&transaction)
}

func (r resource) GetTransaction(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	transaction, err := r.service.Get(req.Context(), vars["id"])
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(transaction)
}

func (r resource) GetAllTransactions(res http.ResponseWriter, req *http.Request) {
	transactions, err := r.service.List(req.Context())
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(transactions)
}

func (r resource) UpdateTransaction(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	// Read request body
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	var input UpdateTransactionRequest
	json.Unmarshal(body, &input)

	r.service.Update(req.Context(), vars["id"], input)

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode("Updated")
}

func (r resource) DeleteTransaction(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	if err := r.service.Delete(req.Context(), vars["id"]); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode("Deleted")
}

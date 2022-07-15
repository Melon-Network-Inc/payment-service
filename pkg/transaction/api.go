package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/Melon-Network-Inc/payment-service/pkg/log"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// APITestCase represents the data needed to describe an API test case.
type APITestCase struct {
	Name         string
	Method, URL  string
	Body         string
	Header       http.Header
	WantStatus   int
	WantResponse string
}

// Endpoint tests an HTTP endpoint using the given APITestCase spec.
func Endpoint(t *testing.T, router *routing.Router, tc APITestCase) {
	t.Run(tc.Name, func(t *testing.T) {
		req, _ := http.NewRequest(tc.Method, tc.URL, bytes.NewBufferString(tc.Body))
		if tc.Header != nil {
			req.Header = tc.Header
		}
		res := httptest.NewRecorder()
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}
		router.ServeHTTP(res, req)
		assert.Equal(t, tc.WantStatus, res.Code, "status mismatch")
		if tc.WantResponse != "" {
			pattern := strings.Trim(tc.WantResponse, "*")
			if pattern != tc.WantResponse {
				assert.Contains(t, res.Body.String(), pattern, "response mismatch")
			} else {
				assert.JSONEq(t, tc.WantResponse, res.Body.String(), "response mismatch")
			}
		}
	})
}

type handler struct {
	DB *gorm.DB
}

func RegisterHandlers(r *mux.Router, service Service, logger log.Logger) {
	// h := handler{DB: db}
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

func (r resource) AddTransaction(writer http.ResponseWriter, response *http.Request) {

	// Read to request body
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	var addTransaction AddTransaction
	json.Unmarshal(body, &addTransaction)

	var transaction Transaction
	json.Unmarshal(body, &transaction)

	// Append to the Transaction table
	// if result := r.service.Create(&transaction); result.Error != nil { //not sure
	r.service.Add(response.Context(), addTransaction)

	// Send a 201 created response
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode("Created")
}

func (r resource) DeleteTransaction(res http.ResponseWriter, req *http.Request) {
	// Read dynamic id parameter
	vars := mux.Vars(req)
	id, _ := strconv.Atoi(vars["id"])

	var transaction Transaction

	// api layer calls service layer get to get transaction
	transaction, err := r.service.Get(req.Context(), id) // fill in transaction table (service)

	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	// Delete that transaction
	r.service.Delete(req.Context(), id)

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(transaction)
}

// need service layer GetAll first
func (r resource) GetAllTransactions(res http.ResponseWriter, req *http.Request) {
	// Read dynamic id parameter
	vars := mux.Vars(req)
	id, _ := strconv.Atoi(vars["id"])

	var transaction Transaction
	// api layer calls service layer get to get transaction
	transaction, err := r.service.Get(req.Context(), id) // fill in transaction table (service)

	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(transaction)
}

func (r resource) GetTransaction(res http.ResponseWriter, req *http.Request) {
	// Read dynamic id parameter
	vars := mux.Vars(req)
	id, _ := strconv.Atoi(vars["id"])

	// Find transaction by Id
	var transaction Transaction

	// api layer calls service layer get to get transaction
	transaction, err := r.service.Get(req.Context(), id) // fill in transaction table (service)
	// use r.service layer (dont need handler class h.DB)
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(transaction)
}

func (r resource) UpdateTransaction(res http.ResponseWriter, req *http.Request) {
	// Read dynamic id parameter
	vars := mux.Vars(req)
	id, _ := strconv.Atoi(vars["id"])

	// Read request body
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
	}

	var updatedTransaction UpdateTransactionRequest
	json.Unmarshal(body, &updatedTransaction)

	var transaction Transaction

	// api layer calls service layer get to get transaction
	transaction, err = r.service.Get(req.Context(), id)
	if err != nil { // fill in transaction table (service)
		// use r.service layer (dont need handler class h.DB)
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	ID, err := strconv.Atoi(updatedTransaction.Id)
	transaction.Id = uint(ID)
	transaction.Name = updatedTransaction.Name
	transaction.Message = updatedTransaction.Message
	transaction.SenderPubkey = updatedTransaction.SenderPubkey

	string_id := strconv.Itoa(id)
	r.service.Update(req.Context(), string_id, updatedTransaction) //has error

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(transaction)
}

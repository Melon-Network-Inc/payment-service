package transaction

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Melon-Network-Inc/payment-service/pkg/entity"
	"github.com/Melon-Network-Inc/payment-service/pkg/log"
	"gopkg.in/go-playground/validator.v9"
)

// Service encapsulates usecase logic for transactions.
type Service interface {
	Add(ctx context.Context, input AddTransactionRequest) (Transaction, error)
	Get(c context.Context, id string) (Transaction, error)
	List(ctx context.Context) ([]Transaction, error)
	Update(ctx context.Context, id string, input UpdateTransactionRequest) (Transaction, error)
	Delete(ctx context.Context, id string) error
}

// address represents the data about an address.
type Transaction struct {
	entity.Transaction
}

//July 7
// AddAddressRequest represents an address creation request.
type AddTransactionRequest struct {
	Name           string `json:"name"        validate:"required"`
	Status         string `json:"status"`
	Amount         uint   `json:"amount"      validate:"required,uint"`
	Currency       string `json:"currency"    validate:"required,iso4217"` //currency code
	SenderId       uint64 `json:"sender_id"   validate:"uuid"`
	SenderPubkey   uint64 `json:"sender_pk"   validate:"required, oneof='eth_addr' 'btc_addr'"` // ETH or BTC address
	ReceiverId     uint64 `json:"receiver_id" validate:"uuid"`
	ReceiverPubkey uint64 `json:"receiver_pk" validate:"required, oneof='eth_addr' 'btc_addr'"` // ETH or BTC address
	Message        string `json:"message"     validate:"ls=200"`
}

//July 7
// Validate validates the AddTransaction fields.
func (m AddTransactionRequest) Validate() error {
	validate := validator.New()
	err := validate.StructExcept(m, "Status")
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

//July 7
// UpdateTransactionRequest represents an address update request.
type UpdateTransactionRequest struct {
	Id      uint   `json:"id"`
	Name    string `json:"name"    validate:"required"`
	Message string `json:"message" validate:"ls=200"`
	Status  string `json:"status"`
}

//July 7
// Validate validates the UpdateTransactionRequest fields.
func (m UpdateTransactionRequest) Validate() error {
	validate := validator.New()
	err := validate.StructExcept(m, "Status")
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new transaction service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Create creates a new transaction.
func (s service) Add(ctx context.Context, req AddTransactionRequest) (Transaction, error) {
	if err := req.Validate(); err != nil {
		return Transaction{}, err
	}

	transaction, err := s.repo.Add(ctx, entity.Transaction{
		Name:           req.Name,
		Status:         req.Status,
		Amount:         req.Amount,
		Currency:       req.Currency,
		SenderId:       req.SenderId,
		SenderPubkey:   req.SenderPubkey,
		ReceiverId:     req.ReceiverId,
		ReceiverPubkey: req.ReceiverPubkey,
		Message:        req.Message,
	})
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

// Get returns the transaction with the specified the transaction ID.
func (s service) Get(ctx context.Context, id string) (Transaction, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return Transaction{}, err
	}
	transaction, err := s.repo.Get(ctx, uid)
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

// Get returns the a list of transactions associated to a user.
func (s service) List(ctx context.Context) ([]Transaction, error) {
	transaction, err := s.repo.List(ctx)
	if err != nil {
		return []Transaction{}, err
	}
	listTransaction := []Transaction{}
	for _, transaction := range transaction {
		listTransaction = append(listTransaction, Transaction{transaction})
	}
	return listTransaction, nil
}

// Update updates the transaction with the specified the transaction ID.
func (s service) Update(
	ctx context.Context,
	id string,
	input UpdateTransactionRequest,
) (Transaction, error) {
	if err := input.Validate(); err != nil {
		return Transaction{}, err
	}

	uid, err := strconv.Atoi(id)
	if err != nil {
		return Transaction{}, err
	}

	transaction, err := s.repo.Get(ctx, uid)
	if err != nil {
		return Transaction{}, err
	}
	transaction.Name = input.Name
	transaction.Status = input.Status
	transaction.Message = input.Message

	if err := s.repo.Update(ctx, transaction); err != nil {
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

// Delete deletes the transaction with the specified ID.
func (s service) Delete(ctx context.Context, id string) error {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, uid)
}

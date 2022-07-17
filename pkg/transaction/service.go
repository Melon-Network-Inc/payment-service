package transaction

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Melon-Network-Inc/payment-service/pkg/entity"
	"github.com/Melon-Network-Inc/payment-service/pkg/log"
	"gopkg.in/go-playground/validator.v9"
)

// Service encapsulates usecase logic for transactions.
type Service interface {
	Get(c context.Context, id int) (Transaction, error)
	Update(ctx context.Context, id string, input UpdateTransactionRequest) (Transaction, error)
	List(ctx context.Context) ([]Transaction, error)
	Delete(ctx context.Context, id int) (Transaction, error)
}

// address represents the data about an address.
type Transaction struct {
	entity.Transaction
}

//July 7
// AddAddressRequest represents an address creation request.
type AddTransaction struct {
	Id             uint   `json:"id"` // string of hex? then use "hexadecimal"
	Name           string `json:"name" validate:"required"`
	Status         string
	Amount         uint      `json:"amount" validate:"required,uint"`
	Currency       string    `json:"currency" validate:"required,iso4217"` //currency code
	SenderId       uint64    `json:"sender_id" validate:"uuid"`
	SenderPubkey   uint64    `json:"sender_pk" validate:"required, oneof='eth_addr' 'btc_addr'"` // ETH or BTC address
	ReceiverId     uint64    `json:"receiver_id" validate:"uuid"`
	ReceiverPubkey uint64    `json:"receiver_pk" validate:"required, oneof='eth_addr' 'btc_addr'"` // ETH or BTC address
	CreatedAt      time.Time `json:"creat_at" validate:"required,datetime"`
	UpdatedAt      time.Time `json:"update_at" validate:"required,datetime"`
	// message should be less than 200 characters
	Message string `json:"message" validate:"ls=200"`
}

//July 7
// Validate validates the AddTransaction fields.
func (m AddTransaction) Validate() error {
	// return validation.ValidateStruct(&m,
	// 	validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	// )

	// run `go get gopkg.in/go-playground/validator.v9`
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
	Id           uint   `json:"id"`
	Name         string `json:"name" validate:"required"`
	Message      string `json:"message" validate:"ls=200"`
	SenderPubkey uint64 `json:"sender_pk" validate:"required, oneof='eth_addr' 'btc_addr'"` // ETH or BTC address
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

// Get returns the transaction with the specified the transaction ID.
func (s service) Get(ctx context.Context, id int) (Transaction, error) {
	transaction, err := s.repo.Get(ctx, id)
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

// Create creates a new transaction.
func (s service) Add(ctx context.Context, req AddTransaction) (Transaction, error) {
	if err := req.Validate(); err != nil {
		return Transaction{}, err
	}

	err := s.repo.Add(ctx, entity.Transaction{
		Id:             req.Id,
		Name:           req.Name,
		Status:         req.Status,
		Amount:         req.Amount,
		Currency:       req.Currency,
		SenderId:       req.SenderId,
		SenderPubkey:   req.SenderPubkey,
		ReceiverId:     req.ReceiverId,
		ReceiverPubkey: req.ReceiverPubkey,
		CreatedAt:      req.CreatedAt,
		UpdatedAt:      req.UpdatedAt,
		Message:        req.Message,
	})
	if err != nil {
		return Transaction{}, err
	}

	return s.Get(ctx, int(req.Id))
}

// July 7
func (s service) Update(ctx context.Context, id string, input UpdateTransactionRequest) (Transaction, error) {
	if err := input.Validate(); err != nil {
		return Transaction{}, err
	}

	uid, err := strconv.Atoi(id)

	transaction, err := s.repo.Get(ctx, uid)
	if err != nil {
		return Transaction{}, err
	}
	transaction.Name = input.Name
	transaction.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, transaction); err != nil {
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

// Get returns the transaction with the specified the transaction ID.
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

// Delete deletes the transaction with the specified ID.
func (s service) Delete(ctx context.Context, id int) (Transaction, error) {
	transaction, err := s.repo.Get(ctx, id)
	if err != nil {
		return Transaction{}, err
	}

	if err = s.repo.Delete(ctx, id); err != nil {
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

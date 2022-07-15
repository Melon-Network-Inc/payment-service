package transaction

import (
	"context"
	"strconv"
	"time"

	"github.com/Melon-Network-Inc/payment-service/pkg/entity"
	"github.com/Melon-Network-Inc/payment-service/pkg/log"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Service encapsulates usecase logic for transactions.
type Service interface {
	Add(ctx context.Context, req AddTransaction) (Transaction, error)
	Get(c context.Context, id int) (Transaction, error)
	Update(ctx context.Context, id string, input UpdateTransaction) (Transaction, error)
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
	Id             uint
	Name           string
	Status         string
	Amount         uint
	Currency       string
	SenderId       uint64
	SenderPubkey   uint64
	ReceiverId     uint64
	ReceiverPubkey uint64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Message        string
}

//July 7
// Validate validates the CreateAddressRequest fields.
func (m AddTransaction) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

//July 7
// UpdateAddressRequest represents an address update request.
type UpdateTransaction struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Message      string // 控制字符200 character
	SenderPubkey uint64 //public key有const：看metamask PK（改成string of hex）
}

//July 7
// Validate validates the CreateAddressRequest fields.
func (m UpdateTransaction) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new address service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the address with the specified the address ID.
func (s service) Get(ctx context.Context, id int) (Transaction, error) {
	transaction, err := s.repo.Get(ctx, id)
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

// Create creates a new address.
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

// need service update
// July 7
func (s service) Update(ctx context.Context, id string, input UpdateTransaction) (Transaction, error) {
	if err := input.Validate(); err != nil {
		return Transaction{}, err
	}

	uid, err := strconv.Atoi(id)

	transaction, err := s.repo.Get(ctx, uid)
	if err != nil {
		return Transaction{}, err
	}
	transaction.Name = input.Name
	// transaction.IsPrimary = req.IsPrimary
	transaction.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, transaction); err != nil { // not sure Update takes in transaction or Transaction{t}
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

// add, delete, List/GetAll

// Get returns the address with the specified the address ID.
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

// Delete deletes the address with the specified ID.
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

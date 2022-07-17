package transaction

import (
	"context"

	"github.com/Melon-Network-Inc/payment-service/pkg/entity"
	"github.com/Melon-Network-Inc/payment-service/pkg/log"
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

// Validate validates the CreateAddressRequest fields.
func (m AddTransaction) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

// UpdateTransactionRequest represents an address update request.
type UpdateTransactionRequest struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Message      string // 控制字符200 character
	SenderPubkey uint64 //public key有const：看metamask PK（改成string of hex）
}

// Validate validates the CreateAddressRequest fields.
func (m UpdateTransactionRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
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
	// transaction.IsPrimary = req.IsPrimary
	transaction.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, transaction); err != nil { // not sure Update takes in transaction or Transaction{t}
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

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
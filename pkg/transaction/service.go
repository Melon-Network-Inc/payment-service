package transaction

import (
	"context"
	"strconv"

	"github.com/Melon-Network-Inc/entity-repo/pkg/api"
	"github.com/Melon-Network-Inc/entity-repo/pkg/entity"
	"github.com/Melon-Network-Inc/payment-service/pkg/log"
)

// Service encapsulates usecase logic for transactions.
type Service interface {
	Add(ctx context.Context, input api.AddTransactionRequest) (Transaction, error)
	Get(c context.Context, id string) (Transaction, error)
	List(ctx context.Context) ([]Transaction, error)
	Update(ctx context.Context, id string, input api.UpdateTransactionRequest) (Transaction, error)
	Delete(ctx context.Context, id string) (Transaction, error)
}

// transaction represents the data about a transaction.
type Transaction struct {
	entity.Transaction
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
func (s service) Add(ctx context.Context, req api.AddTransactionRequest) (Transaction, error) {
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
	input api.UpdateTransactionRequest,
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
func (s service) Delete(ctx context.Context, id string) (Transaction, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return Transaction{}, err
	}
	transaction, err := s.repo.Delete(ctx, uid)
	if err != nil {
		return Transaction{}, err
	}
	return Transaction{transaction}, nil
}

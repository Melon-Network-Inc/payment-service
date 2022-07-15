package transaction

import (
	"context"

	"github.com/Melon-Network-Inc/payment-service/pkg/entity"
	"github.com/Melon-Network-Inc/payment-service/pkg/log"
)

// Service encapsulates usecase logic for transactions.
type Service interface {
	Get(c context.Context, id int) (Transaction, error)
}

// address represents the data about an address.
type Transaction struct {
	entity.Transaction
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

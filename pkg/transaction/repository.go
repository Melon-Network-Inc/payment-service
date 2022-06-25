package transaction

import (
	"context"

	db "github.com/Melon-Network-Inc/payment-service/pkg/dbcontext"
	"github.com/Melon-Network-Inc/payment-service/pkg/entity"
	"github.com/Melon-Network-Inc/payment-service/pkg/log"
)

// Repository encapsulates the logic to access transactions from the data source.
type Repository interface {
	// Get returns the address with the specified address ID.
	Get(c context.Context, id int) (entity.Transaction, error)
}

// repository persists addresses in database
type repository struct {
	db     *db.DB
	logger log.Logger
}

// NewRepository creates a new transaction repository
func NewRepository(db *db.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the transaction with the specified ID from the database.
func (r repository) Get(c context.Context, id int) (entity.Transaction, error) {
	var address entity.Transaction
	result := r.db.With(c).First(&address, id)
	return address, result.Error
}

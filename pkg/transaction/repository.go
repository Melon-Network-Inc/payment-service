package transaction

import (
	"context"

	"github.com/Melon-Network-Inc/entity-repo/pkg/entity"
	db "github.com/Melon-Network-Inc/payment-service/pkg/dbcontext"
	"github.com/Melon-Network-Inc/payment-service/pkg/log"
)

// Repository encapsulates the logic to access transactions from the data source.
type Repository interface {
	// Add creates the transaction.
	Add(c context.Context, transaction entity.Transaction) (entity.Transaction, error)
	// Get returns the transaction with the specified transaction ID.
	Get(c context.Context, id int) (entity.Transaction, error)
	// List returns the transaction associated to target user.
	List(c context.Context) ([]entity.Transaction, error)
	// Update returns the transaction with the specified transaction ID.
	Update(c context.Context, transaction entity.Transaction) error
	// Delete deletes the transaction with the specified ID.
	Delete(c context.Context, id int) (entity.Transaction, error)
}

// repository persists transactions in database
type repository struct {
	db     *db.DB
	logger log.Logger
}

// NewRepository creates a new transaction repository
func NewRepository(db *db.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Add creates the transaction.
func (r repository) Add(
	c context.Context,
	transaction entity.Transaction,
) (entity.Transaction, error) {
	if result := r.db.With(c).Create(&transaction); result.Error != nil {
		return entity.Transaction{}, result.Error
	}
	return r.Get(c, int(transaction.ID))
}

// Get reads the transaction with the specified ID from the database.
func (r repository) Get(c context.Context, id int) (entity.Transaction, error) {
	var transaction entity.Transaction
	result := r.db.With(c).First(&transaction, id)
	return transaction, result.Error
}

// Lists lists all transactions.
func (r repository) List(c context.Context) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	result := r.db.With(c).Find(&transactions, entity.Transaction{})
	return transactions, result.Error
}

// Update updates the transaction with the specified transaction ID.
func (r repository) Update(c context.Context, transaction entity.Transaction) error {
	if result := r.db.With(c).First(&transaction, transaction.ID); result.Error != nil {
		return result.Error
	}
	return r.db.With(c).Save(&transaction).Error
}

// Delete deletes the transaction with the specified ID.
func (r repository) Delete(c context.Context, id int) (entity.Transaction, error) {
	transaction, err := r.Get(c, id)
	if err != nil {
		return entity.Transaction{}, err
	}
	return transaction, r.db.With(c).Delete(&transaction).Error
}

package transaction

import (
	"context"

	"github.com/Melon-Network-Inc/payment-service/pkg/dbcontext"
	"github.com/Melon-Network-Inc/payment-service/pkg/entity"
	"github.com/Melon-Network-Inc/payment-service/pkg/log"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Repository encapsulates the logic to access transactions from the data source.
type Repository interface {
	Add(ctx context.Context, transaction entity.Transaction) error
	// Get returns the transaction with the specified transaction ID.
	Get(c context.Context, id int) (entity.Transaction, error)
	// Update updates the transaction with given ID in the storage.
	Update(ctx context.Context, transaction entity.Transaction) error
	List(ctx context.Context) ([]entity.Transaction, error)
	Delete(ctx context.Context, id int) error
}

// repository persists transaction in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new transaction repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the transaction with the specified ID from the database.
func (r repository) Get(c context.Context, id int) (entity.Transaction, error) {
	var transaction entity.Transaction
	result := r.db.With(c).First(&transaction, id)
	return transaction, result.Error
}

// July 7 need repository update delete add put *****put?
//also in api and service

// Create saves a new transaction record in the database.
// It returns the ID of the newly inserted transaction record.
func (r repository) Add(ctx context.Context, transaction entity.Transaction) error {
	return r.db.With(ctx).Create(&transaction).Error
}

// Update saves the changes to a transaction in the database.
func (r repository) Update(ctx context.Context, transaction entity.Transaction) error {
	if result := r.db.With(ctx).First(&transaction, transaction.Id); result.Error != nil {
		return result.Error
	}
	return r.db.With(ctx).Save(&transaction).Error
}

// Delete deletes a transaction with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id int) error {
	transaction, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Delete(&transaction).Error
}

// List returns all addresses owned by the user.
func (r repository) List(ctx context.Context) ([]entity.Transaction, error) {
	owner := ctx.Value("user")
	var transaction []entity.Transaction
	// err := r.db.With(ctx).Select().Where(dbx.HashExp{"owner": owner}).All(transaction) //not sure
	err := r.db.With(ctx).Find(dbx.HashExp{"owner": owner}, &transaction).Error
	return transaction, err
}

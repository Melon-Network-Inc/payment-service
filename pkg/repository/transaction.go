package repository

import (
	db "github.com/Melon-Network-Inc/common/pkg/dbcontext"
	"github.com/Melon-Network-Inc/common/pkg/entity"
	"github.com/Melon-Network-Inc/common/pkg/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TransactionRepository encapsulates the logic to access transactions from the data source.
type TransactionRepository interface {
	// Add creates the transaction.
	Add(c *gin.Context, transaction entity.Transaction) (entity.Transaction, error)
	// Get returns the transaction with the specified transaction ID.
	Get(c *gin.Context, ID int) (entity.Transaction, error)
	// List returns the transaction associated to target user.
	List(c *gin.Context, ID int, showType string) ([]entity.Transaction, error)
	// Update returns the transaction with the specified transaction ID.
	Update(c *gin.Context, transaction entity.Transaction) error
	// Delete deletes the transaction.
	Delete(c *gin.Context, transaction entity.Transaction) error
}

// transactionRepository persists transactions in database
type transactionRepository struct {
	db     *db.DB
	logger log.Logger
}

// NewTransactionRepository creates a new transaction transactionRepository
func NewTransactionRepository(db *db.DB, logger log.Logger) TransactionRepository {
	return transactionRepository{db, logger}
}

// Add creates the transaction.
func (r transactionRepository) Add(
	c *gin.Context,
	transaction entity.Transaction,
) (entity.Transaction, error) {
	if result := r.db.With(c).Create(&transaction); result.Error != nil {
		return entity.Transaction{}, result.Error
	}
	return r.Get(c, int(transaction.ID))
}

// Get reads the transaction with the specified ID from the database.
func (r transactionRepository) Get(c *gin.Context, ID int) (entity.Transaction, error) {
	var transaction entity.Transaction
	result := r.db.With(c).First(&transaction, ID)
	return transaction, result.Error
}

// List lists all transactions by show_type.
func (r transactionRepository) List(c *gin.Context, ID int, showType string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	var result *gorm.DB
	tx := r.db.With(c).
		Where("sender_id = ?", ID).
		Or("receiver_id = ?", ID).
		Order("updated_at desc")

	tx = updateTransactionByShowType(showType, tx, transactions)
	result = tx.Find(&transactions)
	return transactions, result.Error
}

func updateTransactionByShowType(showType string, tx *gorm.DB, transactions []entity.Transaction) *gorm.DB {
	if showType == "Private" {
		tx = tx.Find(&transactions)
	} else if showType == "Friend" {
		tx = tx.Where("show_type = ?", "Friend").
			Or("show_type = ?", "Public")
	} else {
		tx = tx.Where("show_type = ?", "Public")
	}
	return tx
}

// Update updates the transaction with the specified transaction ID.
func (r transactionRepository) Update(c *gin.Context, transaction entity.Transaction) error {
	if result := r.db.With(c).First(&transaction, transaction.ID); result.Error != nil {
		return result.Error
	}
	return r.db.With(c).Save(&transaction).Error
}

// Delete deletes the transaction with the specified ID.
func (r transactionRepository) Delete(c *gin.Context, transaction entity.Transaction) error {
	return r.db.With(c).Delete(&transaction).Error
}
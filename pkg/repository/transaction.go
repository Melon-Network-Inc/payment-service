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
	Get(c *gin.Context, userID uint) (entity.Transaction, error)
	// List returns the transaction associated to target user.
	List(c *gin.Context, userID uint, showType string) ([]entity.Transaction, error)
	// Update returns the transaction with the specified transaction ID.
	Update(c *gin.Context, transaction entity.Transaction) error
	// Delete deletes the transaction.
	Delete(c *gin.Context, transaction entity.Transaction) error
	// Count returns the number of user's transactions in the database.
	Count(ctx *gin.Context, userID uint, showType string) (int, error)
	// Query returns the list of user's transactions with the given offset and limit.
	Query(ctx *gin.Context, offset, limit int, userID uint, showType string) ([]entity.Transaction, error)
	// ListByRequester returns the transaction associated to requester.
	ListByRequester(c *gin.Context, requesterID uint) ([]entity.Transaction, error)
	// ListByUserID returns the transaction associated to target user.
	ListByUserID(c *gin.Context, requesterID, ID uint) ([]entity.Transaction, error)
	// CountByFriendIDs returns the number of user's transactions in the database.
	CountByFriendIDs(ctx *gin.Context, requesterID uint, friendsIDs []uint) (int, error)
	// QueryByFriendIDs returns the list of user's transactions with the given offset and limit by the friend IDs.
	QueryByFriendIDs(ctx *gin.Context, offset, limit int, requesterID uint, friendsIDs []uint) ([]entity.Transaction, error)
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
	return r.Get(c, transaction.ID)
}

// Get reads the transaction with the specified ID from the database.
func (r transactionRepository) Get(c *gin.Context, userID uint) (entity.Transaction, error) {
	var transaction entity.Transaction
	result := r.db.With(c).First(&transaction, userID)
	return transaction, result.Error
}

// List lists all transactions by show_type.
func (r transactionRepository) List(c *gin.Context, userID uint, showType string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	var result *gorm.DB
	tx := createTransactionByShowType(showType, r.db.With(c)).
		Where("sender_id = ?", userID).
		Or("receiver_id = ?", userID).
		Order("updated_at desc")

	result = tx.Find(&transactions)
	return transactions, result.Error
}

// Count returns the number of user's transactions in the database.
func (r transactionRepository) Count(ctx *gin.Context, userID uint, showType string) (int, error) {
	var rows int64
	result := r.db.With(ctx).Model(&entity.Transaction{}).
		Where("sender_id = ?", userID).
		Or("receiver_id = ?", userID).
		Order("updated_at desc").
		Count(&rows)
	return int(rows), result.Error
}

// Query returns the list of user's transactions with the given offset and limit.
func (r transactionRepository) Query(ctx *gin.Context, offset, limit int, userID uint, showType string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	tx := r.db.With(ctx).Model(&entity.Transaction{}).
		Where("sender_id = ?", userID).
		Or("receiver_id = ?", userID).
		Order("updated_at desc").
		Offset(offset).
		Limit(limit)
	tx = createTransactionByShowType(showType, tx)
	result := tx.Find(&transactions)
	return transactions, result.Error
}

// createTransactionByShowType updates the transaction with the specified show type.
func createTransactionByShowType(showType string, tx *gorm.DB) *gorm.DB {
	if showType == "Public" {
		tx = tx.Where("show_type = ?", "Public")
	} else if showType == "Friend" {
		tx = tx.Where("show_type = ?", "Friend").
			Or("show_type = ?", "Public")
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

// CountByFriendIDs returns the number of user's transactions in the database by the friend IDs.
func (r transactionRepository) CountByFriendIDs(ctx *gin.Context, requesterID uint, friendsIDs []uint) (int, error) {
	var rows int64
	result := r.db.With(ctx).Model(&entity.Transaction{}).
		Where("sender_id = ? or receiver_id = ?",
			requesterID,
			requesterID).
		Or("(sender_id in ? or receiver_id in ?) and show_type in ?",
			friendsIDs,
			friendsIDs,
			[]string{"Friend", "Public"}).
		Where("deleted_at is null").
		Count(&rows)
	return int(rows), result.Error
}

// QueryByFriendIDs returns the list of user's transactions with the given offset and limit by the friend IDs.
func (r transactionRepository) QueryByFriendIDs(ctx *gin.Context, offset, limit int, requesterID uint, friendsIDs []uint) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	tx := r.db.With(ctx).Model(&entity.Transaction{}).
		Where("sender_id = ? or receiver_id = ?",
			requesterID,
			requesterID).
		Or("(sender_id in ? or receiver_id in ?) and show_type in ?",
			friendsIDs,
			friendsIDs,
			[]string{"Friend", "Public"}).
		Where("deleted_at is null").
		Order("updated_at desc").
		Offset(offset).
		Limit(limit)
	result := tx.Find(&transactions)
	return transactions, result.Error
}

// ListByRequester returns the transaction associated to target user.
func (r transactionRepository) ListByRequester(c *gin.Context, requesterID uint) ([]entity.Transaction, error) {
	var pastTransactions []entity.Transaction
	result := r.db.With(c).
		Where("sender_id = ? OR receiver_id = ?", requesterID, requesterID).
		Find(&pastTransactions)
	if result.Error != nil {
		return []entity.Transaction{}, result.Error
	}

	return pastTransactions, nil
}

// ListByUserID returns the transaction associated to target user.
func (r transactionRepository) ListByUserID(c *gin.Context, requesterID uint, ID uint) ([]entity.Transaction, error) {
	var pastTransactions []entity.Transaction
	result := r.db.With(c).
		Where("sender_id = ? OR receiver_id = ?", ID, ID).
		Where("sender_id != ? AND receiver_id != ?", requesterID, requesterID).
		Where("show_type in ?", []string{"Friend", "Public"}).
		Find(&pastTransactions)
	if result.Error != nil {
		return []entity.Transaction{}, result.Error
	}

	return pastTransactions, nil
}

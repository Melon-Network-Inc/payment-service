package repository

import (
	"time"

	db "github.com/Melon-Network-Inc/common/pkg/dbcontext"
	"github.com/Melon-Network-Inc/common/pkg/entity"
	"github.com/Melon-Network-Inc/common/pkg/log"
	"github.com/gin-gonic/gin"
)

// NewsRepository encapsulates the logic to access news from the data source.
type NewsRepository interface {
	// BatchAdd creates an array of news.
	BatchAdd(c *gin.Context, news []entity.News) (int, error)
	// Delete deletes the news.
	Delete(c *gin.Context, news entity.News) error
	// DeleteBefore deletes the news by timestamp.
	DeleteBefore(c *gin.Context, deadline time.Time) error
	// Count returns the number of user's news in the database.
	Count(ctx *gin.Context) (int, error)
	// Query returns the list of news with the given offset and limit.
	Query(ctx *gin.Context, offset, limit int) ([]entity.News, error)
}

// transactionRepository persists transactions in database
type newsRepository struct {
	db     *db.DB
	logger log.Logger
}

// NewTransactionRepository creates a new transaction transactionRepository
func NewNewsRepository(db *db.DB, logger log.Logger) NewsRepository {
	return newsRepository{db, logger}
}

// Add creates an array of news.
func (r newsRepository) BatchAdd(c *gin.Context, news []entity.News) (int, error) {
	if result := r.db.With(c).CreateInBatches(&news, len(news)); result.Error != nil {
		return 0, result.Error
	}
	return r.Count(c)
}

// Delete deletes the news by the news entity.
func (r newsRepository) Delete(c *gin.Context, news entity.News) error {
	return r.db.With(c).Delete(&news).Error
}

// DeleteBefore deletes the news before certain timestamp.
func (r newsRepository) DeleteBefore(c *gin.Context, timestamp time.Time) error{
	return r.db.With(c).Model(&entity.News{}).Where("created_at < ?", timestamp).Delete(&entity.News{}).Error
}

// Count returns the number of news in the database by the friend IDs.
func (r newsRepository) Count(ctx *gin.Context) (int, error) {
	var rows int64
	result := r.db.With(ctx).Model(&entity.News{}).
		Count(&rows)
	return int(rows), result.Error
}

// Query returns the list of news with the given offset and limit by the friend IDs.
func (r newsRepository) Query(ctx *gin.Context, offset, limit int) ([]entity.News, error) {
	var news []entity.News
	tx := r.db.With(ctx).Model(&entity.News{}).
		Order("id desc").
		Offset(offset).
		Limit(limit)
	result := tx.Find(&news)
	return news, result.Error
}
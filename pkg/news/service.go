package news

import (
	"fmt"
	"github.com/Melon-Network-Inc/account-service/pkg/processor"
	"github.com/Melon-Network-Inc/common/pkg/entity"
	"github.com/Melon-Network-Inc/common/pkg/log"
	"github.com/Melon-Network-Inc/payment-service/pkg/repository"
	"github.com/gin-gonic/gin"
)

// Service encapsulates use case logic for activities.
type Service interface {
	// Count returns the number of news in database.
	Count(c *gin.Context) (int, error)
	// Query returns news by offset and limit.
	Query(c *gin.Context, offset, limit int) ([]entity.News, error)
	// Collect fetches the latest urls from source urls and store them into database.
	Collect()
}

type service struct {
	newsRepo   repository.NewsRepository
	newsClient Client
	logger     log.Logger
}

// NewService creates a new address service.
func NewService(
	newsRepo repository.NewsRepository,
	newsClient Client,
	logger log.Logger) Service {
	return service{newsRepo, newsClient, logger}
}

// Count returns the number of news in database.
func (s service) Count(c *gin.Context) (int, error) {
	_, err := processor.GetContextUserID(c)
	if err != nil {
		return 0, err
	}

	return s.newsRepo.Count(c)
}

// Query returns news by offset and limit.
func (s service) Query(c *gin.Context, offset, limit int) ([]entity.News, error) {
	items, err := s.newsRepo.Query(c, offset, limit)
	if err != nil {
		return []entity.News{}, err
	}
	return items, nil
}

// Collect returns the news from source urls
func (s service) Collect() {
	result := s.newsClient.FetchData()
	for _, newsItem := range result.NewsItems {
		fmt.Println(newsItem)
	}
}

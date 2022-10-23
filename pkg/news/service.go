package news

import (
	"fmt"
	"github.com/Melon-Network-Inc/account-service/pkg/processor"
	"github.com/Melon-Network-Inc/common/pkg/entity"
	"github.com/Melon-Network-Inc/common/pkg/log"
	"github.com/Melon-Network-Inc/payment-service/pkg/repository"
	"github.com/badoux/goscraper"
	"github.com/gin-gonic/gin"
	"time"
)

// Service encapsulates use case logic for activities.
type Service interface {
	// Count returns the number of news in database.
	Count(c *gin.Context) (int, error)
	// Query returns news by offset and limit.
	Query(c *gin.Context, offset, limit int) ([]entity.News, error)
	// Collect fetches the latest urls from source urls and store them into database.
	Collect()
	// InitializeNewsTable initialises the table if no news record found.
	InitializeNewsTable()
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

// InitializeNewsTable initialises news table if no news record found.
func (s service) InitializeNewsTable() {
	count, err := s.newsRepo.CountWithoutContext()
	if err != nil {
		s.logger.Error("initialises news fails with error ", err)
		return
	}
	if count == 0 {
		s.Collect()
	}
}

// Collect returns the news from source urls.
func (s service) Collect() {
	result := s.newsClient.FetchData()
	var fetchedNews []entity.News
	for _, newsItem := range result.NewsItems {
		item, err := goscraper.Scrape(newsItem.Url, 5)
		if err != nil {
			s.logger.Error("collecting news fails with error ", err)
		} else {
			fetchedNews = append(fetchedNews, entity.News{
				Title:              newsItem.Title,
				Url:                newsItem.Url,
				Source:             newsItem.Source,
				PreviewIcon:        item.Preview.Icon,
				PreviewName:        item.Preview.Name,
				PreviewImage:       item.Preview.Images[0],
				PreviewDescription: item.Preview.Description,
			})
		}
	}

	if len(fetchedNews) == 0 {
		s.logger.Error(fmt.Sprintf("fetch 0 records of news"))
		return
	}

	deletionDeadline := time.Now()

	count, err := s.newsRepo.BatchAdd(fetchedNews)
	if err != nil {
		s.logger.Error("collecting news fails with error ", err)
		return
	}

	newRecordCount, err := s.newsRepo.DeleteBefore(deletionDeadline)
	if err != nil {
		s.logger.Error("news deletion fails with error ", err)
		return
	}
	s.logger.Info(fmt.Sprintf(
		"fetch %d records of news and remove %d news records from previous run",
		newRecordCount,
		count-newRecordCount))
}

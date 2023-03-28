package news

import (
	"fmt"
	"time"

	"github.com/Melon-Network-Inc/common/pkg/entity"
	"github.com/Melon-Network-Inc/common/pkg/log"
	"github.com/badoux/goscraper"
)

type Consumer interface {
	InitializeNewsTable()
	Collect()
}


type consumer struct {
	service Service
	logger  log.Logger
}

// NewConsumer creates a new consumer.
func NewConsumer(service Service, logger log.Logger) Consumer {
	return consumer{
		service: service,
		logger:  logger,
	}
}

// InitializeNewsTable initialises news table if no news record found.
func (c consumer) InitializeNewsTable() {
	count, err := c.service.GetRepo().CountWithoutContext()
	if err != nil {
		c.logger.Error("initialises news fails with error ", err)
		return
	}
	if count == 0 {
		c.Collect()
	}
}

// Collect returns the news from source urls.
func (c consumer) Collect() {
	result := c.service.GetClient().FetchData()
	var fetchedNews []entity.News
	for _, newsItem := range result.NewsItems {
		item, err := goscraper.Scrape(newsItem.Url, 5)
		if err != nil {
			c.logger.Error("collecting news fails with error ", err)
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
		c.logger.Error("fetch 0 records of news")
		return
	}

	deletionDeadline := time.Now()

	count, err := c.service.GetRepo().BatchAdd(fetchedNews)
	if err != nil {
		c.logger.Error("collecting news fails with error ", err)
		return
	}

	newRecordCount, err := c.service.GetRepo().DeleteBefore(deletionDeadline)
	if err != nil {
		c.logger.Error("news deletion fails with error ", err)
		return
	}
	c.logger.Info(fmt.Sprintf(
		"fetch %d records of news and remove %d news records from previous run",
		newRecordCount,
		count-newRecordCount))
}

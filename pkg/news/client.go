package news

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Melon-Network-Inc/common/pkg/log"
)

type Client interface {
	FetchData() Result
	rapidAPIFetch(source Source) ([]Item, error)
	parseRegularFetch(source Source) (Result, error)
	parseCustomLinkFetch(source Source) (Result, error)
}

type newsClient struct {
	newsSource []Source
	logger     log.Logger
}

func NewClient(logger log.Logger) Client {
	var newsSources []Source
	newsSources = append(newsSources, CreateCryptoNewsSources())

	return newsClient{
		newsSource: newsSources,
		logger:     logger,
	}
}

type SourceType int

const (
	RapidAPI SourceType = iota
	RegularNewLinks
	CustomLinks
)

func (st SourceType) String() string {
	return [...]string{"Rapid", "Regular", "Custom"}[st]
}

func ToSourceType(st string) SourceType {
	switch st {
	case "Rapid":
		return RapidAPI
	case "Regular":
		return RegularNewLinks
	case "Custom":
		return CustomLinks
	default:
		fmt.Println("Unknown Source Type")
		return CustomLinks
	}
}

type Source struct {
	SourceType SourceType `json:"news_source_type"`
	Url        string     `json:"news_source_url"`
	Key        string     `json:"news_source_key"`
	Host       string     `json:"news_source_host"`
}

func NewNewsSource(sourceType string, url string, key string, host string) Source {
	return Source{SourceType: ToSourceType(sourceType), Url: url, Key: key, Host: host}
}

func CreateCryptoNewsSources() Source {
	return NewNewsSource("Rapid",
		"https://crypto-news-live3.p.rapidapi.com/news",
		"79f6e24995msh6c1d3f6a1bdb266p1aa6fbjsn423a6dc42ba8",
		"crypto-news-live3.p.rapidapi.com")
}

type Item struct {
	Title  string `json:"title"`
	Url    string `json:"url"`
	Source string `json:"source"`
}

type Result struct {
	NewsItems []Item
}

func NewResult() Result {
	return Result{}
}

func (nc newsClient) FetchData() Result {
	result := NewResult()

	for _, source := range nc.newsSource {
		switch source.SourceType {
		case RapidAPI:
			rapidAPIResults, err := nc.rapidAPIFetch(source)
			if err == nil {
				result.NewsItems = append(result.NewsItems, rapidAPIResults...)
			}
		case RegularNewLinks:
		case CustomLinks:
		default:
		}
	}
	return result
}

func (nc newsClient) rapidAPIFetch(source Source) ([]Item, error) {
	req, _ := http.NewRequest("GET", source.Url, nil)

	req.Header.Add("X-RapidAPI-Key", source.Key)
	req.Header.Add("X-RapidAPI-Host", source.Host)

	var result []Item

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		nc.logger.Error("fail to fetch from rapid API ", err)
		return []Item{}, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		nc.logger.Error("fail to read from rapid API result body ", err)
		return []Item{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		nc.logger.Error("fail to unmarshal rapid API body ", err)
		return []Item{}, err
	}
	return result, nil
}

func (nc newsClient) parseRegularFetch(source Source) (Result, error) {
	//TODO implement me
	panic("implement me")
}

func (nc newsClient) parseCustomLinkFetch(source Source) (Result, error) {
	//TODO implement me
	panic("implement me")
}

package service

import "goognion/src/repository"

type Service struct {
	Crawler ICrawlerService
	Logger  ILogger
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		Crawler: NewCrawlerService(repository),
		Logger:  NewLogger(repository)}
}

type ICrawlerService interface {
	Crawl(url string, depth int) <-chan error
	GetNotCrawledUrls() ([]string, error)
	GetNotCrawledUrlsWithLimit(limit int) ([]string, error)
}

type ILogger interface {
	Log(<-chan error)
	LogOnce(error)
}

package service

import "goognion/src/repository"

type Service struct {
	Crawler ICrawlerService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		Crawler: NewCrawlerService(repository)}
}

type ICrawlerService interface {
	Crawl(url string, depth int)
	GetUncrawledUrls(maxDepth int) []string
}

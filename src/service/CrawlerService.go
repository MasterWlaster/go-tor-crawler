package service

import (
	"goognion/src"
	"goognion/src/repository"
)

type CrawlerService struct {
	repository repository.Repository
}

func NewCrawlerService(repository *repository.Repository) *CrawlerService {
	return &CrawlerService{repository: *repository}
}

func (s *CrawlerService) Crawl(url string, depth int) {
	var err error = nil
	defer s.repository.Logger.Log(err)

	isValid, err := s.repository.UrlValidator.IsValid(url)
	if err != nil || !isValid {
		return
	}

	used, err := s.repository.Memory.UsedUrl(url)
	if err != nil || used {
		return
	}

	err = s.repository.Memory.Remember(url)
	if err != nil {
		return
	}

	if depth--; depth < 0 {
		return
	}

	data, urls, err := s.repository.Crawler.Load(url)
	if err != nil {
		return
	}

	for _, u := range urls {
		if u == url {
			continue
		}
		go s.Crawl(u, depth)
	}

	indexes, err := s.repository.Crawler.DoIndexing(data)
	if err != nil {
		return
	}

	err = s.repository.Memory.Save(src.Page{Url: url, Indexes: indexes})
	if err != nil {
		return
	}
}

func (s *CrawlerService) GetNotCrawledUrls() []string {
	var err error = nil
	defer s.repository.Logger.Log(err)

	urls, err := s.repository.Memory.GetUnusedUrls()
	if err != nil {
		return nil
	}

	return urls
}

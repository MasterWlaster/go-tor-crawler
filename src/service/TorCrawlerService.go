package service

import (
	"goognion/src"
	"goognion/src/repository"
)

type TorCrawlerService struct {
}

func NewTorCrawlerService(repository ICrawlerRepository) *TorCrawlerService {
	return &TorCrawlerService{}
}

func (s *TorCrawlerService) Crawl(url string, depth int) {
	if depth--; depth < 0 {
		return
	}

	indexes := map[string]int{}
	data, urls := repository.Crawler.Load(url)

	for _, u := range urls {
		go s.Crawl(u, depth)
	}

	repository.Crawler.DoIndexing(data, indexes)
	repository.Crawler.Save(src.Page{Url: url, Indexes: indexes})
}

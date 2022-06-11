package service

import (
	"goognion/src"
	"goognion/src/repository"
)

type TorCrawlerService struct {
}

func (s *TorCrawlerService) Crawl(url string, depth int) {
	var err error = nil
	defer repository.Logger.Log(err)

	if depth--; depth < 0 {
		return
	}

	data, urls, err := repository.Crawler.Load(url)
	if err != nil {
		return
	}

	for _, u := range urls {
		go s.Crawl(u, depth)
	}
	indexes := map[string]int{}

	err = repository.Crawler.DoIndexing(data, indexes)
	if err != nil {
		return
	}

	err = repository.Crawler.Save(src.Page{Url: url, Indexes: indexes})
	if err != nil {
		return
	}
}

func (s *TorCrawlerService) GetUncrawledUrls(maxDepth int) []string {
	//todo
	return nil
}

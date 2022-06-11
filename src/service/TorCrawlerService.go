package service

import (
	"goognion/src"
	"goognion/src/repository"
)

type CrawlerService struct {
}

func (s *CrawlerService) Crawl(url string, depth int) {
	var err error = nil
	defer repository.Logger.Log(err)

	isValid, err := repository.UrlValidator.IsValid(url)
	if err != nil || !isValid {
		return
	}

	used, err := repository.Crawler.UsedUrl(url, depth)
	if err != nil || used {
		return
	}

	err = repository.Crawler.Remember(url)
	if err != nil {
		return
	}

	if depth--; depth < 0 {
		return
	}

	data, urls, err := repository.Crawler.Load(url)
	if err != nil {
		return
	}

	for _, u := range urls {
		if u == url {
			continue
		}
		go s.Crawl(u, depth)
	}

	indexes, err := repository.Crawler.DoIndexing(data)
	if err != nil {
		return
	}

	err = repository.Crawler.Save(src.Page{Url: url, Indexes: indexes})
	if err != nil {
		return
	}
}

func (s *CrawlerService) GetUncrawledUrls(maxDepth int) []string {
	//todo
	return nil
}

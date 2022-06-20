package service

import (
	"fmt"
	"goognion/src"
	"goognion/src/repository"
	"sync"
)

type CrawlerService struct {
	repository repository.Repository
}

func NewCrawlerService(repository *repository.Repository) *CrawlerService {
	return &CrawlerService{repository: *repository}
}

func (s *CrawlerService) Crawl(url string, depth int) {
	isValid, err := s.repository.UrlValidator.IsValid(url)
	if err != nil {
		s.repository.Logger.Log(err)
		return
	}
	if !isValid {
		s.repository.Logger.Log(fmt.Errorf("invalid url: %s", url))
		return
	}

	used, err := s.repository.Memory.UsedUrl(url)
	if err != nil {
		s.repository.Logger.Log(err)
		return
	}
	if used {
		s.repository.Logger.Log(fmt.Errorf("already used url: %s", url))
		return
	}

	err = s.repository.Memory.Remember(url)
	if err != nil {
		s.repository.Logger.Log(err)
		return
	}

	if depth--; depth < 0 {
		return
	}

	data, urls, err := s.repository.Crawler.Load(url)
	if err != nil {
		s.repository.Logger.Log(err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for u := range urls {
			wg.Add(1)
			if u == url {
				wg.Done()
				continue
			}
			u, err = s.repository.UrlValidator.Relative(u, url)
			if err != nil {
				wg.Done()
				continue
			}
			go func(u string) {
				s.Crawl(u, depth)
				wg.Done()
			}(u)
		}
		wg.Done()
	}()

	defer wg.Wait()

	indexes, err := s.repository.Crawler.DoIndexing(data)
	if err != nil {
		s.repository.Logger.Log(err)
		return
	}

	err = s.repository.Memory.Save(src.Page{Url: url, Indexes: indexes})
	if err != nil {
		s.repository.Logger.Log(err)
		return
	}
}

func (s *CrawlerService) GetNotCrawledUrls() []string {
	urls, err := s.repository.Memory.GetUnusedUrls()
	if err != nil {
		s.repository.Logger.Log(err)
		return nil
	}

	return urls
}

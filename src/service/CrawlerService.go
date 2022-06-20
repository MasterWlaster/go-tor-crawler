package service

import (
	"fmt"
	"goognion/src"
	"goognion/src/repository"
	"sync"
)

type CrawlerService struct {
	repository *repository.Repository
}

func NewCrawlerService(repository *repository.Repository) *CrawlerService {
	return &CrawlerService{repository: repository}
}

func (s *CrawlerService) Crawl(url string, depth int) <-chan error {
	errs := make(chan error)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go s.crawl(url, depth, errs, wg)

	go func() {
		wg.Wait()
		close(errs)
	}()

	return errs
}

func (s *CrawlerService) GetNotCrawledUrls() ([]string, error) {
	urls, err := s.repository.Memory.GetUnusedUrls()
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func (s *CrawlerService) GetNotCrawledUrlsWithLimit(limit int) ([]string, error) {
	urls, err := s.repository.Memory.GetUnusedUrlsWithLimit(limit)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func (s *CrawlerService) crawl(url string, depth int, errOut chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	isValid, err := s.repository.UrlValidator.IsValid(url)
	if err != nil {
		errOut <- err
		return
	}
	if !isValid {
		errOut <- fmt.Errorf("invalid url: %s", url)
		return
	}

	used, err := s.repository.Memory.UsedUrl(url)
	if err != nil {
		errOut <- err
		return
	}
	if used {
		errOut <- fmt.Errorf("already used url: %s", url)
		return
	}

	err = s.repository.Memory.Remember(url)
	if err != nil {
		errOut <- err
		return
	}

	if depth--; depth < 0 {
		return
	}

	data, urls, err := s.repository.Crawler.Load(url)
	if err != nil {
		errOut <- err
		return
	}

	wg.Add(1)
	go func() {
		for u := range urls {
			if u == url {
				continue
			}

			u, err = s.repository.UrlValidator.Relative(u, url)
			if err != nil {
				continue
			}

			wg.Add(1)
			go s.crawl(u, depth, errOut, wg)
		}
		wg.Done()
	}()

	indexes, err := s.repository.Crawler.DoIndexing(data)
	if err != nil {
		errOut <- err
		return
	}

	err = s.repository.Memory.Save(src.Page{Url: url, Indexes: indexes})
	if err != nil {
		errOut <- err
		return
	}
}

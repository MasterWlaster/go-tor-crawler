package service

import (
	"errors"
	"goognion/src"
	"goognion/src/repository"
)

type CrawlerService struct {
	repository repository.Repository
}

func NewCrawlerService(repository *repository.Repository) *CrawlerService {
	return &CrawlerService{repository: *repository}
}

func (s *CrawlerService) Crawl(url string, depth int) { //todo logging
	isValid, err := s.repository.UrlValidator.IsValid(url)
	if err != nil {
		s.repository.Logger.Log(err)
		return
	}
	if !isValid {
		s.repository.Logger.Log(errors.New("url is invalid"))
		return
	}

	used, err := s.repository.Memory.UsedUrl(url)
	if err != nil {
		s.repository.Logger.Log(err)
		return
	}
	if used {
		s.repository.Logger.Log(errors.New("url is used"))
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

	//wg := sync.WaitGroup{}
	//wg.Add(len(urls))
	go func() {
		for u := range urls {
			if u == url {
				//wg.Done()
				continue
			}
			//todo relative path
			go func(u string) {
				s.Crawl(u, depth)
				//wg.Done()
			}(u)
		}
	}()

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

	//wg.Wait()
}

func (s *CrawlerService) GetNotCrawledUrls() []string {
	urls, err := s.repository.Memory.GetUnusedUrls()
	if err != nil {
		s.repository.Logger.Log(err)
		return nil
	}

	return urls
}

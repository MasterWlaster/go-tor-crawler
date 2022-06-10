package service

import "goognion/src"

type ICrawlerService interface {
	Crawl(url string, depth int)
	GetUncrawledUrls(maxDepth int) []string
}

type ICrawlerRepository interface {
	Save(page src.Page) error
	Load(url string) ([]src.Text, []string, error)
	DoIndexing(src []src.Text, dst map[string]int) error
}

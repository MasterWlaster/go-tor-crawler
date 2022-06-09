package service

import "goognion/src"

type ICrawlerService interface {
	Crawl(url string, depth int)
}

type ICrawlerRepository interface {
	Save(page src.Page)
	Load(url string) ([]src.Text, []string)
	DoIndexing(src []src.Text, dst map[string]int)
}

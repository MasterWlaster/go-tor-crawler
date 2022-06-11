package service

type ICrawlerService interface {
	Crawl(url string, depth int)
	GetUncrawledUrls(maxDepth int) []string
}

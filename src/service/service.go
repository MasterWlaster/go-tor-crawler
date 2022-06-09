package service

var Crawler ICrawlerService

func Init(crawler ICrawlerService) {
	Crawler = crawler
}

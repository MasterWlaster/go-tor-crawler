package service

var Crawler ICrawlerService
var Logger ILoggerService

func Init(crawler ICrawlerService, logger ILoggerService) {
	Crawler = crawler
	Logger = logger
}

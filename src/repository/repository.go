package repository

var Crawler ICrawlerRepository
var Logger ILoggerRepository

func Init(crawler ICrawlerRepository, logger ILoggerRepository) {
	Crawler = crawler
	Logger = logger
}

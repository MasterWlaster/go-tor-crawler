package repository

var Crawler ICrawlerRepository
var Logger ILoggerRepository
var UrlValidator IUrlValidator

func Init(crawler ICrawlerRepository, logger ILoggerRepository, urlValidator IUrlValidator) {
	Crawler = crawler
	Logger = logger
	UrlValidator = urlValidator
}

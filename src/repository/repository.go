package repository

import "goognion/src/service"

var Crawler service.ICrawlerRepository
var Logger service.ILoggerRepository

func InitRepositories(crawler service.ICrawlerRepository) {
	Crawler = crawler
}

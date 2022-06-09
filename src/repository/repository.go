package repository

import "goognion/src/service"

var Crawler service.ICrawlerRepository

func InitRepositories(crawler service.ICrawlerRepository) {
	Crawler = crawler
}

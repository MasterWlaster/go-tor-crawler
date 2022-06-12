package repository

import "goognion/src"

type Repository struct {
	Crawler      ICrawlerRepository
	Logger       ILoggerRepository
	UrlValidator IUrlValidator
}

func NewRepository(crawler ICrawlerRepository, logger ILoggerRepository, urlValidator IUrlValidator) *Repository {
	return &Repository{Crawler: crawler, Logger: logger, UrlValidator: urlValidator}
}

func NewTorRepository() *Repository {
	return &Repository{
		Crawler:      NewTorCrawlerRepository(),
		Logger:       NewConsoleLogger(),
		UrlValidator: NewTorUrlValidator()}
}

type IUrlValidator interface {
	IsValid(url string) (bool, error)
}

type ICrawlerRepository interface {
	Save(page src.Page) error
	Load(url string) ([]src.Text, []string, error)
	DoIndexing(src []src.Text) (map[string]int, error)
	UsedUrl(url string, depth int) (bool, error)
	Remember(url string) error
}

type ILoggerRepository interface {
	Log(err error) // error
}

package repository

import "goognion/src"

type Repository struct {
	Crawler      ICrawlerRepository
	Logger       ILoggerRepository
	UrlValidator IUrlValidator
	Memory       IMemoryRepository
}

func NewTorRepository() *Repository {
	return &Repository{
		Crawler:      NewTorCrawlerRepository(),
		Logger:       NewConsoleLogger(),
		UrlValidator: NewTorUrlValidator(),
		Memory:       NewPostgresDb()}
}

type IUrlValidator interface {
	IsValid(url string) (bool, error)
}

type ICrawlerRepository interface {
	Load(url string) ([]src.Text, []string, error)
	DoIndexing(src []src.Text) (map[string]int, error)
}

type IMemoryRepository interface {
	Save(page src.Page) error
	Remember(url string) error
	UsedUrl(url string, depth int) (bool, error)
}

type ILoggerRepository interface {
	Log(err error)
}

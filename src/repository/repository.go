package repository

import (
	"github.com/cretz/bine/tor"
	"github.com/jmoiron/sqlx"
	"goognion/src"
	"net/http"
)

type Repository struct {
	Crawler      ICrawlerRepository
	Logger       ILoggerRepository
	UrlValidator IUrlValidator
	Memory       IMemoryRepository
}

func NewTorRepository(db *sqlx.DB, client *http.Client, tor *tor.Tor) *Repository {
	return &Repository{
		Crawler:      NewTorCrawlerRepository(client, tor),
		Logger:       NewConsoleLogger(),
		UrlValidator: NewTorUrlValidator(),
		Memory:       NewPostgresDb(db)}
}

type IUrlValidator interface {
	IsValid(url string) (bool, error)
	Relative(relative string, url string) (string, error)
}

type ICrawlerRepository interface {
	Load(url string) (<-chan src.Text, <-chan string, error)
	DoIndexing(input <-chan src.Text) (map[string]int, error)
}

type IMemoryRepository interface {
	Save(page src.Page) error
	Remember(url string) error
	UsedUrl(url string) (bool, error)
	GetUnusedUrls() ([]string, error)
	GetUnusedUrlsWithLimit(limit int) ([]string, error)
}

type ILoggerRepository interface {
	Log(err error)
}

type DbConfig struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
	SslMode  string
}

type TorConfig struct {
	ExePath string
	DataDir string
}

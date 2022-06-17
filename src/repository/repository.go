package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"goognion/src"
)

type Repository struct {
	Crawler      ICrawlerRepository
	Logger       ILoggerRepository
	UrlValidator IUrlValidator
	Memory       IMemoryRepository
}

func NewTorRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Crawler:      NewTorCrawlerRepository(),
		Logger:       NewConsoleLogger(),
		UrlValidator: NewTorUrlValidator(),
		Memory:       NewPostgresDb(db)}
}

func ConnectPostgres(host, port, username, password, dbName, sslMode string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, username, dbName, password, sslMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

type IUrlValidator interface {
	IsValid(url string) (bool, error)
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
}

type ILoggerRepository interface {
	Log(err error)
}

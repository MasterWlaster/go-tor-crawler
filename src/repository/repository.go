package repository

import (
	"fmt"
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

func ConnectPostgres(host, port, username, password, dbName, sslMode string) (*sqlx.DB, error) {
	fmt.Println("Подключение к БД...")

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, username, dbName, password, sslMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Успешно подключено!")

	return db, nil
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
}

type ILoggerRepository interface {
	Log(err error)
}

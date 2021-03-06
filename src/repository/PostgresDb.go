package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"goognion/src"
	"strings"
)

const (
	urls    = `urls`
	words   = `words`
	indexes = `indexes`
)

type PostgresDb struct {
	db *sqlx.DB
}

func NewPostgresDb(db *sqlx.DB) *PostgresDb {
	return &PostgresDb{db: db}
}

func (p *PostgresDb) Save(page src.Page) error {
	if len(page.Indexes) == 0 {
		return fmt.Errorf("empty indexes: %s", page.Url)
	}

	wordsB := strings.Builder{}
	indexesB := strings.Builder{}

	iw := 1
	ii := 2
	ws := make([]interface{}, len(page.Indexes)+1)
	ws[0] = page.Url
	for w, c := range page.Indexes {
		ws[iw] = w
		wordsB.WriteString(fmt.Sprintf(` ($%d),`, iw))
		indexesB.WriteString(fmt.Sprintf(` ($%d, $1, %d),`, ii, c))
		iw++
		ii++
	}

	wb := wordsB.String()
	wb = wb[:len(wb)-1]

	i := indexesB.String()
	i = i[:len(i)-1]

	wordsQ := fmt.Sprintf(`INSERT INTO %s VALUES %s ON CONFLICT DO NOTHING`, words, wb)
	_, err := p.db.Exec(wordsQ, ws[1:]...)
	if err != nil {
		return err
	}

	indexesQ := fmt.Sprintf(`INSERT INTO %s VALUES %s ON CONFLICT DO NOTHING`, indexes, i)
	_, err = p.db.Exec(indexesQ, ws...)

	return err
}

func (p *PostgresDb) Remember(url string) error {
	query := fmt.Sprintf(`INSERT INTO %s VALUES ($1) ON CONFLICT DO NOTHING`, urls)
	_, err := p.db.Exec(query, url)

	return err
}

func (p *PostgresDb) UsedUrl(url string) (bool, error) {
	query := fmt.Sprintf(`SELECT COUNT(1) FROM %s WHERE url = $1`, indexes)
	row := p.db.QueryRow(query, url)
	c := 0
	err := row.Scan(&c)

	return c != 0, err
}

func (p *PostgresDb) GetUnusedUrls() ([]string, error) {
	query := fmt.Sprintf(
		`SELECT value FROM %s WHERE value NOT IN (SELECT DISTINCT url FROM indexes)`, urls)
	var urls []string
	err := p.db.Select(&urls, query)

	return urls, err
}

func (p *PostgresDb) GetUnusedUrlsWithLimit(limit int) ([]string, error) {
	query := fmt.Sprintf(
		`SELECT value FROM %s WHERE value NOT IN (SELECT DISTINCT url FROM indexes) LIMIT %d`,
		urls, limit)
	var urls []string
	err := p.db.Select(&urls, query)

	return urls, err
}

func ConnectPostgres(c DbConfig) (*sqlx.DB, error) {
	fmt.Println("?????????????????????? ?? ????...")

	db, err := sqlx.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Name, c.Password, c.SslMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("?????????????? ????????????????????!")

	return db, nil
}

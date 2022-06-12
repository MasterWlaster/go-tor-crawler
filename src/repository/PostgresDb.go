package repository

import "goognion/src"

type PostgresDb struct {
}

func NewPostgresDb() *PostgresDb {
	return &PostgresDb{}
}

func (p PostgresDb) Save(page src.Page) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDb) Remember(url string) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresDb) UsedUrl(url string, depth int) (bool, error) {
	//TODO implement me
	panic("implement me")
}

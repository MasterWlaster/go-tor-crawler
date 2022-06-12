package repository

import (
	"goognion/src"
)

type TorCrawlerRepository struct {
}

func NewTorCrawlerRepository() *TorCrawlerRepository {
	return &TorCrawlerRepository{}
}

func (t *TorCrawlerRepository) Load(url string) ([]src.Text, []string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TorCrawlerRepository) DoIndexing(src []src.Text) (map[string]int, error) {
	//TODO implement me
	panic("implement me")
}

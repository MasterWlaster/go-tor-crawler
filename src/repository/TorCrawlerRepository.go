package repository

import (
	"goognion/src"
)

type TorCrawlerRepository struct {
}

func (t *TorCrawlerRepository) Remember(url string) error {
	//TODO implement me
	panic("implement me")
}

func (t *TorCrawlerRepository) Save(page src.Page) error {
	//TODO implement me
	panic("implement me")
}

func (t *TorCrawlerRepository) Load(url string) ([]src.Text, []string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TorCrawlerRepository) DoIndexing(src []src.Text) (map[string]int, error) {
	//TODO implement me
	panic("implement me")
}

func (t *TorCrawlerRepository) UsedUrl(url string, depth int) (bool, error) {
	//TODO implement me
	panic("implement me")
}

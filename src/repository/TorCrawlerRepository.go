package repository

import (
	"goognion/src"
)

type TorCrawlerRepository struct {
}

func (t TorCrawlerRepository) Save(page src.Page) error {
	//TODO implement me
	panic("implement me")
}

func (t TorCrawlerRepository) Load(url string) ([]src.Text, []string, error) {
	//TODO implement me
	panic("implement me")
}

func (t TorCrawlerRepository) DoIndexing(src []src.Text, dst map[string]int) error {
	//TODO implement me
	panic("implement me")
}



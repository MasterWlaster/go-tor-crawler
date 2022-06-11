package repository

import "goognion/src"

type ICrawlerRepository interface {
	Save(page src.Page) error
	Load(url string) ([]src.Text, []string, error)
	DoIndexing(src []src.Text) (map[string]int, error)
	UsedUrl(url string, depth int) (bool, error)
	Remember(url string) error
}

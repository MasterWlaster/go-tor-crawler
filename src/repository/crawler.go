package repository

import "goognion/src"

type ICrawlerRepository interface {
	Save(page src.Page) error
	Load(url string) ([]src.Text, []string, error)
	DoIndexing(src []src.Text, dst map[string]int) error
}

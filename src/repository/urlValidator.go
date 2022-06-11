package repository

type IUrlValidator interface {
	IsValid(url string) (bool, error)
}

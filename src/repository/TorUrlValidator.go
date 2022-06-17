package repository

import "regexp"

type TorUrlValidator struct {
}

func NewTorUrlValidator() *TorUrlValidator {
	return &TorUrlValidator{}
}

func (t TorUrlValidator) IsValid(url string) (bool, error) {
	return regexp.MatchString(`http(s?)://(.+)\.onion((/(.+))*)((/(.*))?)`, url)
}

package repository

type TorUrlValidator struct {
}

func NewTorUrlValidator() *TorUrlValidator {
	return &TorUrlValidator{}
}

func (t TorUrlValidator) IsValid(url string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

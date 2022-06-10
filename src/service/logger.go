package service

type ILoggerService interface {
	Log(err error)
}

type ILoggerRepository interface {
	Log(err error) error
}

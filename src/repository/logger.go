package repository

type ILoggerRepository interface {
	Log(err error) // error
}

package service

import "goognion/src/repository"

type Logger struct {
	repository *repository.Repository
}

func NewLogger(repository *repository.Repository) *Logger {
	return &Logger{repository: repository}
}

func (l *Logger) Log(errors <-chan error) {
	for err := range errors {
		l.LogOnce(err)
	}
}

func (l *Logger) LogOnce(err error) {
	if err == nil {
		return
	}
	l.repository.Logger.Log(err)
}

package service

import "goognion/src/repository"

type ConsoleLogger struct {
}

func (l *ConsoleLogger) Log(err error) {
	if err == nil {
		return
	}

	err = repository.Logger.Log(err)
	if err != nil {
		//todo
	}
}


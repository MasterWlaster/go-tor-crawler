package repository

import "fmt"

type ConsoleLogger struct {
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (c *ConsoleLogger) Log(err error) {
	if err == nil {
		return
	}
	fmt.Println(fmt.Sprintf("LOG ERROR: %s", err))
}

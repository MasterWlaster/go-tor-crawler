package repository

import "fmt"

type ConsoleLogger struct {
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (c *ConsoleLogger) Log(err error) {
	fmt.Println(fmt.Sprintf("- ERROR: %s", err))
}

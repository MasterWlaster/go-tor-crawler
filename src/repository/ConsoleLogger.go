package repository

import "fmt"

type ConsoleLogger struct {
}

func (c *ConsoleLogger) Log(err error) {
	if err == nil {
		return
	}
	fmt.Printf("LOG ERROR: %s", err)
}

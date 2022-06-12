package main

import (
	"goognion/src/controller"
	"goognion/src/repository"
	"goognion/src/service"
)

func main() {
	r := repository.NewTorRepository()
	s := service.NewService(r)
	c := controller.NewConsoleController(s)

	c.Run()

	//todo shutdown
}

package main

import (
	"goognion/src/controller"
	"goognion/src/repository"
	"goognion/src/service"
)

func main() {
	repository.Init(&repository.TorCrawlerRepository{}, &repository.ConsoleLogger{})
	service.Init(&service.TorCrawlerService{})
	controller.Init(&controller.WindowsConsoleController{})

	controller.Controller.Run()
}

package main

import (
	"goognion/src/controller"
	"goognion/src/repository"
	"goognion/src/service"
)

func main() {
	repository.Init(
		&repository.TorCrawlerRepository{},
		&repository.ConsoleLogger{},
		&repository.TorUrlValidator{})
	service.Init(
		&service.CrawlerService{})
	controller.Init(
		&controller.WindowsConsoleController{})

	controller.Controller.Run()

	//todo shutdown
}

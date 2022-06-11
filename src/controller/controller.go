package controller

type IController interface {
	Run()
}

var Controller IController

func Init(controller IController) {
	Controller = controller
}

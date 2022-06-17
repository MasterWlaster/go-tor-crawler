package main

import (
	_ "github.com/lib/pq"
	"goognion/src/controller"
	"goognion/src/repository"
	"goognion/src/service"
)

func main() {
	db, err := repository.ConnectPostgres(
		"localhost",
		"5432",
		"postgres",
		"postgres",
		"go_crawler",
		"disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := repository.NewTorRepository(db)
	s := service.NewService(r)
	c := controller.NewConsoleController(s)

	c.Run()
}

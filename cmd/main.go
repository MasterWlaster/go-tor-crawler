package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"goognion/src/controller"
	"goognion/src/repository"
	"goognion/src/service"
	"os"
	"os/signal"
	"syscall"
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
	db.SetMaxOpenConns(32)
	defer func() {
		err := db.Close()
		fmt.Println(err)
	}()

	client, tor, err := repository.NewTorClient(
		"D:/Tor Browser/Browser/TorBrowser/Tor/tor.exe",
		"")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := tor.Close()
		fmt.Println(err)
	}()

	r := repository.NewTorRepository(db, client, tor)
	s := service.NewService(r)
	c := controller.NewConsoleController(s)

	go c.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
}

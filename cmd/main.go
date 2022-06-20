package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"goognion/src/controller"
	"goognion/src/repository"
	"goognion/src/service"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	fmt.Println("\nСтарт...")

	err := initConfig()
	if err != nil {
		panic(err)
	}

	db, err := repository.ConnectPostgres(repository.DbConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Name:     viper.GetString("db.name"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		SslMode:  viper.GetString("db.ssl_mode"),
	})
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(2 * runtime.NumCPU())

	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println(fmt.Sprintf("Не удалось закрыть подключение к бд: %s", err))
		}
	}()

	client, tor, err := repository.NewTorClient(repository.TorConfig{
		ExePath: viper.GetString("tor.exe_path"),
		DataDir: viper.GetString("tor.data_dir"),
	})
	if err != nil {
		panic(err)
	}

	defer func() {
		err := tor.Close()
		if err != nil {
			fmt.Println(fmt.Sprintf("Не удалось закрыть Tor: %s", err))
		}
	}()

	r := repository.NewTorRepository(db, client, tor)
	s := service.NewService(r)
	c := controller.NewConsoleController(s)

	go c.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	fmt.Println("Завершение работы...")
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("app-config")
	return viper.ReadInConfig()
}

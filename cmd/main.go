package main

import (
	"fmt"
	"log"

	smarthelpdesc "github.com/arynskiii/help_desk"
	"github.com/arynskiii/help_desk/internal/handler"
	"github.com/arynskiii/help_desk/internal/repository"
	"github.com/arynskiii/help_desk/internal/service"

	"github.com/spf13/viper"
)

func main() {

	if err := InitConfig(); err != nil {
		log.Fatal("failed to init config: ", err)
	}
	db, err := repository.NewMySqlDB(&repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		fmt.Print(err)
		log.Fatal("failed to initialize db: ", err.Error())
	}
	if err := repository.CreateTable(db); err != nil {
		log.Fatal("failed to create table : ", err)
	}

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)
	server := new(smarthelpdesc.Server)
	if err := server.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
		log.Fatal("failed to run server: ", err)
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

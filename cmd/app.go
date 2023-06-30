package main

import (
	"github.com/cost_control"
	"github.com/cost_control/config"
	"github.com/cost_control/internal/handlers/api"
	"github.com/cost_control/internal/handlers/telegram"
	"github.com/cost_control/pkg/repository"
	"log"
)

const botToken = "5836425300:AAG2azf8sY54f_Y9Mod1PM9vY6IzIWRpnq0"

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := repository.NewMongoDB(&repository.Config{
		Host:       cfg.Db.Mongo.Host,
		Port:       cfg.Db.Mongo.Port,
		Database:   cfg.Db.Mongo.Database,
		Collection: cfg.Db.Mongo.Collection,
	})
	if err != nil {
		log.Fatal(err)
	}

	apiHandlers := api.New(db)
	telegramBot, err := telegram.New(botToken, db)
	if err != nil {
		log.Fatal(err)
	}
	err = telegramBot.Start(nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	srv := new(cost_control.Server)
	if err := srv.Run(10000, apiHandlers.InitRoutes()); err != nil {
		log.Fatal(err)
	}
}

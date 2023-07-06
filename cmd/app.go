package main

import (
	"github.com/cost_control"
	"github.com/cost_control/config"
	"github.com/cost_control/internal/handlers/api"
	"github.com/cost_control/internal/handlers/rpc"
	"github.com/cost_control/internal/handlers/telegram"
	"github.com/cost_control/pkg/repository"
	"log"
	"sync"
)

const (
	botToken          = "5836425300:AAG2azf8sY54f_Y9Mod1PM9vY6IzIWRpnq0"
	RestApiServerPort = 10000
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := repository.NewMongoDB(&repository.Config{
		Host:     cfg.Db.Mongo.Host,
		Port:     cfg.Db.Mongo.Port,
		Database: cfg.Db.Mongo.Database,
	})
	if err != nil {
		log.Fatal(err)
	}

	apiHandlers := api.New(db, cfg)
	telegramBot, err := telegram.New(botToken, db)
	if err != nil {
		log.Fatal(err)
	}
	rpcHandler := rpc.New(db)
	var wg sync.WaitGroup
	wg.Add(3)
	srv := new(cost_control.Server)
	go func(wg *sync.WaitGroup, restApiServerPort int) {
		log.Print("Start Rest Api server")
		if err := srv.Run(restApiServerPort, apiHandlers.InitRoutes()); err != nil {
			log.Fatal(err)
			wg.Done()
		}
	}(&wg, RestApiServerPort)

	go func(wg *sync.WaitGroup) {
		err = telegramBot.Start(wg, nil, nil)
		if err != nil {
			log.Fatal(err)
			wg.Done()
		}
	}(&wg)

	err = rpcHandler.Start()
	if err != nil {
		log.Fatal(err)
		wg.Done()
	}

	wg.Wait()
}

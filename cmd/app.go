package main

import (
	"github.com/cost_control"
	"github.com/cost_control/config"
	"github.com/cost_control/internal/handlers/api"
	"github.com/cost_control/internal/handlers/rpc"
	"github.com/cost_control/internal/handlers/telegram"
	"github.com/cost_control/pkg/logger"
	"github.com/cost_control/pkg/repository"
	"sync"
)

func main() {
	log := logger.New()
	log.Info("инициализация конфигурации")
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
	telegramBot, err := telegram.New(cfg.TelegramBotToken, db, log)
	rpcHandler := rpc.New(db, cfg, log)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(3)
	srv := new(cost_control.Server)

	go func(wg *sync.WaitGroup, port int, log logger.ILogger) {
		log.Info("Start Rest Api server")
		if err := srv.Run(port, apiHandlers.InitRoutes()); err != nil {
			log.Fatal(err)
			wg.Done()
		}
	}(&wg, cfg.Rest.Port, log)

	go func(wg *sync.WaitGroup, log logger.ILogger) {
		err = telegramBot.Start(wg, nil, nil)
		if err != nil {
			log.Fatal(err)
			wg.Done()
		}
	}(&wg, log)

	err = rpcHandler.Start()
	if err != nil {
		log.Fatal(err)
		wg.Done()
	}

	wg.Wait()
}

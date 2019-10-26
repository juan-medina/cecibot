package main

import (
	"github.com/juan-medina/cecibot/bot"
	"github.com/juan-medina/cecibot/config"
	"go.uber.org/zap"
)

func main() {

	log,_ := zap.NewProduction()
	defer log.Sync()

	log.Info("Reading config.")
	var cfg, err = config.FromProvider(config.EnvironmentVariables())

	if err != nil {
		log.Error("Error reading config", zap.Error(err))
		return
	}

	log.Info("Creating bot.")

	bt, err := bot.New(cfg)
	if err != nil {
		log.Error("Error creating bot", zap.Error(err))
		return
	}

	log.Info("Starting bot.")

	err = bt.Run()

	if err != nil {
		log.Error("Error running bot,", zap.Error(err))
		return
	}

	log.Info("Bot stopped.")

}

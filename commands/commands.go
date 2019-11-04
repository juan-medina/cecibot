package commands

import (
	"github.com/juan-medina/cecibot/commands/basic"
	"github.com/juan-medina/cecibot/commands/raid"
	"github.com/juan-medina/cecibot/commands/system"
	"github.com/juan-medina/cecibot/prototype"
	"go.uber.org/zap"
)

func New(processor prototype.Processor) []prototype.Provider {
	log, _ := zap.NewProduction()
	defer log.Sync()

	log.Info("creating command providers.")
	var providers = []prototype.Provider{
		basic.New(processor),
		system.New(processor),
		raid.New(processor),
	}

	log.Info("Commands providers created.", zap.Int("number of providers", len(providers)))
	return providers
}

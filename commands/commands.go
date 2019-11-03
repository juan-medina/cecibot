package commands

import (
	"github.com/juan-medina/cecibot/commands/basic"
	"github.com/juan-medina/cecibot/commands/system"
	"github.com/juan-medina/cecibot/prototype"
)

func New(processor prototype.Processor) []prototype.Provider {
	var providers = []prototype.Provider{
		basic.New(processor),
		system.New(processor),
	}
	return providers
}

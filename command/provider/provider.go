package provider

import "github.com/juan-medina/cecibot/prototype"

type BaseProvider struct {
	commands prototype.CommandsMap
	prc      prototype.Processor
}

func (b *BaseProvider) GetCommands() *prototype.CommandsMap {
	return &b.commands
}

func (b *BaseProvider) AddCommand(cmd *prototype.Command) {
	b.commands[cmd.Key] = cmd
}
func (b *BaseProvider) GetProcessor() prototype.Processor {
	return b.prc
}

func New(prc prototype.Processor) *BaseProvider {
	var prov = BaseProvider{
		commands: make(prototype.CommandsMap),
		prc:      prc,
	}
	return &prov
}

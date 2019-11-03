package prototype

import "github.com/juan-medina/cecibot/config"

type Bot interface {
	Run() error
	GetConfig() config.Config
}

type Processor interface {
	ProcessMessage(text string, author string) string
	Init(bot Bot) error
	End()
	IsOwner(userId string) bool
	GetCommandHelp(key string) string
	GetHelp() string
}

type CommandFunction func(args []string, author string) string

type Command struct {
	Key  string
	Desc string
	Fun  CommandFunction
	Help string
}

type CommandsMap map[string]*Command

type Provider interface {
	GetCommands() *CommandsMap
	AddCommand(cmd *Command)
	GetProcessor() Processor
}

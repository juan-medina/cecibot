package basic

import (
	"github.com/juan-medina/cecibot/command"
	"github.com/juan-medina/cecibot/command/provider"
	"github.com/juan-medina/cecibot/prototype"
)

type basicCommands struct {
	*provider.BaseProvider
}

func (d basicCommands) ping(args []string, author string) string {
	return "pong!"
}

func (d basicCommands) hello(args []string, author string) string {
	if d.GetProcessor().IsOwner(author) {
		return "hello master!"
	}
	return "hello!"
}

func New(p prototype.Processor) prototype.Provider {
	var prov = basicCommands{BaseProvider: provider.New(p)}

	prov.AddCommand(command.New("ping",
		"Asks for a ping to the *bot*.",
		"This is a test command for the *bot* that will reply with a pong message",
		prov.ping),
	)
	prov.AddCommand(command.New("hello",
		"Greets the *user*.",
		"This command will greet *you* back.",
		prov.hello),
	)

	return prov
}

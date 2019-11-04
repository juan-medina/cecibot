package system

import (
	"fmt"
	"github.com/juan-medina/cecibot/command"
	"github.com/juan-medina/cecibot/command/provider"
	"github.com/juan-medina/cecibot/prototype"
	"go.uber.org/zap"
)

type systemCommands struct {
	*provider.BaseProvider
}

func (d systemCommands) help(args []string, author string) string {
	argc := len(args)
	if argc > 0 {
		key := args[0]
		help := d.GetProcessor().GetCommandHelp(key)
		if help != "" {
			return fmt.Sprintf("Command **%s** : \n%s", key, help)
		}

		return "Unknown command in help. " + d.GetProcessor().GetHelp()
	}
	return d.GetProcessor().GetHelp()
}

func New(p prototype.Processor) prototype.Provider {
	log, _ := zap.NewProduction()
	defer log.Sync()

	log.Info("Creating system commands")
	var prov = systemCommands{BaseProvider: provider.New(p)}

	prov.AddCommand(command.New("help",
		"Gets help with *commands*.",
		"Usage:\n\t**help** *command*\n\nUse this command to get help with any *command*.",
		prov.help),
	)

	log.Info("System commands created", zap.Int("number of commands", len(*prov.GetCommands())))
	return prov
}

package processor

import "fmt"

func SystemCommands() CommandProvider {
	provider := CommandProviderImpl{}
	provider.Init()

	provider.Add(helpCommand)

	return provider
}

func helpCommand(p *processorImpl) *command {
	return &command{
		key:  "help",
		desc: "Gets help with *commands*.",
		help: "Usage:\n\t**help** *command*\n\nUse this command to get help with any *command*.",
		fun: func(args []string, author string) string {
			argc := len(args)
			if argc > 0 {
				key := args[0]
				cmd, found := p.commands[key]
				if found {
					return fmt.Sprintf("Command **%s** : \n%s", key, cmd.help)
				}

				return "Unknown command in help. " + p.help
			}
			return p.help
		},
	}
}

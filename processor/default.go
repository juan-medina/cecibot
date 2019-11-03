package processor

func DefaultCommands() CommandProvider {
	provider := CommandProviderImpl{}
	provider.Init()

	provider.Add(helloCommand)
	provider.Add(pingCommand)

	return provider
}

func helloCommand(p *processorImpl) *command {
	return &command{
		key:  "hello",
		desc: "Greets the *user*.",
		help: "This command will greet *you* back.",
		fun: func(args []string, author string) string {
			if p.isOwner(author) {
				return "hello master!"
			}
			return "hello!"
		},
	}
}

func pingCommand(_ *processorImpl) *command {
	return &command{
		key:  "ping",
		desc: "Asks for a ping to the *bot*.",
		help: "This is a test command for the *bot* that will reply with a pong message",
		fun: func(args []string, author string) string {
			return "pong!"
		},
	}
}

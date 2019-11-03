package processor

type CommandProviderImpl struct {
	commands []createCommandFun
}

func (d CommandProviderImpl) getCommands() []createCommandFun {
	return d.commands
}

func (d CommandProviderImpl) Init() {
}

func (d *CommandProviderImpl) Add(fun createCommandFun) {
	d.commands = append(d.commands, fun)
}

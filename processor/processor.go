package processor

import (
	"fmt"
	"github.com/juan-medina/cecibot/commands"
	"github.com/juan-medina/cecibot/prototype"
	"go.uber.org/zap"
	"strings"
	"unicode"
)

type processorImpl struct {
	bot      prototype.Bot
	owner    string
	commands prototype.CommandsMap
	help     string
}

func (p *processorImpl) AddCommand(cmd *prototype.Command) {
	p.commands[cmd.Key] = cmd
}

func (p *processorImpl) generateHelp() {
	help := "Available commands are:"
	for key, cmd := range p.commands {
		help += fmt.Sprintf("\n\t **%s** : %q", key, cmd.Desc)
	}
	help += "\n\nTo get help on any *command* send:\n\t**help** *command*"
	p.help = help

}

func (p *processorImpl) addCommands(provider prototype.Provider) {
	for _, cmd := range *provider.GetCommands() {
		p.AddCommand(cmd)
	}
}

func (p *processorImpl) configure() {
	p.owner = p.bot.GetConfig().GetOwner()
}

func (p *processorImpl) Init(bot prototype.Bot) error {
	log, _ := zap.NewProduction()
	defer log.Sync()

	log.Info("Processor initialising.")

	p.bot = bot

	log.Info("Configuring processor.")
	p.configure()

	log.Info("Adding commands.")
	for _, prov := range commands.New(p) {
		p.addCommands(prov)
	}
	log.Info("Commands added.", zap.Int("number of commands", len(p.commands)))

	log.Info("Generating commands help.")
	p.generateHelp()

	log.Info("Processor initialised.")
	return nil
}

func (p processorImpl) End() {
	log, _ := zap.NewProduction()
	defer log.Sync()

	log.Info("Processor end.")
}

func New() *prototype.Processor {
	var prc prototype.Processor = &processorImpl{commands: make(prototype.CommandsMap)}
	return &prc
}

func (p processorImpl) IsOwner(author string) bool {
	return p.owner == author
}

func (p processorImpl) GetCommandHelp(key string) string {
	cmd, found := p.commands[key]
	if found {
		return cmd.Help
	}
	return ""
}

func (p processorImpl) GetHelp() string {
	return p.help
}

func (p processorImpl) parseCommand(text string) (key string, args []string) {
	var m []string = nil
	var s string

	str := strings.TrimSpace(text) + " "

	var lastQuote int32 = 0
	isSpace := false
	for i, c := range str {
		switch {
		// If we're ending a quote, break out and skip this character
		case c == lastQuote:
			lastQuote = 0

		// If we're in a quote, count this character
		case lastQuote != 0:
			s += string(c)

		// If we encounter a quote, enter it and skip this character
		case unicode.In(c, unicode.Quotation_Mark):
			isSpace = false
			lastQuote = c

		// If it's a space, store the string
		case unicode.IsSpace(c):
			if 0 == i || isSpace {
				continue
			}
			isSpace = true
			m = append(m, s)
			s = ""

		default:
			isSpace = false
			s += string(c)
		}

	}

	if m == nil || lastQuote != 0 {
		return "", []string{}
	}

	return m[0], m[1:]
}

func (p processorImpl) ProcessMessage(text string, author string) string {

	key, args := p.parseCommand(text)
	cmd, found := p.commands[key]
	if found {
		return cmd.Fun(args, author)
	}

	return "Unknown command. " + p.help
}

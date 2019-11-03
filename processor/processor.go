package processor

import (
	"fmt"
	"github.com/juan-medina/cecibot/prototype"
	"strings"
	"unicode"
)

type Processor interface {
	ProcessMessage(text string, author string) string
	Init(bot prototype.Bot) error
	End()
}
type commandFunction func(args []string, author string) string
type command struct {
	key  string
	desc string
	fun  commandFunction
	help string
}

type processorImpl struct {
	bot      prototype.Bot
	owner    string
	commands map[string]command
	help     string
}

type createCommandFun func(p *processorImpl) *command
type CommandProvider interface {
	getCommands() []createCommandFun
}

func (p *processorImpl) AddCommand(fun createCommandFun) {
	cmd := fun(p)
	p.commands[cmd.key] = *cmd
}

func (p *processorImpl) generateHelp() {
	help := "Available commands are:"
	for key, cmd := range p.commands {
		help += fmt.Sprintf("\n\t **%s** : %q", key, cmd.desc)
	}
	help += "\n\nTo get help on any *command* send:\n\t**help** *command*"
	p.help = help

}

func (p *processorImpl) addCommands(provider CommandProvider) {
	for _, cmdFun := range provider.getCommands() {
		p.AddCommand(cmdFun)
	}

	p.generateHelp()
}

func (p *processorImpl) configure() {
	p.owner = p.bot.GetConfig().GetOwner()
}

func (p *processorImpl) Init(bot prototype.Bot) error {
	p.bot = bot

	p.configure()
	p.addCommands(SystemCommands())
	p.addCommands(DefaultCommands())

	return nil
}

func (p processorImpl) End() {
}

func New() *Processor {
	var prc Processor = &processorImpl{commands: make(map[string]command)}
	return &prc
}

func (p processorImpl) isOwner(author string) bool {
	return p.owner == author
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
		return cmd.fun(args, author)
	}

	return "Unknown command. " + p.help
}

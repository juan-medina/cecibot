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

type defaultProcessor struct {
	bot      prototype.Bot
	owner    string
	commands map[string]command
	help     string
}

func (p *defaultProcessor) AddCommand(key string, fun commandFunction, desc string, help string) {
	p.commands[key] = command{
		key:  key,
		desc: desc,
		fun:  fun,
		help: help,
	}
}

func (p *defaultProcessor) Init(bot prototype.Bot) error {
	p.bot = bot
	p.owner = bot.GetConfig().GetOwner()

	p.AddCommand("ping", func(args []string, author string) string {
		return "pong!"
	}, "Asks for a ping to the *bot*.", "This is a test command for the *bot* that will reply with a pong message")

	p.AddCommand("hello", func(args []string, author string) string {
		if p.isOwner(author) {
			return "hello master!"
		}
		return "hello!"
	}, "Greets the *user*.", "This command will greet *you* back.")

	p.AddCommand("help", func(args []string, author string) string {
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

	}, "Gets help with *commands*.", "Usage:\n\t**help** *command*\n\nUse this command to get help with any *command*.")

	help := "Available commands are:"
	for key, cmd := range p.commands {
		help += fmt.Sprintf("\n\t **%s** : %q", key, cmd.desc)
	}
	help += "\n\nTo get help on any *command* send:\n\t**help** *command*"
	p.help = help

	return nil
}

func (p defaultProcessor) End() {
}

func New() *Processor {
	var prc Processor = &defaultProcessor{commands: make(map[string]command)}
	return &prc
}

func (p defaultProcessor) isOwner(author string) bool {
	return p.owner == author
}

func (p defaultProcessor) parseCommand(text string) (key string, args []string) {
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

func (p defaultProcessor) ProcessMessage(text string, author string) string {

	key, args := p.parseCommand(text)
	cmd, found := p.commands[key]
	if found {
		return cmd.fun(args, author)
	}

	return "Unknown command. " + p.help
}

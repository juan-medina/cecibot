package processor

import (
	"github.com/juan-medina/cecibot/prototype"
)

type Processor interface {
	ProcessMessage(text string, author string) string
	Init(bot prototype.Bot) error
	End()
}

type defaultProcessor struct {
	bot   prototype.Bot
	owner string
}

func (p *defaultProcessor) Init(bot prototype.Bot) error {
	p.bot = bot
	p.owner = bot.GetConfig().GetOwner()
	return nil
}

func (p defaultProcessor) End() {
}

func New() *Processor {
	var prc Processor = &defaultProcessor{}
	return &prc
}

func (p defaultProcessor) isOwner(author string) bool {
	return p.owner == author
}
func (p defaultProcessor) ProcessMessage(text string, author string) string {
	if text == "ping" {
		return "pong!"
	} else if text == "pong" {
		return "ping!"
	} else if text == "hello" {
		if p.isOwner(author) {
			return "hello master!"
		}
		return "hello!"
	}

	return ""
}

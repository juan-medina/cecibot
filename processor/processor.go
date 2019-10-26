package processor

import "github.com/juan-medina/cecibot/config"

type Processor interface {
	ProcessMessage(text string, author string) string
}

type defaultProcessor struct {
	cfg config.Config
}

func New(cfg config.Config) Processor {
	return defaultProcessor{cfg: cfg}
}

func (p defaultProcessor) isOwner(author string) bool {
	return p.cfg.GetOwner() == author
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

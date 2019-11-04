package data

import (
	"github.com/juan-medina/cecibot/commands/raid/data/memory"
	"github.com/juan-medina/cecibot/prototype"
)

func New() prototype.RaidDataProvider {
	return memory.New()
}

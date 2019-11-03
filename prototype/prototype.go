package prototype

import "github.com/juan-medina/cecibot/config"

type Bot interface {
	Run() error
	GetConfig() config.Config
}

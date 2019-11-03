package command

import "github.com/juan-medina/cecibot/prototype"

func New(key string, desc string, help string, fun prototype.CommandFunction) *prototype.Command {
	return &prototype.Command{
		Key:  key,
		Desc: desc,
		Fun:  fun,
		Help: help,
	}
}

package raid

import (
	"fmt"
	"github.com/juan-medina/cecibot/command"
	"github.com/juan-medina/cecibot/command/provider"
	"github.com/juan-medina/cecibot/commands/raid/data"
	"github.com/juan-medina/cecibot/prototype"
	"go.uber.org/zap"
)

type raidCommands struct {
	*provider.BaseProvider
	data prototype.RaidDataProvider
}

func (d *raidCommands) officers() string {
	result := "raid officers:\n"
	for _, officer := range d.data.GetOfficers() {
		result += fmt.Sprintf("\t<@%s>\n", officer.Id)
	}
	return result
}

func (d *raidCommands) deleteOfficer(args []string) string {
	argc := len(args)
	if argc > 0 {
		id := args[0]
		d.data.DeleteOfficer(id)
		return fmt.Sprintf("officer <@%s> deleted", id)
	}

	return ""
}

func (d *raidCommands) addOfficer(args []string) string {
	argc := len(args)
	if argc > 0 {
		id := args[0]
		d.data.AddOfficer(id)
		return fmt.Sprintf("officer <@%s> added", id)
	}

	return ""
}

func (d *raidCommands) officer(args []string) string {
	argc := len(args)
	if argc > 0 {
		sub := args[0]
		if sub == "delete" {
			return d.deleteOfficer(args[1:])
		} else if sub == "add" {
			return d.addOfficer(args[1:])
		}
	}
	return ""
}

func (d *raidCommands) raid(args []string, author string) string {
	argc := len(args)
	if argc > 0 {
		sub := args[0]
		if sub == "officers" {
			return d.officers()
		} else if sub == "officer" {
			if d.isOfficer(author) {
				return d.officer(args[1:])
			} else {
				return "this command is for **officers** only"
			}
		}
	}

	return ""
}

func (d *raidCommands) isOfficer(id string) bool {
	if d.GetProcessor().IsOwner(id) {
		return true
	}

	for _, officer := range d.data.GetOfficers() {
		if officer.Id == id {
			return true
		}
	}

	return false
}

func New(p prototype.Processor) prototype.Provider {
	log, _ := zap.NewProduction()
	defer log.Sync()

	log.Info("Creating raid commands")
	var prov = raidCommands{
		BaseProvider: provider.New(p),
		data:         data.New(),
	}

	prov.AddCommand(command.New("raid",
		"Manage *raid* attendance.",
		`With this command you could create, list and confirm raid attendance
Usage:
	**raid** *option* *parameters*
*Options* for *members* and *officer* are:
	**list**
		list next raids, and their *raid-id*
	**sign up** *raid-id* *char* *class* *spec*
		confirm/change attendance for the desired *raid-id* wth the *char* using the given *class* and *spec*
	**sign down** *raid-id*
		sign down for attendance for the desired *raid-id*
	**rooster** *raid-id*
		shows the rooster for the given *raid-id*
	**officers**
		list raid officers
*Options* for *officers* only are:
    **create** *name* *date*
		creates a raid with the given *name* and *date*. Shows the *raid-id*
	**cancels** *raid-id*
		cancel the raid indicated by the *raid-id*
	**officer add** *discord-id*
		add a raid officer with it *discord-id*
	**officer delete** *discord-id*
		delete a raid officer with it *discord-id*
`,
		prov.raid),
	)

	log.Info("Raid commands created", zap.Int("number of commands", len(*prov.GetCommands())))
	return prov
}

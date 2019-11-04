package memory

import (
	"github.com/juan-medina/cecibot/commands/raid/data/entities"
	"github.com/juan-medina/cecibot/prototype"
	"sort"
)

type inMemory struct {
	officers map[string]entities.Officer
}

func (d *inMemory) AddOfficer(id string) {
	officer := entities.Officer{Id: id}
	d.officers[id] = officer
}

func (d *inMemory) DeleteOfficer(id string) {
	delete(d.officers, id)
}

func (d *inMemory) GetOfficers() []entities.Officer {
	result := make([]entities.Officer, 0)

	keys := make([]string, 0)
	for key := range d.officers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		result = append(result, d.officers[key])
	}
	return result
}

func New() prototype.RaidDataProvider {
	return &inMemory{
		officers: make(map[string]entities.Officer),
	}
}

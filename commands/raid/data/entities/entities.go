package entities

import "time"

type Officer struct {
	Id string
}

type Raid struct {
	Id   string
	name string
	date time.Time
}

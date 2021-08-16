package models

import "time"

type Event struct {
	Id        uint
	Title     string
	CreatedBy string
	StartDate time.Time
}

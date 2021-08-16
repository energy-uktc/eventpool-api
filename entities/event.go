package entities

import (
	"time"

	"gorm.io/gorm"
)

//Event ...
type Event struct {
	gorm.Model
	Title     string
	CreatedBy string
	StartDate time.Time
}

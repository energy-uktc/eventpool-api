package models

import "time"

type Activity struct {
	Id          string     `json:"id"`
	EventID     string     `json:"eventId"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DateTime    *time.Time `json:"dateTime"`
	Location    *Location  `json:"location,omitempty"`
}

type CreateUpdateActivity struct {
	EventID     string
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description" binding:"required"`
	DateTime    *time.Time `json:"dateTime" binding:"required"`
	Location    *Location  `json:"location"`
}

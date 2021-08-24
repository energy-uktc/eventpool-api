package models

import "time"

type Event struct {
	Id                 string       `json:"id"`
	Title              string       `json:"title"`
	Description        string       `json:"description"`
	CreatedBy          *UserModel   `json:"createdBy,omitempty"`
	StartDate          *time.Time   `json:"startDate"`
	EndDate            *time.Time   `json:"endDate,omitempty"`
	NumberOfAtendees   int          `json:"numberOfAtendees"`
	Atendees           []*UserModel `json:"atendees,omitempty"`
	NumberOfActivities int          `json:"numberOfActivities"`
	Location           *Location    `json:"location,omitempty"`
}

type CreateEvent struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description" binding:"required"`
	CreatedBy   string     `json:"createdBy"`
	StartDate   *time.Time `json:"startDate" binding:"required"`
	EndDate     *time.Time `json:"endDate"`
}

type UpdateEvent struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Location    *Location  `json:"location"`
}

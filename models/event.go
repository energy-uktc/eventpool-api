package models

import (
	"time"
)

type Event struct {
	Id                 string             `json:"id"`
	Title              string             `json:"title"`
	Description        string             `json:"description"`
	CreatedBy          *UserModel         `json:"createdBy,omitempty"`
	StartDate          *time.Time         `json:"startDate"`
	EndDate            *time.Time         `json:"endDate,omitempty"`
	NumberOfAttendees  int                `json:"numberOfAttendees"`
	Attendees          []*SimpleUserModel `json:"attendees,omitempty"`
	NumberOfActivities int                `json:"numberOfActivities"`
	Activities         []*Activity        `json:"activities,omitempty"`
	NumberOfPolls      int                `json:"numberOfPolls"`
	Polls              []*PollModel       `json:"polls,omitempty"`
	Location           *Location          `json:"location,omitempty"`
}

type CreateEvent struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description" binding:"required"`
	CreatedBy   string     `json:"createdBy"`
	StartDate   *time.Time `json:"startDate" binding:"required,gte"`
	EndDate     *time.Time `json:"endDate"`
}

type UpdateEvent struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Location    *Location  `json:"location"`
}

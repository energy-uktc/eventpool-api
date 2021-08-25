package entities

import (
	"encoding/json"
	"time"

	"github.com/energy-uktc/eventpool-api/models"
	"github.com/gofrs/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

//Event ...
type Event struct {
	ID          string `gorm:"primarykey;type:uuid"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Description string
	CreatedByID string `gorm:"type:uuid"`
	CreatedBy   *User  `gorm:"foreignKey:CreatedByID"`
	StartDate   *time.Time
	EndDate     *time.Time
	Location    datatypes.JSON
	Atendees    []User     `gorm:"many2many:user_events;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Activities  []Activity `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Polls       []Poll     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (e *Event) ToModel() *models.Event {
	model := &models.Event{
		Id:                 e.ID,
		Title:              e.Title,
		Description:        e.Description,
		StartDate:          e.StartDate,
		EndDate:            e.EndDate,
		NumberOfAtendees:   len(e.Atendees),
		NumberOfActivities: len(e.Activities),
		NumberOfPolls:      len(e.Polls),
	}

	if e.CreatedBy != nil {
		model.CreatedBy = e.CreatedBy.ToModel()
	}
	if e.Atendees != nil {
		for _, user := range e.Atendees {
			model.Atendees = append(model.Atendees, user.ToSimpleModel())
		}
	}
	if e.Activities != nil {
		for _, activity := range e.Activities {
			model.Activities = append(model.Activities, activity.ToModel())
		}
	}
	if e.Polls != nil {
		for _, poll := range e.Polls {
			model.Polls = append(model.Polls, poll.ToModel())
		}
	}

	if e.Location != nil {
		locationBytes, _ := e.Location.MarshalJSON()
		location := &models.Location{}
		json.Unmarshal(locationBytes, location)
		model.Location = location
	}
	return model
}

func (e *Event) ParseUpdateModel(eventModel *models.UpdateEvent) {
	e.Title = eventModel.Title
	e.Description = eventModel.Description
	e.StartDate = eventModel.StartDate
	e.EndDate = eventModel.EndDate
	e.Location = nil
	if eventModel.Location != nil && eventModel.Location.Latitude != 0 && eventModel.Location.Longitude != 0 {
		locationBytes, _ := json.Marshal(eventModel.Location)
		e.Location = datatypes.JSON(locationBytes)
	}

}

func (e *Event) ParseCreateModel(eventModel *models.CreateEvent) {
	e.Title = eventModel.Title
	e.Description = eventModel.Description
	e.CreatedByID = eventModel.CreatedBy
	e.StartDate = eventModel.StartDate
	e.EndDate = eventModel.EndDate
	e.Atendees = []User{{ID: eventModel.CreatedBy}}
}

func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	uid, _ := uuid.NewV4()
	e.ID = uid.String()
	return
}

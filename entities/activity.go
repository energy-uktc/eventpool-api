package entities

import (
	"encoding/json"
	"time"

	"github.com/energy-uktc/eventpool-api/models"
	"github.com/gofrs/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

//Activity ...
type Activity struct {
	ID          string `gorm:"primarykey;type:uuid"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Description string
	DateTime    *time.Time
	Location    datatypes.JSON
	EventID     string
}

func (a *Activity) BeforeCreate(tx *gorm.DB) (err error) {
	uid, _ := uuid.NewV4()
	a.ID = uid.String()
	return
}

func (a *Activity) ParseModel(model *models.CreateUpdateActivity) {
	a.EventID = model.EventID
	a.Title = model.Title
	a.Description = model.Description
	a.DateTime = model.DateTime
	a.Location = nil
	if model.Location != nil && model.Location.Latitude != 0 && model.Location.Longitude != 0 {
		locationBytes, _ := json.Marshal(model.Location)
		a.Location = datatypes.JSON(locationBytes)
	}
}

func (a *Activity) ToModel() *models.Activity {
	model := &models.Activity{
		Id:          a.ID,
		Title:       a.Title,
		Description: a.Description,
		DateTime:    a.DateTime,
		EventID:     a.EventID,
	}
	if a.Location != nil {
		locationBytes, _ := a.Location.MarshalJSON()
		location := &models.Location{}
		json.Unmarshal(locationBytes, location)
		model.Location = location
	}
	return model
}

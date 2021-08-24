package entities

import (
	"time"

	"github.com/energy-uktc/eventpool-api/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

//User ...
type User struct {
	ID        string `gorm:"primarykey;type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserName  string
	Email     string `gorm:"uniqueIndex"`
	Password  string
	Verified  bool
	Events    []Event `gorm:"many2many:user_events;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	uid, _ := uuid.NewV4()
	u.ID = uid.String()
	return
}

func NewUserFromID(id string) *User {
	return &User{ID: id}
}

func (u *User) ToModel() *models.UserModel {
	return &models.UserModel{
		Id:       u.ID,
		Email:    u.Email,
		UserName: u.UserName,
	}
}

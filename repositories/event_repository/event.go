package event_repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/energy-uktc/eventpool-api/database"
	"github.com/energy-uktc/eventpool-api/entities"
	"github.com/energy-uktc/eventpool-api/repositories/poll_repository"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Create(event *entities.Event) (string, error) {
	response := database.DbConn.Omit("Attendees.*").Create(event)
	if response.Error != nil {
		log.Println(response.Error)
		return "", fmt.Errorf("Something went wrong during event creation.")
	}
	return event.ID, nil
}

func Update(event *entities.Event) (*entities.Event, error) {
	response := database.DbConn.Omit(clause.Associations).Save(event)
	if response.Error != nil {
		log.Println(response.Error)
		return nil, fmt.Errorf("Something went wrong during event update.")
	}
	return FindById(event.ID)
}

func UpdatePartial(event *entities.Event) (*entities.Event, error) {
	response := database.DbConn.Updates(&event)
	if response.Error != nil {
		log.Println(response.Error)
		return nil, fmt.Errorf("Something went wrong during event update.")
	}
	return FindById(event.ID)
}

func FindById(id string) (*entities.Event, error) {
	var event *entities.Event
	var users []entities.User
	var activities []entities.Activity
	var polls []*entities.Poll
	var err error

	response := database.DbConn.Joins("CreatedBy").First(&event, uuid.FromStringOrNil(id))
	database.DbConn.Model(&event).Association("Attendees").Find(&users)
	database.DbConn.Model(&event).Association("Activities").Find(&activities)
	if polls, err = poll_repository.FindForEvent(id); err != nil {
		log.Println(response.Error)
		return nil, fmt.Errorf("Something went wrong")
	}

	event.Attendees = users
	event.Activities = activities
	for _, poll := range polls {
		event.Polls = append(event.Polls, *poll)
	}

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Event Not Found")
		}
		log.Println(response.Error)
		return nil, fmt.Errorf("Something went wrong")
	}

	return event, nil
}

func FindUserEvents(user *entities.User) ([]entities.Event, error) {
	var events []entities.Event
	err := database.DbConn.Model(&user).Association("Events").Find(&events)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Something went wrong")
	}
	return events, nil
}

func AppendAtendee(event *entities.Event, user *entities.User) error {
	err := database.DbConn.Omit("Attendees.*").Model(&event).Association("Attendees").Append(&entities.User{ID: user.ID})
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Something went wrong")
	}
	return nil
}

func RemoveAtendee(event *entities.Event, user *entities.User) error {
	err := database.DbConn.Omit("Attendees.*").Model(event).Association("Attendees").Delete(user)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Something went wrong")
	}
	return nil
}

func CountAttendees(event *entities.Event) int {
	return int(database.DbConn.Model(&event).Association("Attendees").Count())
}

func CountActivities(event *entities.Event) int {
	return int(database.DbConn.Model(&event).Association("Activities").Count())
}

func CountPolls(event *entities.Event) int {
	return int(database.DbConn.Model(&event).Association("Polls").Count())
}

func Delete(id string) error {
	uid, err := uuid.FromString(id)
	if err != nil {
		return fmt.Errorf("Event Not Found")
	}
	response := database.DbConn.Delete(&entities.Event{}, uid)
	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Event Not Found")
		}
		log.Println(response.Error)
		return fmt.Errorf("Something went wrong")
	}
	return nil
}

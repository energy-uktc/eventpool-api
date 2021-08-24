package event_repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/energy-uktc/eventpool-api/database"
	"github.com/energy-uktc/eventpool-api/entities"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

func Create(event *entities.Event) (string, error) {
	response := database.DbConn.Omit("Atendees.*").Create(event)
	if response.Error != nil {
		log.Println(response.Error)
		return "", fmt.Errorf("Something went wrong during event creation.")
	}
	return event.ID.String(), nil
}

func Update(event *entities.Event) (*entities.Event, error) {
	response := database.DbConn.Save(event)
	if response.Error != nil {
		log.Println(response.Error)
		return nil, fmt.Errorf("Something went wrong during event update.")
	}
	return FindById(event.ID.String())
}

func FindById(id string) (*entities.Event, error) {
	var event *entities.Event
	var users []entities.User
	response := database.DbConn.Joins("CreatedBy").First(&event, uuid.FromStringOrNil(id))
	database.DbConn.Model(&event).Association("Atendees").Find(&users)
	event.Atendees = users
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

func CountActivities(event *entities.Event) int {
	// var events []entities.Event
	// err := database.DbConn.Model(&event).Association("A").Find(&events)
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, fmt.Errorf("Something went wrong")
	// }
	// return events, nil
	return 0
}

func CountAtendees(event *entities.Event) int {
	return int(database.DbConn.Model(&event).Association("Atendees").Count())
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

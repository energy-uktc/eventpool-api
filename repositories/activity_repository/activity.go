package activity_repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/energy-uktc/eventpool-api/database"
	"github.com/energy-uktc/eventpool-api/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Create(activity *entities.Activity) (string, error) {
	response := database.DbConn.Create(activity)
	if response.Error != nil {
		log.Println(response.Error)
		return "", fmt.Errorf("Something went wrong during activity creation.")
	}
	return activity.ID, nil
}

func Update(activity *entities.Activity) (*entities.Activity, error) {
	response := database.DbConn.Omit(clause.Associations).Save(activity)
	if response.Error != nil {
		log.Println(response.Error)
		return nil, fmt.Errorf("Something went wrong during activity update.")
	}
	return FindById(activity.EventID, activity.ID)
}

func FindById(eventId string, id string) (*entities.Activity, error) {
	var activity *entities.Activity
	response := database.DbConn.Where("event_id = ? AND id = ?", eventId, id).First(&activity)
	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Activity Not Found")
		}
		log.Println(response.Error)
		return nil, fmt.Errorf("Something went wrong")
	}

	return activity, nil
}

func FindForEvent(eventId string) ([]*entities.Activity, error) {
	var activities []*entities.Activity
	response := database.DbConn.Where("event_id = ?", eventId).Find(&activities)
	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Activity Not Found")
		}
		log.Println(response.Error)
		return nil, fmt.Errorf("Something went wrong")
	}

	return activities, nil
}

func Delete(eventId string, id string) error {
	response := database.DbConn.Delete(&entities.Activity{ID: id, EventID: eventId})
	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Activity Not Found")
		}
		log.Println(response.Error)
		return fmt.Errorf("Something went wrong")
	}
	return nil
}

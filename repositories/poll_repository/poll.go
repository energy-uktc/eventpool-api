package poll_repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/energy-uktc/eventpool-api/database"
	"github.com/energy-uktc/eventpool-api/entities"
	"gorm.io/gorm"
)

func Create(poll *entities.Poll) (string, error) {
	response := database.DbConn.Debug().Create(poll)
	if response.Error != nil {
		log.Println(response.Error)
		return "", fmt.Errorf("Something went wrong during poll creation.")
	}
	return poll.ID, nil
}

func FindForEvent(eventId string) ([]*entities.Poll, error) {
	var polls []*entities.Poll
	response := database.DbConn.Joins("CreatedBy").Where("event_id = ?", eventId).Find(&polls)
	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Polls Not Found")
		}
		log.Println(response.Error)
		return nil, fmt.Errorf("Something went wrong")
	}

	for _, poll := range polls {
		var answers []entities.PollAnswer
		if err := database.DbConn.Preload("Votes").Where(&entities.PollAnswer{PollID: poll.ID}).Find(&answers).Error; err != nil {
			log.Println(err)
			return nil, fmt.Errorf("Something went wrong")
		}
		poll.Options = answers
	}
	return polls, nil
}

func FindById(eventId string, id string) (*entities.Poll, error) {
	var poll *entities.Poll
	response := database.DbConn.Joins("CreatedBy").Where(&entities.Poll{ID: id, EventID: eventId}).First(&poll)
	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Poll Not Found")
		}
		log.Println(response.Error)
		return nil, fmt.Errorf("Something went wrong")
	}
	var answers []entities.PollAnswer
	if err := database.DbConn.Preload("Votes").Where(&entities.PollAnswer{PollID: poll.ID}).Find(&answers).Error; err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Something went wrong")
	}
	poll.Options = answers
	return poll, nil
}

func Delete(eventId string, id string) error {
	response := database.DbConn.Delete(&entities.Poll{ID: id, EventID: eventId})
	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Poll Not Found")
		}
		log.Println(response.Error)
		return fmt.Errorf("Something went wrong")
	}
	return nil
}

func AppendVote(option *entities.PollAnswer, user *entities.User) error {
	err := database.DbConn.Omit("Votes.*").Model(&option).Association("Votes").Append(user)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Something went wrong")
	}
	return nil
}

func RemoveVote(option *entities.PollAnswer, user *entities.User) error {
	err := database.DbConn.Omit("Votes.*").Model(&option).Association("Votes").Delete(user)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Something went wrong")
	}
	return nil
}

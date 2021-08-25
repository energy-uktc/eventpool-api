package poll_service

import (
	"fmt"
	"time"

	"github.com/energy-uktc/eventpool-api/entities"
	"github.com/energy-uktc/eventpool-api/models"
	"github.com/energy-uktc/eventpool-api/repositories/poll_repository"
	"github.com/energy-uktc/eventpool-api/services/user_service"
)

func Create(model *models.CreatePollModel) (*models.PollModel, error) {
	poll := new(entities.Poll)
	poll.ParseCreateModel(model)
	id, err := poll_repository.Create(poll)
	if err != nil {
		return nil, err
	}
	poll, err = poll_repository.FindById(model.EventID, id)
	if err != nil {
		return nil, err
	}
	return poll.ToModel(), nil
}

func FindAll(eventId string) ([]*models.PollModel, error) {
	polls, err := poll_repository.FindForEvent(eventId)
	if err != nil {
		return nil, err
	}
	var pollModels []*models.PollModel
	for _, poll := range polls {
		model := poll.ToModel()
		pollModels = append(pollModels, model)
	}
	return pollModels, nil
}

func FindById(eventId string, id string) (*models.PollModel, error) {
	poll, err := poll_repository.FindById(eventId, id)
	if err != nil {
		return nil, err
	}
	return poll.ToModel(), nil
}

func Delete(eventId string, id string) error {
	return poll_repository.Delete(eventId, id)
}

func Vote(positive bool, userId string, eventId string, pollId string, optionId string) error {
	var vote func(*entities.PollAnswer, *entities.User) error
	if positive {
		vote = poll_repository.AppendVote
	} else {
		vote = poll_repository.RemoveVote
	}

	if _, err := user_service.GetUser(userId); err != nil {
		return err
	}
	user := &entities.User{ID: userId}
	poll, err := poll_repository.FindById(eventId, pollId)
	if err != nil {
		return err
	}
	if poll.EndTime.Before(time.Now()) {
		return fmt.Errorf("The poll is closed for voting")
	}

	for _, option := range poll.Options {
		if option.ID == optionId {
			if err := vote(&option, user); err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("Voting Option Not Found")
}

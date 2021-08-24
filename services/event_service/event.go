package event_service

import (
	"time"

	"github.com/energy-uktc/eventpool-api/entities"
	"github.com/energy-uktc/eventpool-api/models"
	"github.com/energy-uktc/eventpool-api/repositories/event_repository"
	"github.com/energy-uktc/eventpool-api/repositories/user_repository"
)

func Create(model *models.CreateEvent) (*models.Event, error) {
	event := new(entities.Event)
	event.ParseCreateModel(model)
	id, err := event_repository.Create(event)
	if err != nil {
		return nil, err
	}
	event, err = event_repository.FindById(id)
	if err != nil {
		return nil, err
	}
	return event.ToModel(), nil
}

func Update(id string, model *models.UpdateEvent) (*models.Event, error) {
	event, err := event_repository.FindById(id)
	if err != nil {
		return nil, err
	}
	event.ParseUpdateModel(model)
	updatedEvent, err := event_repository.Update(event)
	if err != nil {
		return nil, err
	}

	return updatedEvent.ToModel(), nil
}

func FindAllForUser(userId string) ([]*models.Event, error) {
	user, err := user_repository.FindById(userId)
	if err != nil {
		return nil, err
	}
	events, err := event_repository.FindUserEvents(user)
	if err != nil {
		return nil, err
	}
	var eventModels []*models.Event
	for _, event := range events {
		eventModel := event.ToModel()
		eventModel.NumberOfAtendees = event_repository.CountAtendees(&event)
		eventModel.NumberOfActivities = event_repository.CountActivities(&event)
		eventModels = append(eventModels, eventModel)
	}
	return eventModels, nil
}

func FindById(id string) (*models.Event, error) {
	event, err := event_repository.FindById(id)
	if err != nil {
		return nil, err
	}
	return event.ToModel(), nil
}

func Delete(id string) error {
	return event_repository.Delete(id)
}

func FindActiveEvents(userId string) ([]*models.Event, error) {
	events, err := FindAllForUser(userId)
	if err != nil {
		return nil, err
	}
	var activeEvents []*models.Event
	date := time.Now()
	for _, event := range events {
		if (event.EndDate != nil && event.EndDate.After(date)) || (event.StartDate.After(date)) {
			activeEvents = append(activeEvents, event)
		}
	}
	return activeEvents, nil
}
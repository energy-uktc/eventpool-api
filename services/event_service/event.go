package event_service

import (
	"time"

	"github.com/energy-uktc/eventpool-api/entities"
	"github.com/energy-uktc/eventpool-api/models"
	"github.com/energy-uktc/eventpool-api/repositories/event_repository"
	"github.com/energy-uktc/eventpool-api/repositories/user_repository"
	"github.com/energy-uktc/eventpool-api/services/user_service"
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

func Update(id string, partial bool, model *models.UpdateEvent) (*models.Event, error) {
	event, err := event_repository.FindById(id)
	if err != nil {
		return nil, err
	}

	var updatedEvent *entities.Event
	if partial {
		eventToUpdate := &entities.Event{
			ID: event.ID,
		}
		eventToUpdate.ParseUpdateModel(model)
		updatedEvent, err = event_repository.UpdatePartial(eventToUpdate)
	} else {
		event.ParseUpdateModel(model)
		updatedEvent, err = event_repository.Update(event)
	}

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
		eventModel.NumberOfAttendees = event_repository.CountAttendees(&event)
		eventModel.NumberOfActivities = event_repository.CountActivities(&event)
		eventModel.NumberOfPolls = event_repository.CountPolls(&event)
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

func AssignUser(eventId string, userId string) error {
	event, err := event_repository.FindById(eventId)
	if err != nil {
		return err
	}
	user, err := user_repository.FindById(userId)
	if err != nil {
		return err
	}

	return event_repository.AppendAtendee(event, user)
}

func RemoveUser(eventId string, userId string) error {
	event, err := event_repository.FindById(eventId)
	if err != nil {
		return err
	}
	user, err := user_service.GetUser(userId)
	if err != nil {
		return err
	}

	return event_repository.RemoveAtendee(event, &entities.User{ID: user.Id})
}

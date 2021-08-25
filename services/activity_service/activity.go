package activity_service

import (
	"github.com/energy-uktc/eventpool-api/entities"
	"github.com/energy-uktc/eventpool-api/models"
	"github.com/energy-uktc/eventpool-api/repositories/activity_repository"
)

func Create(model *models.CreateUpdateActivity) (*models.Activity, error) {
	activity := new(entities.Activity)
	activity.ParseModel(model)
	id, err := activity_repository.Create(activity)
	if err != nil {
		return nil, err
	}
	activity, err = activity_repository.FindById(model.EventID, id)
	if err != nil {
		return nil, err
	}
	return activity.ToModel(), nil
}

func Update(id string, model *models.CreateUpdateActivity) (*models.Activity, error) {
	activity, err := activity_repository.FindById(model.EventID, id)
	if err != nil {
		return nil, err
	}
	activity.ParseModel(model)
	updatedActivity, err := activity_repository.Update(activity)

	if err != nil {
		return nil, err
	}

	return updatedActivity.ToModel(), nil
}

func FindAll(eventId string) ([]*models.Activity, error) {
	activities, err := activity_repository.FindForEvent(eventId)
	if err != nil {
		return nil, err
	}
	var activityModels []*models.Activity
	for _, activity := range activities {
		activityModel := activity.ToModel()
		activityModels = append(activityModels, activityModel)
	}
	return activityModels, nil
}

func FindById(eventId string, id string) (*models.Activity, error) {
	activity, err := activity_repository.FindById(eventId, id)
	if err != nil {
		return nil, err
	}
	return activity.ToModel(), nil
}

func Delete(eventId, id string) error {
	return activity_repository.Delete(eventId, id)
}

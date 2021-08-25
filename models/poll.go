package models

import (
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validatepolltype", validatePollType)
		v.RegisterValidation("validatepolloptions", validatePollOptions)
	}
}

type PollModel struct {
	ID        string             `json:"id"`
	CreatedBy *UserModel         `json:"createdBy"`
	Type      string             `json:"type" binding:"required,validatepolltype"`
	EndTime   time.Time          `json:"endTime" binding:"required,gt"`
	Question  string             `json:"question" binding:"required,min=5"`
	Options   []*PollAnswerModel `json:"options" binding:"validatepolloptions"`
}

type CreatePollModel struct {
	EventID   string
	CreatedBy string             `json:"createdBy"`
	Type      string             `json:"type" binding:"required,validatepolltype"`
	EndTime   time.Time          `json:"endTime" binding:"required,gt"`
	Question  string             `json:"question" binding:"required,min=5"`
	Options   []*PollAnswerModel `json:"options" binding:"validatepolloptions"`
}

type PollAnswerModel struct {
	ID            string             `json:"id"`
	Text          string             `json:"text" binding:"required"`
	NumberOfVotes int                `json:"numberOfVotes"`
	ShowOrder     uint8              `json:"showOrder" binding:"required"`
	Votes         []*SimpleUserModel `json:"votes"`
}

type UpdatePollModel struct {
	EndTime time.Time `json:"endTime" binding:"required,gt"`
}

var validatePollType validator.Func = func(fl validator.FieldLevel) bool {
	pollType := fl.Field().Interface().(string)
	switch pollType {
	case "SINGLE_OPTION", "MULTIPLE_OPTIONS":
		return true
	}
	return false
}

var validatePollOptions validator.Func = func(fl validator.FieldLevel) bool {
	options := fl.Field().Interface().([]*PollAnswerModel)
	if options == nil {
		return false
	}
	if len(options) < 2 || len(options) > 50 {
		return false
	}
	for _, option := range options {
		if len(option.Text) < 1 {
			return false
		}
	}
	return true
}

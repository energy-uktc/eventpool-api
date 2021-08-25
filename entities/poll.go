package entities

import (
	"sort"
	"time"

	"github.com/energy-uktc/eventpool-api/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type PollType uint

const (
	SINGE_OPTION PollType = iota + 1
	MULTIPLE_OPTIONS
)

var pollTypeOptions = [...]string{"SINGLE_OPTION", "MULTIPLE_OPTIONS"}

func (t PollType) String() string {
	return pollTypeOptions[t-1]
}

func ParsePollType(typeString string) PollType {
	for i, value := range pollTypeOptions {
		if value == typeString {
			return PollType(i + 1)
		}
	}
	return 0
}

type Poll struct {
	ID          string `gorm:"primarykey;type:uuid"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	EventID     string
	CreatedByID string   `gorm:"type:uuid"`
	CreatedBy   *User    `gorm:"foreignKey:CreatedByID"`
	Type        PollType `gorm:"check:type > 0;check:type < 3"`
	EndTime     time.Time
	Question    string
	Options     []PollAnswer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (p *Poll) BeforeCreate(tx *gorm.DB) (err error) {
	uid, _ := uuid.NewV4()
	p.ID = uid.String()
	return
}

func (p *Poll) ToModel() *models.PollModel {
	model := &models.PollModel{
		ID:        p.ID,
		CreatedBy: p.CreatedBy.ToModel(),
		Type:      p.Type.String(),
		EndTime:   p.EndTime,
		Question:  p.Question,
		Options:   make([]*models.PollAnswerModel, 0),
	}

	if p.Options != nil {
		for _, option := range p.Options {
			model.Options = append(model.Options, option.ToModel())
		}
		sort.Slice(p.Options, func(i, j int) bool {
			return p.Options[i].ShowOrder < p.Options[j].ShowOrder
		})
	}
	return model
}

func (p *Poll) ParseCreateModel(model *models.CreatePollModel) {
	p.Type = ParsePollType(model.Type)
	p.EventID = model.EventID
	p.CreatedByID = model.CreatedBy
	p.EndTime = model.EndTime
	p.Question = model.Question
	p.ParsePollAnswers(model.Options)
}

func (p *Poll) ParsePollAnswers(answers []*models.PollAnswerModel) {
	p.Options = make([]PollAnswer, 0)
	if answers == nil {
		return
	}
	for _, option := range answers {
		p.Options = append(p.Options, PollAnswer{
			Text: option.Text,
		})
	}
}

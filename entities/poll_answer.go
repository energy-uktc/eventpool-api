package entities

import (
	"time"

	"github.com/energy-uktc/eventpool-api/models"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type PollAnswer struct {
	ID        string `gorm:"primarykey;type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	PollID    string
	Text      string
	ShowOrder uint8
	Votes     []User `gorm:"many2many:user_poll_answers;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (p *PollAnswer) BeforeCreate(tx *gorm.DB) (err error) {
	uid, _ := uuid.NewV4()
	p.ID = uid.String()
	return
}

func (p *PollAnswer) ParseModel(model *models.PollAnswerModel) {
	p.Text = model.Text
	p.ShowOrder = model.ShowOrder
}

func (p *PollAnswer) ToModel() *models.PollAnswerModel {
	model := &models.PollAnswerModel{
		ID:            p.ID,
		Text:          p.Text,
		ShowOrder:     p.ShowOrder,
		NumberOfVotes: 0,
		Votes:         make([]*models.SimpleUserModel, 0),
	}
	if p.Votes != nil {
		model.NumberOfVotes = len(p.Votes)
		for _, vote := range p.Votes {
			model.Votes = append(model.Votes, vote.ToSimpleModel())

		}
	}
	return model
}

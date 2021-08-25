package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type TokenType uint8

const (
	RefreshToken TokenType = iota + 1
	VerificationCode
	ResetPasswordCode
)

func (t TokenType) String() string {
	return [...]string{"RefreshToken", "VerificationCode", "ResetPasswordCode"}[t-1]
}

type UserToken struct {
	ID        string `gorm:"primarykey;type:uuid"`
	CreatedAt time.Time
	ExpiresAt time.Time
	Type      TokenType `gorm:"index:idx_user_tokens_token,unique;not null;check:type > 0"`
	Token     string    `gorm:"index:idx_user_tokens_token,unique;not null"`
	UserID    string
	User      User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (t *UserToken) BeforeCreate(tx *gorm.DB) (err error) {
	uid, _ := uuid.NewV4()
	t.ID = uid.String()
	t.ExpiresAt = time.Now().Add(time.Hour * 24 * 7)
	return
}

package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	ExpiresAt time.Time
	User      User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

func (session *Session) BeforeCreate(tx *gorm.DB) error {
	session.ID = uuid.New()

	return nil
}

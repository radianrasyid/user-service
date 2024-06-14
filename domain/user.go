package domain

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/prithuadhikary/user-service/helper"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;primarykey"`
	Username string
	Password string
	Role     string
	Token    string
	Email    string
	Session  []Session `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.ID = uuid.New()

	userPassword, err := helper.HashPassword(user.Password)

	if err != nil {
		return fmt.Errorf("error on password hashing: %w", err)
	}
	// Digest and store the hex representation.
	user.Password = userPassword

	return nil
}

package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/prithuadhikary/user-service/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	Save(user *domain.User)
	ExistsByUsername(username string) bool
	FindSpecificUsername(username string) (domain.User, error)
	EditUser(username string, value any, column string) error
	FindUserBySessionID(sessionID uuid.UUID) (*domain.User, error)
	CreateSession(session *domain.Session) error
}

type userRepository struct {
	db *gorm.DB
}

func (repository *userRepository) Save(user *domain.User) {
	repository.db.Save(user)
}

func (repository *userRepository) ExistsByUsername(username string) bool {
	var count int64
	repository.db.Model(&domain.User{}).Where("username = ?", username).Count(&count)
	return count > 0
}

func (repository *userRepository) FindSpecificUsername(username string) (domain.User, error) {
	var user domain.User

	result := repository.db.Model(&domain.User{}).Where("username = ?", username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.User{}, fmt.Errorf("user with username '%s' not found", username)
		}

		return domain.User{}, fmt.Errorf("database error: %v", result.Error)
	}

	return user, nil
}

func (repository *userRepository) EditUser(username string, value any, column string) error {
	result := repository.db.Model(&domain.User{}).Where("username = ?", username).Update(column, value)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrInvalidValue) {
			return fmt.Errorf("invalid value '%s'", value)
		}

		return fmt.Errorf("database error: %v", result.Error)
	}

	return nil
}

func (repository *userRepository) FindUserBySessionID(sessionID uuid.UUID) (*domain.User, error) {
	var session domain.Session
	err := repository.db.Preload("User").First(&session, `id = ?`, sessionID).Error
	if err != nil {
		return nil, err
	}

	return &session.User, nil
}

func (repository *userRepository) CreateSession(session *domain.Session) error {
	return repository.db.Create(session).Error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	var repository UserRepository

	repository = &userRepository{
		db: db,
	}

	return repository
}

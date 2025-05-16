package database

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
	"errors"

	"gorm.io/gorm"
)

type PostgreSQLUserSessionRepository struct {
	db *gorm.DB
}

func NewPostgreSQLUserSessionRepository(db *gorm.DB) repository.UsersSessionRepository {
	return &PostgreSQLUserSessionRepository{db: db}
}

func (r *PostgreSQLUserSessionRepository) Create(usersSessions *entity.UsersSesions, tx *gorm.DB) error {
	result := tx.Create(usersSessions)
	return result.Error
}

func (r *PostgreSQLUserSessionRepository) FindByRefreshToken(refreshToken string) (*entity.UsersSesions, error) {
	var userSession = &entity.UsersSesions{}
	result := r.db.Model(&entity.UsersSesions{}).Where("refresh_token = ?", refreshToken).First(userSession)
	if result.Error != nil {
		return nil, result.Error
	}
	return userSession, nil
}

func (r *PostgreSQLUserSessionRepository) UpdateIsRevoked(user *entity.UsersSesions) error {
	result := r.db.Model(&entity.UsersSesions{}).Where("id = ?", user.Id).Update("is_revoked", true)
	return result.Error
}

func (r *PostgreSQLUserSessionRepository) CheckUserInSession(userId int64) bool {
	var userSessions = &entity.UsersSesions{}
	result := r.db.Model(&entity.UsersSesions{}).Where("user_id = ? and is_revoked = ?", userId, false).First(userSessions)
	return !errors.Is(result.Error, gorm.ErrRecordNotFound)
}

func (r *PostgreSQLUserSessionRepository) FindByUserIdInSession(userId int64) (*entity.UsersSesions, error) {
	var userSessions = &entity.UsersSesions{}
	result := r.db.Model(&entity.UsersSesions{}).Where("user_id = ? and is_revoked = ?", userId, false).First(userSessions)
	if result != nil {
		return nil, result.Error
	}
	return userSessions, nil
}

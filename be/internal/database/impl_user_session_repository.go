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

func (r *PostgreSQLUserSessionRepository) Create(usersSessions *entity.UsersSessions, tx *gorm.DB) error {
	result := tx.Create(usersSessions)
	return result.Error
}

func (r *PostgreSQLUserSessionRepository) FindByRefreshToken(refreshToken string) (*entity.UsersSessions, error) {
	var userSession = &entity.UsersSessions{}
	result := r.db.Model(&entity.UsersSessions{}).Where("refresh_token = ?", refreshToken).First(userSession)
	if result.Error != nil {
		return nil, result.Error
	}
	return userSession, nil
}

func (r *PostgreSQLUserSessionRepository) UpdateIsRevoked(user *entity.UsersSessions) error {
	result := r.db.Model(&entity.UsersSessions{}).Where("id = ?", user.Id).Update("is_revoked", true)
	return result.Error
}

func (r *PostgreSQLUserSessionRepository) CheckUserInSession(userId int64) bool {
	var userSessions = &entity.UsersSessions{}
	result := r.db.Model(&entity.UsersSessions{}).Where("user_id = ? and is_revoked = ?", userId, false).First(userSessions)
	return !errors.Is(result.Error, gorm.ErrRecordNotFound)
}

func (r *PostgreSQLUserSessionRepository) FindByUserIdInSession(userId int64) (*entity.UsersSessions, error) {
	var userSessions = &entity.UsersSessions{}
	result := r.db.Model(&entity.UsersSessions{}).Where("user_id = ? and is_revoked = ?", userId, false).First(userSessions)
	if result.Error != nil {
		return nil, result.Error
	}
	return userSessions, nil
}

func (r *PostgreSQLUserSessionRepository) CheckTokenWasInVoked(accessToken string) bool {
	var userSessions = &entity.UsersSessions{}
	r.db.Model(&entity.UsersSessions{}).Where("access_token = ?", accessToken).First(userSessions)
	return userSessions.IsRevoked
}

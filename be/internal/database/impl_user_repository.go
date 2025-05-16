package database

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
	"errors"

	"gorm.io/gorm"
)

type PostgreSQLUserRepository struct {
	db *gorm.DB
}

func NewPostgreSQLUserRepository(db *gorm.DB) repository.UserRepository {
	return &PostgreSQLUserRepository{db: db}
}

func (r *PostgreSQLUserRepository) Create(users *entity.Users) error {
	if users.FirstName == "" {
		return errors.New("name can't not blank")
	}
	if users.LastName == "" {
		return errors.New("name can't not blank")
	}
	if users.Password == "" {
		return errors.New("password can't not blank")
	}
	if err := r.db.Create(users).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostgreSQLUserRepository) FindByToken(token string) (*entity.Users, error) {
	var users = &entity.Users{}
	result := r.db.Model(&entity.Users{}).Where("token = ?", token).Find(users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *PostgreSQLUserRepository) Update(users *entity.Users) error {
	result := r.db.Model(&entity.Users{}).Where("email = ?", users.Email).Update("is_active", true)
	return result.Error
}

func (r *PostgreSQLUserRepository) UpdatePassword(users *entity.Users) error {
	result := r.db.Model(&entity.Users{}).Where("email = ?", users.Email).Update("password", users.Password)
	return result.Error
}

func (r *PostgreSQLUserRepository) FindByEmail(email string) (*entity.Users, error) {
	users := &entity.Users{}
	result := r.db.Model(entity.Users{}).Where("email = ?", email).First(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *PostgreSQLUserRepository) FindByUserId(userId int64) (*entity.Users, error) {
	users := &entity.Users{}
	result := r.db.Model(entity.Users{}).Where("id = ?", userId).First(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *PostgreSQLUserRepository) DeleteUser(email string) error {
	result := r.db.Where("email = ?", email).Delete(&entity.Users{})
	return result.Error
}

func (r *PostgreSQLUserRepository) GetDB() *gorm.DB {
	return r.db
}

package database

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"

	"gorm.io/gorm"
)

type PostgreSQLCategoriesRepository struct {
	db *gorm.DB
}

func NewPostgreSQLCategoriesRepository(db *gorm.DB) repository.CategoriesRepository {
	return &PostgreSQLCategoriesRepository{db: db}
}

func (r *PostgreSQLCategoriesRepository) Create(Category *entity.Categories) (*entity.Categories, error) {
	result := r.db.Create(Category)
	return Category, result.Error
}

func (r *PostgreSQLCategoriesRepository) GetAll() ([]*entity.Categories, error) {
	categories := []*entity.Categories{}
	result := r.db.Model(entity.Categories{}).Find(categories)
	return categories, result.Error
}

func (r *PostgreSQLCategoriesRepository) Delete(id int64) error {
	result := r.db.Model(entity.Categories{}).Where("id = ?", id).Delete(entity.Categories{})
	return result.Error
}

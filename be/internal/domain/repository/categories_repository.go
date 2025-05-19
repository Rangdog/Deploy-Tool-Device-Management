package repository

import "BE_Manage_device/internal/domain/entity"

type CategoriesRepository interface {
	Create(*entity.Categories) (*entity.Categories, error)
	GetAll() ([]*entity.Categories, error)
	Delete(id int64) error
}

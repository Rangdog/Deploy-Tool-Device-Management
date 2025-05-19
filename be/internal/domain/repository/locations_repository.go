package repository

import "BE_Manage_device/internal/domain/entity"

type LocationRepository interface {
	Create(*entity.Locations) (*entity.Locations, error)
	GetAll() ([]*entity.Locations, error)
	Delete(id int64) error
}

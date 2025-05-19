package database

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"

	"gorm.io/gorm"
)

type PostgreSQLLocationRepository struct {
	db *gorm.DB
}

func NewPostgreSQLLocationRepository(db *gorm.DB) repository.LocationRepository {
	return &PostgreSQLLocationRepository{db: db}
}

func (r *PostgreSQLLocationRepository) Create(location *entity.Locations) (*entity.Locations, error) {
	result := r.db.Create(location)
	return location, result.Error
}

func (r *PostgreSQLLocationRepository) GetAll() ([]*entity.Locations, error) {
	locations := []*entity.Locations{}
	result := r.db.Model(entity.Locations{}).Find(locations)
	return locations, result.Error
}

func (r *PostgreSQLLocationRepository) Delete(id int64) error {
	result := r.db.Model(entity.Locations{}).Where("id = ?", id).Delete(entity.Locations{})
	return result.Error
}

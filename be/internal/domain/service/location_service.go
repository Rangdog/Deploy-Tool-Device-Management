package service

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
)

type LocationService struct {
	repo repository.LocationRepository
}

func NewLocationService(repo repository.LocationRepository) *LocationService {
	return &LocationService{repo: repo}
}

func (service *LocationService) Create(LocationName string) (*entity.Locations, error) {
	var location = &entity.Locations{
		LocationName: LocationName,
	}
	locationCreate, err := service.repo.Create(location)
	if err != nil {
		return nil, err
	}
	return locationCreate, nil
}

func (service *LocationService) GetAll() ([]*entity.Locations, error) {
	locations, err := service.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return locations, err
}

func (service *LocationService) Delete(id int64) error {
	err := service.repo.Delete(id)
	return err
}

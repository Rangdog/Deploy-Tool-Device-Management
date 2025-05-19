package service

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
)

type DepartmentsService struct {
	repo repository.DepartmentsRepository
}

func NewDepartmentsService(repo repository.DepartmentsRepository) *DepartmentsService {
	return &DepartmentsService{repo: repo}
}

func (service *DepartmentsService) Create(departmentsName string, locationId int64) (*entity.Departments, error) {
	var departments = &entity.Departments{
		DepartmentName: departmentsName,
		LocationId:     locationId,
	}
	departmentsCreate, err := service.repo.Create(departments)
	if err != nil {
		return nil, err
	}
	return departmentsCreate, nil
}

func (service *DepartmentsService) GetAll() ([]*entity.Departments, error) {
	departments, err := service.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return departments, err
}

func (service *DepartmentsService) Delete(id int64) error {
	err := service.repo.Delete(id)
	return err
}

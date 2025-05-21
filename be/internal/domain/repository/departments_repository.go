package repository

import "BE_Manage_device/internal/domain/entity"

type DepartmentsRepository interface {
	Create(*entity.Departments) (*entity.Departments, error)
	GetAll() ([]*entity.Departments, error)
	Delete(id int64) error
	GetDepartmentById(id int64) (*entity.Departments, error)
}

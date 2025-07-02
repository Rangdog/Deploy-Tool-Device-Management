package repository

import "BE_Manage_device/internal/domain/entity"

type CompanyRepository interface {
	Create(*entity.Company) (*entity.Company, error)
	GetCompanyById(id int64) (*entity.Company, error)
	GetCompanyBySuffixEmail(email string) (*entity.Company, error)
	GetAllCompany() ([]*entity.Company, error)
}

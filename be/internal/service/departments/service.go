package service

import (
	"BE_Manage_device/internal/domain/entity"
	company "BE_Manage_device/internal/repository/company"
	department "BE_Manage_device/internal/repository/departments"
	user "BE_Manage_device/internal/repository/user"
	"BE_Manage_device/pkg/utils"
)

type DepartmentsService struct {
	repo        department.DepartmentsRepository
	userRepo    user.UserRepository
	companyRepo company.CompanyRepository
}

func NewDepartmentsService(repo department.DepartmentsRepository, userRepo user.UserRepository, companyRepo company.CompanyRepository) *DepartmentsService {
	return &DepartmentsService{repo: repo, userRepo: userRepo, companyRepo: companyRepo}
}

func (service *DepartmentsService) Create(userId int64, departmentsName string, locationId int64) (*entity.Departments, error) {
	user, err := service.userRepo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	company, err := service.companyRepo.GetCompanyBySuffixEmail(utils.GetSuffixEmail(user.Email))
	if err != nil {
		return nil, err
	}
	var departments = &entity.Departments{
		DepartmentName: departmentsName,
		LocationId:     locationId,
		CompanyId:      company.Id,
	}
	departmentsCreate, err := service.repo.Create(departments)
	if err != nil {
		return nil, err
	}
	return departmentsCreate, nil
}

func (service *DepartmentsService) GetAll(userId int64) ([]*entity.Departments, error) {
	user, err := service.userRepo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	departments, err := service.repo.GetAll(user.CompanyId)
	if err != nil {
		return nil, err
	}
	return departments, err
}

func (service *DepartmentsService) Delete(id int64) error {
	err := service.repo.Delete(id)
	return err
}

func (service *DepartmentsService) GetCompanyId(userId int64) (int64, error) {
	user, err := service.userRepo.FindByUserId(userId)
	return user.CompanyId, err
}

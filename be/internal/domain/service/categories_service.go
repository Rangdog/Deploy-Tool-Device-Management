package service

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
)

type CategoriesService struct {
	repo repository.CategoriesRepository
}

func NewCategoriesService(repo repository.CategoriesRepository) *CategoriesService {
	return &CategoriesService{repo: repo}
}

func (service *CategoriesService) Create(categoryName string) (*entity.Categories, error) {
	var category = &entity.Categories{
		CategoryName: categoryName,
	}
	categoryCreate, err := service.repo.Create(category)
	if err != nil {
		return nil, err
	}
	return categoryCreate, nil
}

func (service *CategoriesService) GetAll() ([]*entity.Categories, error) {
	categories, err := service.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return categories, err
}

func (service *CategoriesService) Delete(id int64) error {
	err := service.repo.Delete(id)
	return err
}

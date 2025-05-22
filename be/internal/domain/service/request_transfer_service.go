package service

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
)

type RequestTransferService struct {
	repo repository.RequestTransferRepository
}

func NewRequestTransferService(repo repository.RequestTransferRepository) *RequestTransferService {
	return &RequestTransferService{repo: repo}
}

func (service *RequestTransferService) Create(userId int64, assetId int64, departmentId int64) (*entity.RequestTransfer, error) {
	requestTransfer := entity.RequestTransfer{
		UserId:       userId,
		AssetId:      assetId,
		DepartmentId: departmentId,
		Status:       "Pending",
	}
	requestTransferCreate, err := service.repo.Create(&requestTransfer)
	if err != nil {
		return nil, err
	}
	return requestTransferCreate, err
}

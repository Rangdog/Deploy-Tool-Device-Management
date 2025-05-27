package service

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
)

type RequestTransferService struct {
	repo              repository.RequestTransferRepository
	assignmentService *AssignmentService
	userRepo          repository.UserRepository
}

func NewRequestTransferService(repo repository.RequestTransferRepository, assignmentService *AssignmentService, userRepo repository.UserRepository) *RequestTransferService {
	return &RequestTransferService{repo: repo, assignmentService: assignmentService, userRepo: userRepo}
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

func (service *RequestTransferService) Accept(userId int64, id int64) (*entity.RequestTransfer, error) {
	tx := service.repo.GetDB().Begin()
	request, err := service.repo.UpdateStatusConfirm(id, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	assignment, err := service.assignmentService.repo.GetAssignmentByAssetId(request.AssetId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	userAssign, err := service.userRepo.GetUserAssetManageOfDepartment(request.DepartmentId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = service.assignmentService.Update(userId, request.AssetId, assignment.Id, &userAssign.Id, &request.DepartmentId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return request, nil
}

func (service *RequestTransferService) Deny(userId int64, id int64) (*entity.RequestTransfer, error) {
	request, err := service.repo.UpdateStatusDeny(id)
	if err != nil {
		return nil, err
	}
	return request, nil
}

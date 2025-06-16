package service

import (
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/filter"
	"BE_Manage_device/internal/domain/repository"
	"BE_Manage_device/pkg/utils"
	"errors"
	"fmt"
)

type RequestTransferService struct {
	repo              repository.RequestTransferRepository
	assignmentService *AssignmentService
	userRepo          repository.UserRepository
	assetRepo         repository.AssetsRepository
}

func NewRequestTransferService(repo repository.RequestTransferRepository, assignmentService *AssignmentService, userRepo repository.UserRepository, assetRepo repository.AssetsRepository) *RequestTransferService {
	return &RequestTransferService{repo: repo, assignmentService: assignmentService, userRepo: userRepo, assetRepo: assetRepo}
}

func (service *RequestTransferService) Create(userId int64, categoryId int64, description string) (*entity.RequestTransfer, error) {
	requestTransfer := entity.RequestTransfer{
		UserId:      userId,
		CategoryId:  categoryId,
		Status:      "Pending",
		Description: description,
	}
	requestTransferCreate, err := service.repo.Create(&requestTransfer)
	if err != nil {
		return nil, err
	}
	return requestTransferCreate, err
}

func (service *RequestTransferService) Accept(userId int64, id int64, assetId int64) (*entity.RequestTransfer, error) {
	var err error
	if service.IsRoleHeadDep(userId) {
		return nil, errors.New("head Department can't accept request")
	}
	requestCheck, err := service.repo.GetRequestTransferById(id)
	if err != nil {
		return nil, err
	}
	if requestCheck.Status == "Deny" {
		return nil, errors.New("can't change request")
	}
	assetCheck, err := service.assetRepo.GetAssetById(assetId)
	if err != nil {
		return nil, err
	}
	if assetCheck.DepartmentId == *requestCheck.User.DepartmentId {
		return nil, errors.New("asset department same request department")
	}
	tx := service.repo.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()
	request, err := service.repo.UpdateStatusConfirm(id, tx)
	if err != nil {
		return nil, err
	}
	assignment, err := service.assignmentService.repo.GetAssignmentByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	if requestCheck.User.DepartmentId == nil {
		return nil, errors.New("user don't have department")
	}
	userAssign, err := service.userRepo.GetUserAssetManageOfDepartment(*requestCheck.User.DepartmentId)
	if err != nil {
		return nil, err
	}
	_, err = service.assignmentService.Update(userId, assignment.Id, &userAssign.Id, requestCheck.User.DepartmentId)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}
	return request, nil
}

func (service *RequestTransferService) Deny(userId int64, id int64) (*entity.RequestTransfer, error) {
	var err error
	if service.IsRoleHeadDep(userId) {
		return nil, errors.New("head Department can't deny request")
	}
	requestCheck, err := service.repo.GetRequestTransferById(id)
	if err != nil {
		return nil, err
	}
	if requestCheck.Status == "Confirm" {
		return nil, errors.New("can't change request")
	}
	tx := service.repo.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()
	request, err := service.repo.UpdateStatusDeny(id, tx)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}
	return request, nil
}

func (service *RequestTransferService) GetRequestTransferById(userId int64, id int64) (*entity.RequestTransfer, error) {
	request, err := service.repo.GetRequestTransferById(id)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (service *RequestTransferService) Filter(userId int64, status *string) ([]dto.RequestTransferResponse, error) {
	users, err := service.userRepo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	var filter = filter.RequestTransferFilter{
		Status: status,
	}
	if users.Role.Slug != "admin" {
		filter.DepId = users.DepartmentId
	}
	db := service.repo.GetDB()
	dbFilter := filter.ApplyFilter(db.Model(&entity.RequestTransfer{}), userId)
	var requests []entity.RequestTransfer
	result := dbFilter.Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	requestRes := utils.ConvertRequestTransfersToResponses(requests)

	return requestRes, nil
}

func (service *RequestTransferService) IsRoleHeadDep(userId int64) bool {
	user, err := service.userRepo.FindByUserId(userId)
	if err != nil {
		return false
	}
	return user.IsHeadDepartment

}

package service

import (
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/filter"
	"BE_Manage_device/internal/domain/repository"
	"errors"
	"math"
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
	requestCheck, err := service.repo.GetRequestTransferById(id)
	if err != nil {
		return nil, err
	}
	if requestCheck.Status == "Deny" {
		return nil, errors.New("can't change request")
	}
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
	_, err = service.assignmentService.Update(userId, assignment.Id, &userAssign.Id, &request.DepartmentId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return request, nil
}

func (service *RequestTransferService) Deny(userId int64, id int64) (*entity.RequestTransfer, error) {
	requestCheck, err := service.repo.GetRequestTransferById(id)
	if err != nil {
		return nil, err
	}
	if requestCheck.Status == "Confirm" {
		return nil, errors.New("can't change request")
	}
	tx := service.repo.GetDB().Begin()
	request, err := service.repo.UpdateStatusDeny(id, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
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

func (service *RequestTransferService) Filter(userId int64, status *string, page int, limit int) (*map[string]any, error) {
	var filter = filter.RequestTransferFilter{
		Status: status,
		Page:   page,
		Limit:  limit,
	}
	db := service.repo.GetDB()
	dbFilter := filter.ApplyFilter(db.Model(&entity.RequestTransfer{}), userId)
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	var total int64
	dbFilter.Count(&total)
	offset := (filter.Page - 1) * filter.Limit
	var requests []entity.RequestTransfer
	result := dbFilter.Offset(offset).Limit(filter.Limit).Find(&requests)
	if result.Error != nil {
		return nil, result.Error
	}
	var requestRes []dto.RequestTransferResponse
	for _, request := range requests {
		requestTransferResponse := dto.RequestTransferResponse{}
		requestTransferResponse.Id = request.Id
		requestTransferResponse.Status = request.Status
		requestTransferResponse.User.Id = request.User.Id
		requestTransferResponse.User.FirstName = request.User.FirstName
		requestTransferResponse.User.LastName = request.User.LastName
		requestTransferResponse.User.Email = request.User.Email
		requestTransferResponse.Asset.Id = request.Asset.Id
		requestTransferResponse.Asset.AssetName = request.Asset.AssetName
		requestTransferResponse.Asset.SerialNumber = request.Asset.SerialNumber
		requestTransferResponse.Asset.ImageUpload = *request.Asset.ImageUpload
		requestTransferResponse.Asset.FileAttachment = *request.Asset.FileAttachment
		requestTransferResponse.Asset.QrUrl = *request.Asset.QrUrl
		requestTransferResponse.Department.Id = request.DepartmentId
		requestTransferResponse.Department.DepartmentName = request.Department.DepartmentName
		requestTransferResponse.Department.Location.ID = request.Department.Location.Id
		requestTransferResponse.Department.Location.LocationName = request.Department.Location.LocationName
		requestRes = append(requestRes, requestTransferResponse)
	}
	data := map[string]any{
		"data":       requestRes,
		"page":       filter.Page,
		"limit":      filter.Limit,
		"total":      total,
		"total_page": int(math.Ceil(float64(total) / float64(filter.Limit))),
	}
	return &data, nil
}

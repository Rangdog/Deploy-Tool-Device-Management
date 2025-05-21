package service

import (
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/filter"
	"BE_Manage_device/internal/domain/repository"
	"math"
)

type AssignmentService struct {
	repo repository.AssignmentRepository
}

func NewAssignmentService(repo repository.AssignmentRepository) *AssignmentService {
	return &AssignmentService{repo: repo}
}

func (service *AssignmentService) Create(userId, userIdAssign, assetId int64) (*entity.Assignments, error) {
	assignment := entity.Assignments{
		UserId:   &userIdAssign,
		AssetId:  &assetId,
		AssignBy: userId,
	}
	assignmentCreated, err := service.repo.Create(&assignment)
	if err != nil {
		return nil, err
	}
	return assignmentCreated, err
}

func (service *AssignmentService) Update(userId, userIdAssign, assetId, assignmentId int64) (*entity.Assignments, error) {
	assigmentUpdated, err := service.repo.Update(assignmentId, &userIdAssign, &assetId, userId)
	if err != nil {
		return nil, err
	}
	return assigmentUpdated, nil
}

func (service *AssignmentService) Filter(userId int64, emailAssigned *string, emailAssign *string, assetName *string, page int, limit int) (*map[string]any, error) {
	var filter = filter.AssignmentFilter{
		EmailAssigned: emailAssigned,
		EmailAssign:   emailAssign,
		AssetName:     assetName,
		Page:          page,
		Limit:         limit,
	}
	db := service.repo.GetDB()
	dbFilter := filter.ApplyFilter(db.Model(&entity.Assignments{}), userId)
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	var total int64
	dbFilter.Count(&total)
	offset := (filter.Page - 1) * filter.Limit
	var assignments []entity.Assignments
	resutl := dbFilter.Offset(offset).Limit(filter.Limit).Find(&assignments)
	if resutl.Error != nil {
		return nil, resutl.Error
	}
	assignmentsRes := []dto.AssignmentResponse{}
	for _, assignment := range assignments {
		assignResponse := dto.AssignmentResponse{}
		assignResponse.Id = assignment.Id
		assignResponse.UserAssigned.Id = assignment.UserAssigned.Id
		assignResponse.UserAssigned.FirstName = assignment.UserAssigned.FirstName
		assignResponse.UserAssigned.LastName = assignment.UserAssigned.LastName
		assignResponse.UserAssigned.Email = assignment.UserAssigned.Email

		assignResponse.UserAssign.Id = assignment.UserAssign.Id
		assignResponse.UserAssign.FirstName = assignment.UserAssign.FirstName
		assignResponse.UserAssign.LastName = assignment.UserAssign.LastName
		assignResponse.UserAssign.Email = assignment.UserAssign.Email

		assignResponse.Asset.Id = assignment.Asset.Id
		assignResponse.Asset.AssetName = assignment.Asset.AssetName
		assignResponse.Asset.Status = assignment.Asset.Status
		assignResponse.Asset.FileAttachment = *assignment.Asset.FileAttachment
		assignResponse.Asset.ImageUpload = *assignment.Asset.ImageUpload
		assignmentsRes = append(assignmentsRes, assignResponse)
	}
	data := map[string]any{
		"data":       assignmentsRes,
		"page":       filter.Page,
		"limit":      filter.Limit,
		"total":      total,
		"total_page": int(math.Ceil(float64(total) / float64(filter.Limit))),
	}
	return &data, nil
}

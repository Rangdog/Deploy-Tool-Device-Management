package service

import (
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/filter"
	"BE_Manage_device/internal/domain/repository"
	"errors"
	"fmt"
	"math"
	"time"
)

type AssignmentService struct {
	repo           repository.AssignmentRepository
	assetLogRepo   repository.AssetsLogRepository
	assetRepo      repository.AssetsRepository
	departmentRepo repository.DepartmentsRepository
	userRepo       repository.UserRepository
}

func NewAssignmentService(repo repository.AssignmentRepository, assetLogRepo repository.AssetsLogRepository, assetRepo repository.AssetsRepository, departmentRepo repository.DepartmentsRepository, userRepo repository.UserRepository) *AssignmentService {
	return &AssignmentService{repo: repo, assetLogRepo: assetLogRepo, assetRepo: assetRepo, departmentRepo: departmentRepo, userRepo: userRepo}
}

func (service *AssignmentService) Create(userIdAssign, departmentId *int64, userId, assetId int64) (*entity.Assignments, error) {
	assignment := entity.Assignments{
		UserId:       userIdAssign,
		AssetId:      assetId,
		AssignBy:     userId,
		DepartmentId: departmentId,
	}
	tx := service.repo.GetDB().Begin()
	assignmentCreated, err := service.repo.Create(&assignment, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return assignmentCreated, err
}

func (service *AssignmentService) Update(userId, assignmentId int64, userIdAssign, departmentId *int64) (*entity.Assignments, error) {
	assignment, err := service.repo.GetAssignmentById(assignmentId)
	if err != nil {
		return nil, err
	}

	asset, err := service.assetRepo.GetAssetById(assignment.AssetId)
	if err != nil {
		return nil, err
	}

	tx := service.repo.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	assignmentUpdated, err := service.repo.Update(assignmentId, userId, assignment.AssetId, userIdAssign, departmentId, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	assetLog := entity.AssetLog{
		Timestamp: time.Now(),
		Action:    "Transfer",
		AssetId:   assignment.AssetId,
	}

	oldAssetLog, err := service.assetLogRepo.GetNewLogByAssetId(assignment.AssetId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Kiểm tra nil trước khi dereference
	if oldAssetLog != nil && departmentId != nil && userIdAssign != nil &&
		oldAssetLog.DepartmentAssignedId != nil && *oldAssetLog.DepartmentAssignedId == *departmentId &&
		oldAssetLog.UserAssignedId == *userIdAssign {
		tx.Commit()
		return assignmentUpdated, nil
	}
	if departmentId != nil && (*departmentId == asset.DepartmentId) && userIdAssign == nil {
		tx.Commit()
		return assignmentUpdated, nil
	}
	if userIdAssign != nil && (*userIdAssign == asset.OnwerUser.Id) && departmentId == nil {
		tx.Commit()
		return assignmentUpdated, nil
	}
	if (departmentId == nil && userIdAssign == nil) ||
		(userIdAssign != nil && asset.OnwerUser != nil && *userIdAssign == asset.OnwerUser.Id &&
			(departmentId == nil || (*departmentId == asset.DepartmentId))) {
		tx.Rollback()
		return nil, errors.New("invalid request: cannot assign to current owner with same department or missing info")
	}

	// Chuyển phòng ban
	if departmentId != nil && (*departmentId != asset.DepartmentId) {
		department, err := service.departmentRepo.GetDepartmentById(*departmentId)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		assetLog.DepartmentAssignedId = departmentId
		if userIdAssign != nil {
			assetLog.UserAssignedId = *userIdAssign
		} else {
			user, err := service.userRepo.GetUserAssetManageOfDepartment(*departmentId)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			assignmentUpdated, err = service.repo.Update(assignmentId, userId, assignment.AssetId, &user.Id, departmentId, tx)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			assetLog.UserAssignedId = oldAssetLog.UserAssignedId
		}
		assetLog.ChangeSummary = fmt.Sprintf("Transfer from department %v to department %v\n",
			asset.Department.DepartmentName, department.DepartmentName)
		if _, err := service.assetRepo.UpdateAssetDepartment(assignment.AssetId, *departmentId); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Chuyển người dùng
	if userIdAssign != nil && (*userIdAssign != asset.OnwerUser.Id) {
		user, err := service.userRepo.FindByUserId(*userIdAssign)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		assetLog.UserAssignedId = *userIdAssign
		if departmentId != nil {
			assetLog.DepartmentAssignedId = departmentId
		} else {
			assetLog.DepartmentAssignedId = oldAssetLog.DepartmentAssignedId
		}
		assetLog.ChangeSummary += fmt.Sprintf("Transfer from user: %v to user: %v\n",
			asset.OnwerUser.Email, user.Email)
		if _, err := service.assetRepo.UpdateAssetOwner(assignment.AssetId, *userIdAssign); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if _, err := service.assetLogRepo.Create(&assetLog, tx); err != nil {
		tx.Rollback()
		return nil, err
	}
	if _, err := service.assetRepo.UpdateAssetLifeCycleStage(assignment.AssetId, "In Use", tx); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := service.assetRepo.UpdateOwner(assignment.AssetId, assetLog.UserAssignedId, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return assignmentUpdated, nil
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
	result := dbFilter.Offset(offset).Limit(filter.Limit).Find(&assignments)
	if result.Error != nil {
		return nil, result.Error
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
		assignResponse.Department.ID = *assignment.DepartmentId
		assignResponse.Department.DepartmentName = assignment.Department.DepartmentName
		assignResponse.Department.Location.ID = assignment.Department.LocationId
		assignResponse.Department.Location.LocationName = assignment.Department.Location.LocationName
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

func (service *AssignmentService) GetAssignmentById(userId int64, id int64) (*dto.AssignmentResponse, error) {
	assignment, err := service.repo.GetAssignmentById(id)
	if err != nil {
		return nil, err
	}
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
	assignResponse.Department.ID = *assignment.DepartmentId
	assignResponse.Department.DepartmentName = assignment.Department.DepartmentName
	assignResponse.Department.Location.ID = assignment.Department.LocationId
	assignResponse.Department.Location.LocationName = assignment.Department.Location.LocationName
	return &assignResponse, nil
}

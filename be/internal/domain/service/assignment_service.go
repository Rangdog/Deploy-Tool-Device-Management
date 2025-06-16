package service

import (
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/filter"
	"BE_Manage_device/internal/domain/repository"
	"BE_Manage_device/pkg/utils"
	"fmt"
	"time"
)

type AssignmentService struct {
	repo                repository.AssignmentRepository
	assetLogRepo        repository.AssetsLogRepository
	assetRepo           repository.AssetsRepository
	departmentRepo      repository.DepartmentsRepository
	userRepo            repository.UserRepository
	NotificationService *NotificationService
}

func NewAssignmentService(repo repository.AssignmentRepository, assetLogRepo repository.AssetsLogRepository, assetRepo repository.AssetsRepository, departmentRepo repository.DepartmentsRepository, userRepo repository.UserRepository, NotificationService *NotificationService) *AssignmentService {
	return &AssignmentService{repo: repo, assetLogRepo: assetLogRepo, assetRepo: assetRepo, departmentRepo: departmentRepo, userRepo: userRepo, NotificationService: NotificationService}
}

func (service *AssignmentService) Create(userIdAssign, departmentId *int64, userId, assetId int64) (*entity.Assignments, error) {
	var err error
	assignment := entity.Assignments{
		UserId:       userIdAssign,
		AssetId:      assetId,
		AssignBy:     userId,
		DepartmentId: departmentId,
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
	assignmentCreated, err := service.repo.Create(&assignment, tx)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return assignmentCreated, err
}

func (service *AssignmentService) Update(userId, assignmentId int64, userIdAssign, departmentId *int64) (*entity.Assignments, error) {
	var err error
	assignment, err := service.repo.GetAssignmentById(assignmentId)
	if err != nil {
		return nil, err
	}
	byUser, err := service.userRepo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	var assignUser *entity.Users
	if userIdAssign != nil {
		assignUser, err = service.userRepo.FindByUserId(*userIdAssign)
		if err != nil {
			return nil, err
		}
	} else {
		assignUser = nil
	}
	asset, err := service.assetRepo.GetAssetById(assignment.AssetId)
	if err != nil {
		return nil, err
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
	var assignmentUpdated *entity.Assignments
	if departmentId != nil {
		assignmentUpdated, err = service.repo.Update(assignmentId, userId, assignment.AssetId, userIdAssign, departmentId, tx)
	} else {
		assignmentUpdated, err = service.repo.Update(assignmentId, userId, assignment.AssetId, userIdAssign, assignUser.DepartmentId, tx)
	}
	if err != nil {
		return nil, err
	}

	assetLog := entity.AssetLog{
		Timestamp: time.Now(),
		Action:    "Transfer",
		AssetId:   asset.Id,
		ByUserId:  &byUser.Id,
	}

	// Chuyển phòng ban
	if departmentId != nil && (*departmentId != asset.DepartmentId) {
		department, err := service.departmentRepo.GetDepartmentById(*departmentId)
		if err != nil {
			return nil, err
		}
		assetLog.ChangeSummary = fmt.Sprintf("Transfer from department %v to department %v by user %v\n",
			asset.Department.DepartmentName, department.DepartmentName, byUser.Email)
		if _, err := service.assetRepo.UpdateAssetDepartment(assignment.AssetId, *departmentId, tx); err != nil {
			return nil, err
		}
		if _, err := service.assetLogRepo.Create(&assetLog, tx); err != nil {
			return nil, err
		}
	}

	// Chuyển người dùng
	if userIdAssign != nil && (*userIdAssign != asset.OnwerUser.Id) {
		assetLog.AssignUserId = &assignUser.Id
		assetLog.ChangeSummary += fmt.Sprintf("Transfer from user: %v to user: %v\n",
			byUser.Email, assignUser.Email)
		if _, err := service.assetRepo.UpdateAssetOwner(assignment.AssetId, *userIdAssign, tx); err != nil {
			return nil, err
		}
		if err := service.assetRepo.UpdateOwner(assignment.AssetId, assignUser.Id, tx); err != nil {
			return nil, err
		}
		if *assignUser.DepartmentId != asset.DepartmentId {
			_, err := service.assetRepo.UpdateAssetDepartment(asset.Id, *assignUser.DepartmentId, tx)
			if err != nil {
				return nil, err
			}
		}
		if _, err := service.assetLogRepo.Create(&assetLog, tx); err != nil {
			return nil, err
		}
	}

	if _, err := service.assetRepo.UpdateAssetLifeCycleStage(assignment.AssetId, "In Use", tx); err != nil {
		return nil, err
	}
	err = service.assetRepo.UpdateAcquisitionDate(assignment.AssetId, time.Now(), tx)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	userHeadDepart, _ := service.userRepo.GetUserHeadDepartment(assetLog.Asset.DepartmentId)
	userManagerAsset, _ := service.userRepo.GetUserAssetManageOfDepartment(assetLog.Asset.DepartmentId)
	usersToNotifications := []*entity.Users{asset.OnwerUser, userHeadDepart, userManagerAsset}
	message := fmt.Sprintf("The asset '%v' (ID: %v) has just been updated by %v", asset.AssetName, asset.Id, byUser.Email)
	userNotificationUnique := utils.ConvertUsersToNotificationsToMap(userId, usersToNotifications)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("SendNotificationToUsers panic:", r)
			}
		}()
		service.NotificationService.SendNotificationToUsers(userNotificationUnique, message, *asset)
	}()
	return assignmentUpdated, nil
}

func (service *AssignmentService) Filter(userId int64, emailAssigned *string, emailAssign *string, assetName *string) ([]dto.AssignmentResponse, error) {
	var filter = filter.AssignmentFilter{
		EmailAssigned: emailAssigned,
		EmailAssign:   emailAssign,
		AssetName:     assetName,
	}
	db := service.repo.GetDB()
	dbFilter := filter.ApplyFilter(db.Model(&entity.Assignments{}), userId)

	var total int64
	dbFilter.Count(&total)
	var assignments []entity.Assignments
	result := dbFilter.Find(&assignments)
	if result.Error != nil {
		return nil, result.Error
	}
	assignmentsRes := utils.ConvertAssignmentsToResponses(assignments)
	return assignmentsRes, nil
}

func (service *AssignmentService) GetAssignmentById(userId int64, id int64) (*dto.AssignmentResponse, error) {
	assignment, err := service.repo.GetAssignmentById(id)
	if err != nil {
		return nil, err
	}
	assignResponse := utils.ConvertAssignmentToResponse(assignment)
	return &assignResponse, nil
}

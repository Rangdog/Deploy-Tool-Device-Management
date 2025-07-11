package service

import (
	"BE_Manage_device/internal/domain/entity"
	asset "BE_Manage_device/internal/repository/assets"
	maintenanceSchedules "BE_Manage_device/internal/repository/maintenance_schedules"
	user "BE_Manage_device/internal/repository/user"
	notificationS "BE_Manage_device/internal/service/notification"
	"errors"
	"fmt"
	"time"
)

type MaintenanceSchedulesService struct {
	repo                maintenanceSchedules.MaintenanceSchedulesRepository
	assetRepo           asset.AssetsRepository
	userRepository      user.UserRepository
	NotificationService *notificationS.NotificationService
}

func NewMaintenanceSchedulesService(repo maintenanceSchedules.MaintenanceSchedulesRepository, assetRepo asset.AssetsRepository, userRepository user.UserRepository, NotificationService *notificationS.NotificationService) *MaintenanceSchedulesService {
	return &MaintenanceSchedulesService{repo: repo, assetRepo: assetRepo, NotificationService: NotificationService, userRepository: userRepository}
}

func (service *MaintenanceSchedulesService) Create(userId int64, assetId int64, startDate, endDate time.Time) (*entity.MaintenanceSchedules, error) {
	loc, _ := time.LoadLocation("Asia/Bangkok") // GMT+7
	startDate = startDate.In(loc)
	endDate = endDate.In(loc)
	userUpdate, err := service.userRepository.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	assetCheck, err := service.assetRepo.GetAssetById(assetId)
	if err != nil {
		return nil, err
	}
	if assetCheck.Status == "Disposed" || assetCheck.Status == "Retired" || assetCheck.Status == "Under Maintenance" {
		return nil, errors.New("can't set maintenance schedules because status")
	}
	maintenance := entity.MaintenanceSchedules{
		AssetId:   assetId,
		StartDate: startDate,
		EndDate:   endDate,
	}
	timeRange, err := service.repo.GetDateMaintenanceSchedulesInFuture(assetId)
	if err != nil {
		return nil, err
	}
	for _, r := range timeRange {
		if !(endDate.Before(r.Start) || startDate.After(r.End)) {
			return nil, errors.New("maintenance time overlaps with existing schedule")
		}
	}
	maintenanceCreate, err := service.repo.Create(&maintenance)
	if err != nil {
		return nil, err
	}
	userManagerAsset, _ := service.userRepository.GetUserAssetManageOfDepartment(assetCheck.DepartmentId)
	usersToNotifications := []*entity.Users{}
	usersToNotifications = append(usersToNotifications, assetCheck.OnwerUser)
	usersToNotifications = append(usersToNotifications, userManagerAsset)
	filteredUsers := []*entity.Users{}
	for _, user := range usersToNotifications {
		if user.Id != userUpdate.Id {
			filteredUsers = append(filteredUsers, user)
		}
	}
	usersToNotifications = filteredUsers
	message := fmt.Sprintf("The maintenance schedules (ID: %v) has just been created by %v", maintenance.Id, userUpdate.Email)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("SendNotificationToUsers panic:", r)
			}
		}()
		service.NotificationService.SendNotificationToUsers(usersToNotifications, message, *assetCheck)
	}()
	return maintenanceCreate, nil
}

func (service *MaintenanceSchedulesService) GetAllMaintenanceSchedulesByAssetId(userId int64, assetId int64) ([]*entity.MaintenanceSchedules, error) {
	maintenances, err := service.repo.GetAllMaintenanceSchedulesByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	return maintenances, nil
}

func (service *MaintenanceSchedulesService) Update(userId int64, id int64, startDate time.Time, endDate time.Time) (*entity.MaintenanceSchedules, error) {
	loc, _ := time.LoadLocation("Asia/Bangkok") // GMT+7
	startDate = startDate.In(loc)
	endDate = endDate.In(loc)
	userUpdate, err := service.userRepository.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	maintenaceUpdateOld, _ := service.repo.GetMaintenanceSchedulesById(id)
	if maintenaceUpdateOld.StartDate.Before(time.Now()) {
		return nil, errors.New("start date <= now")
	}
	maintenance, err := service.repo.Update(id, startDate, endDate)
	if err != nil {
		return nil, err
	}
	maintenaceUpdate, _ := service.repo.GetMaintenanceSchedulesById(id)

	userManagerAsset, _ := service.userRepository.GetUserAssetManageOfDepartment(maintenaceUpdate.Asset.DepartmentId)
	usersToNotifications := []*entity.Users{}
	usersToNotifications = append(usersToNotifications, maintenaceUpdate.Asset.OnwerUser)
	usersToNotifications = append(usersToNotifications, userManagerAsset)
	filteredUsers := []*entity.Users{}
	for _, user := range usersToNotifications {
		if user.Id != userUpdate.Id {
			filteredUsers = append(filteredUsers, user)
		}
	}
	usersToNotifications = filteredUsers
	message := fmt.Sprintf("The maintenance schedules (ID: %v) has just been updated by %v", maintenaceUpdate.Id, userUpdate.Email)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("SendNotificationToUsers panic:", r)
			}
		}()
		service.NotificationService.SendNotificationToUsers(usersToNotifications, message, maintenaceUpdate.Asset)
	}()
	return maintenance, nil
}

func (service *MaintenanceSchedulesService) Delete(userId int64, id int64) error {
	userUpdate, err := service.userRepository.FindByUserId(userId)
	if err != nil {
		return err
	}
	maintenanceCheck, err := service.repo.GetMaintenanceSchedulesById(id)
	if err != nil {
		return err
	}
	if maintenanceCheck.StartDate.After(time.Now()) {
		return errors.New("start date <= now")
	}
	err = service.repo.Delete(id)
	maintenaceUpdate, _ := service.repo.GetMaintenanceSchedulesById(id)

	userManagerAsset, _ := service.userRepository.GetUserAssetManageOfDepartment(maintenaceUpdate.Asset.DepartmentId)
	usersToNotifications := []*entity.Users{}
	usersToNotifications = append(usersToNotifications, maintenaceUpdate.Asset.OnwerUser)
	usersToNotifications = append(usersToNotifications, userManagerAsset)
	filteredUsers := []*entity.Users{}
	for _, user := range usersToNotifications {
		if user.Id != userUpdate.Id {
			filteredUsers = append(filteredUsers, user)
		}
	}
	usersToNotifications = filteredUsers
	message := fmt.Sprintf("The maintenance schedules (ID: %v) has just been deleted by %v", maintenaceUpdate.Id, userUpdate.Email)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("SendNotificationToUsers panic:", r)
			}
		}()
		service.NotificationService.SendNotificationToUsers(usersToNotifications, message, maintenaceUpdate.Asset)
	}()
	return err
}

func (service *MaintenanceSchedulesService) GetAllMaintenanceSchedules() ([]*entity.MaintenanceSchedules, error) {
	maintenances, err := service.repo.GetAllMaintenanceSchedules()
	if err != nil {
		return nil, err
	}
	return maintenances, nil
}

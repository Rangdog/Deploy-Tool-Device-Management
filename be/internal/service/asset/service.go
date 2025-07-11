package service

import (
	"BE_Manage_device/config"
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/filter"
	asset_log "BE_Manage_device/internal/repository/asset_log"
	asset "BE_Manage_device/internal/repository/assets"
	assignment "BE_Manage_device/internal/repository/assignments"
	company "BE_Manage_device/internal/repository/company"
	department "BE_Manage_device/internal/repository/departments"
	role "BE_Manage_device/internal/repository/role"
	user "BE_Manage_device/internal/repository/user"
	userRBAC "BE_Manage_device/internal/repository/user_rbac"
	notificationS "BE_Manage_device/internal/service/notification"
	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
	"encoding/json"

	"context"
	"sync"

	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type AssetsService struct {
	repo                 asset.AssetsRepository
	assertLogRepository  asset_log.AssetsLogRepository
	roleRepository       role.RoleRepository
	userRBACRepository   userRBAC.UserRBACRepository
	userRepository       user.UserRepository
	assignRepository     assignment.AssignmentRepository
	departmentRepository department.DepartmentsRepository
	NotificationService  *notificationS.NotificationService
	companyRepo          company.CompanyRepository
}

func NewAssetsService(repo asset.AssetsRepository, assertLogRepository asset_log.AssetsLogRepository, roleRepository role.RoleRepository, userRBACRepository userRBAC.UserRBACRepository, userRepository user.UserRepository, assignRepository assignment.AssignmentRepository, departmentRepository department.DepartmentsRepository, NotificationService *notificationS.NotificationService, companyRepo company.CompanyRepository) *AssetsService {
	return &AssetsService{repo: repo, assertLogRepository: assertLogRepository, roleRepository: roleRepository, userRBACRepository: userRBACRepository, userRepository: userRepository, assignRepository: assignRepository, departmentRepository: departmentRepository, NotificationService: NotificationService, companyRepo: companyRepo}
}

func (service *AssetsService) Create(userId int64, assetName string, purchaseDate time.Time, warrantExpiry time.Time, serialNumber string, image *multipart.FileHeader, fileAttachment *multipart.FileHeader, categoryId int64, departmentId int64, url string, cost float64) (*entity.Assets, error) {
	imgFile, err := image.Open()
	if err != nil {
		return nil, fmt.Errorf("cannot open image: %w", err)
	}
	defer imgFile.Close()
	fileFile, err := fileAttachment.Open()
	if err != nil {
		return nil, fmt.Errorf("cannot open fileAttachment: %w", err)
	}
	defer fileFile.Close()
	uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), image.Filename)
	imagePath := "images/" + uniqueName
	uniqueName = fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileAttachment.Filename)
	filePath := "files/" + uniqueName
	uploader := utils.NewSupabaseUploader()
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, 2)
	)
	var imageUrl string
	var fileUrl string
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Add(2)
	go func() {
		defer wg.Done()
		i, err := uploader.Upload(imagePath, imgFile, image.Header.Get("Content-Type"))
		if err != nil {
			errChan <- err
			cancel()
			return
		}
		imageUrl = i
	}()
	go func() {
		defer wg.Done()
		f, err := uploader.Upload(filePath, fileFile, fileAttachment.Header.Get("Content-Type"))
		if err != nil {
			errChan <- err
			cancel()
			return
		}
		fileUrl = f
	}()
	wg.Wait()
	close(errChan)
	if err, ok := <-errChan; ok {
		return nil, err
	}
	wg = sync.WaitGroup{}
	errChan = make(chan error, 2)
	var userAssetManager *entity.Users
	var company *entity.Company
	wg.Add(1)
	go func() {
		defer wg.Done()
		uam, err := service.userRepository.GetUserAssetManageOfDepartment(departmentId)
		if err != nil {
			errChan <- err
			cancel()
			return
		}
		userAssetManager = uam
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		u, err := service.userRepository.FindByUserId(userId)
		if err != nil {
			errChan <- err
			cancel()
			return
		}
		c, err := service.companyRepo.GetCompanyBySuffixEmail(utils.GetSuffixEmail(u.Email))
		if err != nil {
			errChan <- err
			cancel()
			return
		}
		company = c
	}()
	wg.Wait()
	close(errChan)
	if err, ok := <-errChan; ok {
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
	asset := entity.Assets{
		AssetName:      assetName,
		PurchaseDate:   purchaseDate,
		Cost:           cost,
		WarrantExpiry:  warrantExpiry,
		Status:         "New",
		SerialNumber:   serialNumber,
		ImageUpload:    &imageUrl,
		FileAttachment: &fileUrl,
		CategoryId:     categoryId,
		DepartmentId:   departmentId,
		Owner:          &userAssetManager.Id,
		CompanyId:      company.Id,
	}
	assetCreate, err := service.repo.Create(&asset, tx)
	if err != nil {
		return nil, err
	}
	changeSummary := "Create asset"
	assetLog := entity.AssetLog{
		Action:        "Create",
		Timestamp:     time.Now(),
		ByUserId:      &userId,
		AssignUserId:  &userAssetManager.Id,
		ChangeSummary: changeSummary,
		AssetId:       assetCreate.Id,
		CompanyId:     company.Id,
	}
	_, err = service.assertLogRepository.Create(&assetLog, tx)
	if err != nil {
		return nil, err
	}
	assign := entity.Assignments{
		AssetId:      assetCreate.Id,
		UserId:       &userAssetManager.Id,
		AssignBy:     userId,
		DepartmentId: &departmentId,
		CompanyId:    company.Id,
	}
	_, err = service.assignRepository.Create(&assign, tx)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	go service.SetRole(assetCreate.Id)
	go utils.GenQrAndUpdate(service.repo, assetCreate.Id, url)
	return assetCreate, nil
}

func (service *AssetsService) GetAssetById(userId int64, assertId int64) (*entity.Assets, error) {
	assert, err := service.repo.GetAssetById(assertId)
	if err != nil {
		return nil, err
	}
	return assert, err
}

func (service *AssetsService) GetAllAsset(userId int64, cacheKey string) ([]*entity.Assets, error) {
	user, err := service.userRepository.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	cacheKeyAssetCompanyId := fmt.Sprintf("%v:%v", cacheKey, user.CompanyId)
	val, err := config.Rdb.Get(config.Ctx, cacheKeyAssetCompanyId).Result()
	var assets []*entity.Assets
	if err == nil {
		var cached []entity.Assets
		if err := json.Unmarshal([]byte(val), &cached); err == nil {
			for _, a := range cached {
				copy := a
				assets = append(assets, &copy)
			}
		} else {
			log.Error("Happened error when get all asset. Error", err)
			pkg.PanicExeption(constant.UnknownError, "Happened error when get all asset in redis")
		}
	} else {
		assets, err = service.repo.GetAllAsset(user.CompanyId)
		if err != nil {
			log.Error("Happened error when get all asset. Error", err)
			pkg.PanicExeption(constant.UnknownError, "Happened error when get all asset")
		}
	}
	// ✅ Cache lại dữ liệu
	bytes, _ := json.Marshal(assets)
	config.Rdb.Set(config.Ctx, cacheKeyAssetCompanyId, bytes, 10*time.Minute)
	return assets, err
}

func (service *AssetsService) SetRole(assetId int64) error {
	db := service.repo.GetDB()
	assets, err := service.repo.GetAssetById(assetId)
	if err != nil {
		return err
	}
	users := service.userRepository.GetAllUser(assets.CompanyId)
	for _, user := range users {
		userRbac := entity.UserRbac{
			AssetId: assetId,
			UserId:  user.Id,
			RoleId:  user.RoleId,
		}
		err := service.userRBACRepository.Create(&userRbac, db)
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *AssetsService) UpdateAsset(userId int64, assetId int64, assetName string, purchaseDate time.Time, warrantExpiry time.Time, serialNumber string, image *multipart.FileHeader, fileAttachment *multipart.FileHeader, categoryId int64, cost float64) (*entity.Assets, error) {
	var err error
	var imgFile multipart.File
	imgFile, err = image.Open()
	if err != nil {
		return nil, fmt.Errorf("cannot open image: %w", err)
	}
	defer imgFile.Close()
	fileFile, err := fileAttachment.Open()
	if err != nil {
		return nil, fmt.Errorf("cannot open fileAttachment: %w", err)
	}
	defer fileFile.Close()
	uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), image.Filename)
	imagePath := "images/" + uniqueName
	uniqueName = fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileAttachment.Filename)
	filePath := "files/" + uniqueName
	oldAsset, err := service.repo.GetAssetById(assetId)
	if err != nil {
		return nil, fmt.Errorf("cannot find asset: %w", err)
	}
	var filedUpdate []string
	uploader := utils.NewSupabaseUploader()
	if oldAsset.ImageUpload != nil && *oldAsset.ImageUpload != "" {
		oldImagePath, _ := utils.ExtractFilePath(*oldAsset.ImageUpload)
		_ = uploader.Delete(oldImagePath)
		filedUpdate = append(filedUpdate, "image")
	}
	if oldAsset.FileAttachment != nil && *oldAsset.FileAttachment != "" {
		oldFilePath, _ := utils.ExtractFilePath(*oldAsset.FileAttachment)
		_ = uploader.Delete(oldFilePath)
		filedUpdate = append(filedUpdate, "file")
	}
	var (
		wg      sync.WaitGroup
		errChan = make(chan error, 2)
	)
	var imageUrl string
	var fileUrl string
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Add(2)
	go func() {
		defer wg.Done()
		i, err := uploader.Upload(imagePath, imgFile, image.Header.Get("Content-Type"))
		if err != nil {
			errChan <- err
			cancel()
			return
		}
		imageUrl = i
	}()
	go func() {
		defer wg.Done()
		f, err := uploader.Upload(filePath, fileFile, fileAttachment.Header.Get("Content-Type"))
		if err != nil {
			errChan <- err
			cancel()
			return
		}
		fileUrl = f
	}()
	wg.Wait()
	close(errChan)
	if err, ok := <-errChan; ok {
		return nil, err
	}
	if oldAsset.AssetName != assetName {
		filedUpdate = append(filedUpdate, "asset name")
	}
	if oldAsset.PurchaseDate != purchaseDate {
		filedUpdate = append(filedUpdate, "purchase date")
	}
	if oldAsset.Cost != cost {
		filedUpdate = append(filedUpdate, "cost")
	}
	if oldAsset.WarrantExpiry != warrantExpiry {
		filedUpdate = append(filedUpdate, "warrant expiry")
	}
	if oldAsset.SerialNumber != serialNumber {
		filedUpdate = append(filedUpdate, "serial number")
	}
	if oldAsset.CategoryId != categoryId {
		filedUpdate = append(filedUpdate, "category")
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
	asset := entity.Assets{
		Id:             assetId,
		AssetName:      assetName,
		PurchaseDate:   purchaseDate,
		Cost:           cost,
		WarrantExpiry:  warrantExpiry,
		SerialNumber:   serialNumber,
		ImageUpload:    &imageUrl,
		FileAttachment: &fileUrl,
		CategoryId:     categoryId,
	}
	assetUpdated, err := service.repo.UpdateAsset(&asset, tx)
	if err != nil {
		return nil, err
	}

	userUpdate, err := service.userRepository.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	changeSummary := ""
	if len(filedUpdate) > 0 {
		changeSummary = "Update fields: " + strings.Join(filedUpdate, ", ")
	} else {
		changeSummary = "no changes"
	}
	assetLog := entity.AssetLog{
		Action:        "Update",
		Timestamp:     time.Now(),
		ByUserId:      &userId,
		ChangeSummary: changeSummary,
		AssetId:       assetId,
		CompanyId:     userUpdate.CompanyId,
	}
	_, err = service.assertLogRepository.Create(&assetLog, tx)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}
	userManagerAsset, _ := service.userRepository.GetUserAssetManageOfDepartment(assetUpdated.DepartmentId)
	usersToNotifications := []*entity.Users{}
	if assetUpdated.OnwerUser != nil {
		usersToNotifications = append(usersToNotifications, assetUpdated.OnwerUser)
	}
	if userManagerAsset != nil {
		usersToNotifications = append(usersToNotifications, userManagerAsset)
	}
	message := fmt.Sprintf("The asset '%v' (ID: %v) has just been updated by %v", assetName, assetId, userUpdate.Email)
	userNotificationUnique := utils.ConvertUsersToNotificationsToMap(userId, usersToNotifications)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("SendNotificationToUsers panic:", r)
			}
		}()
		service.NotificationService.SendNotificationToUsers(userNotificationUnique, message, asset)
	}()
	return assetUpdated, nil
}

func (service *AssetsService) DeleteAsset(userId int64, id int64) error {
	var err error
	userUpdate, err := service.userRepository.FindByUserId(userId)
	if err != nil {
		return err
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
	err = service.repo.DeleteAsset(id, tx)
	if err != nil {
		return err
	}
	asset, err := service.repo.GetAssetById(id)
	if err != nil {
		return err
	}
	changeSummary := "Delete asset chair"
	assetLog := entity.AssetLog{
		Action:        "Delete",
		Timestamp:     time.Now(),
		ByUserId:      &userId,
		ChangeSummary: changeSummary,
		AssetId:       asset.Id,
		CompanyId:     userUpdate.CompanyId,
	}
	_, err = service.assertLogRepository.Create(&assetLog, tx)
	if err != nil {
		return err
	}
	if err = tx.Commit().Error; err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}
	userManagerAsset, _ := service.userRepository.GetUserAssetManageOfDepartment(asset.DepartmentId)
	usersToNotifications := []*entity.Users{}
	usersToNotifications = append(usersToNotifications, asset.OnwerUser)
	usersToNotifications = append(usersToNotifications, userManagerAsset)
	message := fmt.Sprintf("The asset '%v' (ID: %v) has just been delete by %v", asset.AssetName, asset.Id, userUpdate.Email)
	userNotificationUnique := utils.ConvertUsersToNotificationsToMap(userId, usersToNotifications)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("SendNotificationToUsers panic:", r)
			}
		}()
		service.NotificationService.SendNotificationToUsers(userNotificationUnique, message, *asset)
	}()
	return nil
}

func (service *AssetsService) UpdateAssetRetired(userId int64, id int64, ResidualValue float64) (*entity.Assets, error) {
	var err error
	userUpdate, err := service.userRepository.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	assetCheck, err := service.repo.GetAssetById(id)
	if err != nil {
		return nil, err
	}
	if assetCheck.Status == "Disposed" {
		return nil, errors.New("can't retired asset")
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
	asset, err := service.repo.UpdateAssetLifeCycleStage(id, "Retired", tx)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	asset.AcquisitionDate = &now
	asset.ResidualValue = &ResidualValue
	yearsUsed := now.Year() - asset.AcquisitionDate.Year()
	if now.YearDay() < asset.AcquisitionDate.YearDay() {
		yearsUsed--
	}
	if yearsUsed <= 0 {
		fmt.Println("Asset not in use for a full year.")
		return nil, errors.New("a")
	}
	// Tính khấu hao hàng năm
	annualDepreciation := (asset.Cost - *asset.ResidualValue) / float64(yearsUsed)
	asset.AnnualDepreciation = &annualDepreciation
	asset, err = service.repo.UpdateAsset(asset, tx)
	if err != nil {
		return nil, err
	}
	changeSummary := "Retired asset"
	assetLog := entity.AssetLog{
		Action:        "Update",
		Timestamp:     time.Now(),
		ByUserId:      &userId,
		ChangeSummary: changeSummary,
		AssetId:       asset.Id,
		CompanyId:     userUpdate.CompanyId,
	}
	_, err = service.assertLogRepository.Create(&assetLog, tx)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit failed: %w", err)
	}
	userManagerAsset, _ := service.userRepository.GetUserAssetManageOfDepartment(asset.DepartmentId)
	usersToNotifications := []*entity.Users{}
	usersToNotifications = append(usersToNotifications, asset.OnwerUser)
	usersToNotifications = append(usersToNotifications, userManagerAsset)
	message := fmt.Sprintf("The asset '%v' (ID: %v) has just been updated by %v", asset.AssetName, asset.Id, userUpdate.Email)
	userNotificationUnique := utils.ConvertUsersToNotificationsToMap(userId, usersToNotifications)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("SendNotificationToUsers panic:", r)
			}
		}()
		service.NotificationService.SendNotificationToUsers(userNotificationUnique, message, *asset)
	}()
	return asset, nil

}

func (service *AssetsService) Filter(userId int64, assetName *string, status *string, categoryId *string, cost *string, serialNumber *string, email *string, departmentId *string) ([]dto.AssetResponse, error) {
	var filter = filter.AssetFilter{
		AssetName:    assetName,
		CategoryId:   categoryId,
		Cost:         cost,
		SerialNumber: serialNumber,
		Email:        email,
		DepartmentId: departmentId,
		Status:       status,
	}
	user, err := service.userRepository.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	filter.CompanyId = user.CompanyId
	if user.Role.Slug == "departmentHead" {
		var deptStr *string
		if user.DepartmentId != nil {
			s := strconv.FormatInt(*user.DepartmentId, 10)
			deptStr = &s
		} else {
			deptStr = nil
		}
		filter.DepartmentId = deptStr
	}
	db := service.repo.GetDB()
	dbFilter := filter.ApplyFilter(db.Model(&entity.Assets{}), userId)
	var total int64
	dbFilter.Count(&total)
	var assets []entity.Assets
	result := dbFilter.Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	assetsResponse := []dto.AssetResponse{}
	for _, asset := range assets {
		assetResponse := dto.AssetResponse{
			ID:             asset.Id,
			AssetName:      asset.AssetName,
			PurchaseDate:   asset.PurchaseDate.Format("2006-01-02"),
			Cost:           asset.Cost,
			WarrantExpiry:  asset.WarrantExpiry.Format("2006-01-02"),
			Status:         asset.Status,
			SerialNumber:   asset.SerialNumber,
			FileAttachment: *asset.FileAttachment,
			ImageUpload:    *asset.ImageUpload,
			QrURL:          *asset.QrUrl,
			Category: dto.CategoryResponse{
				ID:           asset.Category.Id,
				CategoryName: asset.Category.CategoryName,
			},
			Department: dto.DepartmentResponse{
				ID:             asset.Department.Id,
				DepartmentName: asset.Department.DepartmentName,
				Location: dto.LocationResponse{
					ID:           asset.Department.Location.Id,
					LocationName: asset.Department.Location.LocationName,
				},
			},
		}
		if asset.OnwerUser != nil {
			assetResponse.Owner = dto.OwnerResponse{
				ID:        asset.OnwerUser.Id,
				FirstName: asset.OnwerUser.FirstName,
				LastName:  asset.OnwerUser.LastName,
				Email:     asset.OnwerUser.Email,
			}
		}
		assetsResponse = append(assetsResponse, assetResponse)
	}
	return assetsResponse, nil
}

func (service *AssetsService) ApplyFilterDashBoard(userId int64, status *string, categoryId *string, departmentId *string, export *string) (*dto.DashboardSummary, []*entity.Assets, error) {
	user, err := service.userRepository.FindByUserId(userId)
	if err != nil {
		return nil, nil, err
	}
	roleHeadDep := service.roleRepository.GetRoleBySlug("department-head")
	roleView := service.roleRepository.GetRoleBySlug("employee")
	if user.RoleId == roleHeadDep.Id || user.RoleId == roleView.Id {
		if user.DepartmentId != nil {
			// Convert *int64 to string
			depIdStr := strconv.FormatInt(*user.DepartmentId, 10)
			departmentId = &depIdStr // assign pointer to string
		} else {
			departmentId = nil
		}
	}

	var filter = filter.AssetFilterDashboard{
		CategoryId:   categoryId,
		DepartmentId: departmentId,
		Status:       status,
	}
	db := service.repo.GetDB()
	dbFilter := filter.ApplyFilterDashBoard(db.Model(&entity.Assets{}), userId)
	var assets []*entity.Assets
	result := dbFilter.Find(&assets)
	if result.Error != nil {
		return nil, nil, result.Error
	}
	summary := CountDashboard(assets)
	return &summary, assets, nil
}

func CountDashboard(assets []*entity.Assets) dto.DashboardSummary {
	var s dto.DashboardSummary
	s.TotalAssets = len(assets)
	for _, a := range assets {
		switch a.Status {
		case "In Use":
			s.Assigned++
		case "Under Maintenance":
			s.UnderMaintenance++
		case "Retired":
			s.Retired++
		}
	}
	return s
}

func (service *AssetsService) GetAssetsByCateOfDepartment(userId, categoryId int64) ([]*entity.Assets, error) {
	user, err := service.userRepository.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	assets, err := service.repo.GetAssetsByCateOfDepartment(categoryId, *user.DepartmentId)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (service *AssetsService) CheckPermissionForManager(userId int64, depId int64) error {
	user, err := service.userRepository.FindByUserId(userId)
	role := service.roleRepository.GetRoleBySlug("admin")
	if user.RoleId == role.Id {
		return nil
	}
	if err != nil {
		return err
	}
	if user.DepartmentId != nil && *user.DepartmentId == depId {
		return nil
	}
	return errors.New("you are not allowed to manage departmental assets")
}

func (service *AssetsService) GetUserById(id int64) (entity.Users, error) {
	user, err := service.userRepository.FindByUserId(id)
	return *user, err
}

func (service *AssetsService) GetAllAssetNotHaveMaintenance(userId int64) ([]*entity.Assets, error) {
	user, err := service.userRepository.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	assets, err := service.repo.GetAllAssetNotHaveMaintenance(user.CompanyId)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

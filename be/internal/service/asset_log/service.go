package service

import (
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/filter"
	asset_log "BE_Manage_device/internal/repository/asset_log"
	asset "BE_Manage_device/internal/repository/assets"
	role "BE_Manage_device/internal/repository/role"
	user "BE_Manage_device/internal/repository/user"

	"errors"
)

type AssetLogService struct {
	repo      asset_log.AssetsLogRepository
	userRepo  user.UserRepository
	roleRepo  role.RoleRepository
	assetRepo asset.AssetsRepository
}

func NewAssetLogService(repo asset_log.AssetsLogRepository, userRepo user.UserRepository, roleRepo role.RoleRepository, assetRepo asset.AssetsRepository) *AssetLogService {
	return &AssetLogService{repo: repo, userRepo: userRepo, roleRepo: roleRepo, assetRepo: assetRepo}
}

func (service *AssetLogService) GetLogByAssetId(assetId int64) ([]*entity.AssetLog, error) {
	assetlogs, err := service.repo.GetLogByAssetId(assetId)
	if err != nil {
		return nil, err
	}
	return assetlogs, nil
}
func (service *AssetLogService) Filter(userId int64, assetId int64, action, startTime, endTime *string) ([]dto.AssetLogsResponse, error) {
	userCheck, err := service.userRepo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	var filter = filter.AssetLogFilter{
		Action:    action,
		StartTime: startTime,
		EndTime:   endTime,
	}
	if userCheck.Role.Slug == "assetManager" {
		if userCheck.DepartmentId != nil {
			filter.DepId = userCheck.DepartmentId
		}
	}

	db := service.repo.GetDB()
	dbFilter := filter.ApplyFilter(db.Model(&entity.AssetLog{}), assetId)
	var total int64
	dbFilter.Count(&total)
	var asset_logs []entity.AssetLog
	result := dbFilter.Find(&asset_logs)
	if result.Error != nil {
		return nil, result.Error
	}
	assetLogResponses := []dto.AssetLogsResponse{}
	for _, assetLog := range asset_logs {
		var assetLogResponse dto.AssetLogsResponse

		assetLogResponse.Action = assetLog.Action
		assetLogResponse.Timestamp = assetLog.Timestamp.Format("2006-01-02")
		assetLogResponse.ChangeSummary = assetLog.ChangeSummary

		// byUser
		if assetLog.ByUserId != nil {
			assetLogResponse.ByUser.Id = *assetLog.ByUserId
		}
		if assetLog.ByUser != nil {
			assetLogResponse.ByUser.FirstName = assetLog.ByUser.FirstName
			assetLogResponse.ByUser.LastName = assetLog.ByUser.LastName
			assetLogResponse.ByUser.Email = assetLog.ByUser.Email
		}

		// assignUser
		if assetLog.AssignUser != nil {
			assetLogResponse.AssignUser = &dto.UserResponseInAssetLog{
				Id:        assetLog.AssignUser.Id,
				FirstName: assetLog.AssignUser.FirstName,
				LastName:  assetLog.AssignUser.LastName,
				Email:     assetLog.AssignUser.Email,
			}
		}

		// asset

		assetLogResponse.Asset.AssetName = assetLog.Asset.AssetName
		assetLogResponse.Asset.Id = assetLog.Asset.Id
		assetLogResponse.Asset.SerialNumber = assetLog.Asset.SerialNumber

		if assetLog.Asset.ImageUpload != nil {
			assetLogResponse.Asset.ImageUpload = *assetLog.Asset.ImageUpload
		}
		if assetLog.Asset.FileAttachment != nil {
			assetLogResponse.Asset.FileAttachment = *assetLog.Asset.FileAttachment
		}
		if assetLog.Asset.QrUrl != nil {
			assetLogResponse.Asset.QrUrl = *assetLog.Asset.QrUrl
		}

		assetLogResponses = append(assetLogResponses, assetLogResponse)
	}
	return assetLogResponses, nil
}

func (service *AssetLogService) CheckPermissionForManager(userId int64, depId int64) error {
	user, err := service.userRepo.FindByUserId(userId)
	role := service.roleRepo.GetRoleBySlug("admin")
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

func (service *AssetLogService) GetAssetById(userId int64, assertId int64) (*entity.Assets, error) {
	assert, err := service.assetRepo.GetAssetById(assertId)
	if err != nil {
		return nil, err
	}
	return assert, err
}

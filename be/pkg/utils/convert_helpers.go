package utils

import (
	"BE_Manage_device/internal/domain/dto"
	"BE_Manage_device/internal/domain/entity"
)

func ConvertUserToUserResponse(user *entity.Users) dto.UserResponse {
	usersResponse := dto.UserResponse{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		IsActive:  user.IsActive,
		Avatar:    user.Avatar,
		Role: dto.UserRoleResponse{
			Id:   user.RoleId,
			Slug: user.Role.Slug,
		},
	}

	if user.DepartmentId != nil {
		usersResponse.Department = &dto.UserDepartmentResponse{
			Id:             *user.DepartmentId,
			DepartmentName: user.Department.DepartmentName,
		}
	}
	return usersResponse
}

func ConvertUsersToUserResponses(users []*entity.Users) []dto.UserResponse {
	responses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, ConvertUserToUserResponse(user))
	}
	return responses
}

func ConvertRequestTransferToResponse(rt *entity.RequestTransfer) dto.RequestTransferResponse {
	return dto.RequestTransferResponse{
		Id:          rt.Id,
		Status:      rt.Status,
		Description: rt.Description,
		User: dto.UserResponseInRequestTransfer{
			Id:           rt.User.Id,
			FirstName:    rt.User.FirstName,
			LastName:     rt.User.LastName,
			Email:        rt.User.Email,
			DepartmentId: derefInt64(rt.User.DepartmentId),
		},
		Category: dto.CategoryResponseInRequestTransfer{
			Id:           rt.CategoryId,
			CategoryName: rt.Category.CategoryName,
		},
	}
}

func ConvertRequestTransfersToResponses(rts []entity.RequestTransfer) []dto.RequestTransferResponse {
	res := make([]dto.RequestTransferResponse, 0, len(rts))
	for _, rt := range rts {
		res = append(res, ConvertRequestTransferToResponse(&rt))
	}
	return res
}

func ConvertAssignmentToResponse(assignment *entity.Assignments) dto.AssignmentResponse {
	return dto.AssignmentResponse{
		Id: assignment.Id,
		UserAssigned: dto.UsersAssignmentResponse{
			Id:        assignment.UserAssigned.Id,
			FirstName: assignment.UserAssigned.FirstName,
			LastName:  assignment.UserAssigned.LastName,
			Email:     assignment.UserAssigned.Email,
		},
		UserAssign: dto.UsersAssignmentResponse{
			Id:        assignment.UserAssign.Id,
			FirstName: assignment.UserAssign.FirstName,
			LastName:  assignment.UserAssign.LastName,
			Email:     assignment.UserAssign.Email,
		},
		Asset: dto.UserAssignmentAssetResponse{
			Id:             assignment.Asset.Id,
			AssetName:      assignment.Asset.AssetName,
			Status:         assignment.Asset.Status,
			FileAttachment: derefString(assignment.Asset.FileAttachment),
			ImageUpload:    derefString(assignment.Asset.ImageUpload),
		},
		Department: dto.DepartmentResponse{
			ID:             derefInt64(assignment.DepartmentId),
			DepartmentName: assignment.Department.DepartmentName,
			Location: dto.LocationResponse{
				ID:           assignment.Department.LocationId,
				LocationName: assignment.Department.Location.LocationName,
			},
		},
	}
}

// Helper for pointer value (nếu có thể nil)
func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
func derefInt64(i *int64) int64 {
	if i != nil {
		return *i
	}
	return 0
}

func ConvertAssignmentsToResponses(assignments []entity.Assignments) []dto.AssignmentResponse {
	res := make([]dto.AssignmentResponse, 0, len(assignments))
	for _, as := range assignments {
		res = append(res, ConvertAssignmentToResponse(&as))
	}
	return res
}

func ConvertMaintenanceSchedulesToResponses(maintenanceSchedules *entity.MaintenanceSchedules) dto.MaintenanceSchedulesResponse {
	return dto.MaintenanceSchedulesResponse{
		Id:        maintenanceSchedules.Id,
		StartDate: maintenanceSchedules.StartDate.Format("2006-01-02"),
		EndDate:   maintenanceSchedules.EndDate.Format("2006-01-02"),
		Asset: dto.AssetResponseInMaintenanceSchedules{
			Id:             maintenanceSchedules.AssetId,
			AssetName:      maintenanceSchedules.Asset.AssetName,
			Status:         maintenanceSchedules.Asset.Status,
			FileAttachment: derefString(maintenanceSchedules.Asset.FileAttachment),
			ImageUpload:    derefString(maintenanceSchedules.Asset.ImageUpload),
		},
	}
}

func ConvertMaintenanceSchedulesToResponsesArray(maintenanceSchedules []*entity.MaintenanceSchedules) []dto.MaintenanceSchedulesResponse {
	res := make([]dto.MaintenanceSchedulesResponse, 0, len(maintenanceSchedules))
	for _, as := range maintenanceSchedules {
		res = append(res, ConvertMaintenanceSchedulesToResponses(as))
	}
	return res
}

func ConvertRoleToResponse(role *entity.Roles) dto.RoleResponse {
	return dto.RoleResponse{
		Id:          role.Id,
		Slug:        role.Slug,
		Description: role.Description,
	}
}

func ConvertRolesToResponsesArray(role []*entity.Roles) []dto.RoleResponse {
	res := make([]dto.RoleResponse, 0, len(role))
	for _, as := range role {
		res = append(res, ConvertRoleToResponse(as))
	}
	return res
}

func ConvertUsersToNotificationsToMap(userId int64, usersToNotifications []*entity.Users) []*entity.Users {
	uniqueMap := make(map[int64]bool)
	uniqueUsers := []*entity.Users{}

	for _, user := range usersToNotifications {
		if user == nil {
			continue
		}
		if user.Id == userId {
			continue
		}
		if !uniqueMap[user.Id] {
			uniqueMap[user.Id] = true
			uniqueUsers = append(uniqueUsers, user)
		}
	}
	return uniqueUsers
}

func ConvertAssetToResponse(asset entity.Assets) dto.AssetResponse {
	var qrURL string
	if asset.QrUrl != nil {
		qrURL = *asset.QrUrl
	}
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
		QrURL:          qrURL,
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
	return assetResponse
}

func ConvertBillToResponse(bill *entity.Bill) dto.BillResponse {
	return dto.BillResponse{
		BillNumber:  bill.BillNumber,
		Amount:      bill.Amount,
		Description: bill.Description,
		CreateAt:    bill.CreateAt,
		Asset:       ConvertAssetToResponse(bill.Asset),
		CreateBy:    ConvertUserToUserResponse(&bill.CreateBy),
	}
}

func ConvertBillsToResponsesArray(bills []*entity.Bill) []dto.BillResponse {
	res := make([]dto.BillResponse, 0, len(bills))
	for _, as := range bills {
		res = append(res, ConvertBillToResponse(as))
	}
	return res
}

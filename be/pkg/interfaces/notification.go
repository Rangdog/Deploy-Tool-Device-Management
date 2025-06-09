package interfaces

import "BE_Manage_device/internal/domain/entity"

type Notification interface {
	SendNotificationToUsers(users []*entity.Users, message string, asset entity.Assets) error
}

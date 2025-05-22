package repository

import "BE_Manage_device/internal/domain/entity"

type RequestTransferRepository interface {
	Create(*entity.RequestTransfer) (*entity.RequestTransfer, error)
}

package service

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
	"fmt"
	"sync"
	"time"
)

type NotificationService struct {
	clients                map[string][]chan string
	mu                     sync.RWMutex
	notificationRepository repository.NotificationRepository
}

func NewNotificationService(notificationRepository repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		clients: make(map[string][]chan string), notificationRepository: notificationRepository,
	}
}

func (ns *NotificationService) Register(userId string) chan string {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ch := make(chan string, 10)
	ns.clients[userId] = append(ns.clients[userId], ch)
	return ch
}

func (ns *NotificationService) UnRegister(userId string, ch chan string) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	chans := ns.clients[userId]
	for i, c := range chans {
		if c == ch {
			ns.clients[userId] = append(chans[:i], chans[i+1:]...)
			close(c)
			break
		}
	}
	// Xoá map khi không còn ai connect
	if len(ns.clients[userId]) == 0 {
		delete(ns.clients, userId)
	}
}

// Push notification cho 1 user
func (ns *NotificationService) Push(userId, message string) {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	for _, ch := range ns.clients[userId] {
		ch <- message
	}
}

func (ns *NotificationService) IsOnline(userId string) bool {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	chans, ok := ns.clients[userId]
	return ok && len(chans) > 0
}

func (service *NotificationService) SendNotificationToUsers(users []*entity.Users, message string, asset entity.Assets) error {
	status := "pending"
	typeNotify := "Info"
	timeNotify := time.Now()
	for _, u := range users {
		if u == nil {
			continue
		}
		notify := entity.Notifications{
			Content:    &message,
			Status:     &status,
			Type:       &typeNotify,
			UserId:     &u.Id,
			AssetId:    &asset.Id,
			NotifyDate: &timeNotify,
		}
		_, err := service.notificationRepository.Create(&notify)
		if err != nil {
			// log lỗi, tuỳ quyết định dừng hay tiếp tục
			fmt.Printf("Lỗi lưu notification cho user %v: %v\n", u.Id, err)
		}
		isOnline := service.IsOnline(fmt.Sprintf("%v", u.Id))
		if isOnline {
			service.Push(fmt.Sprintf("%v", u.Id), message)
		} else {
			fmt.Printf("User %v đang offline, chỉ lưu notification DB\n", u.Id)
		}
	}
	return nil
}

func (service *NotificationService) GetNotificationsByUserId(userId int64) ([]*entity.Notifications, error) {
	notifications, err := service.notificationRepository.GetNotificationsByUserId(userId)
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

func (service *NotificationService) UpdateStatusToSeen(id int64) error {
	err := service.notificationRepository.UpdateStatus(id)
	return err
}

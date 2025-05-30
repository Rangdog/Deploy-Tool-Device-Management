package service

import (
	"sync"
)

type NotificationService struct {
	clients map[string][]chan string
	mu      sync.RWMutex
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		clients: make(map[string][]chan string),
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

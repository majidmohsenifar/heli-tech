package notification

import (
	"context"
	"fmt"
	"sync"
)

type Service struct {
	allNotifications []SendNotificationParams
	mutex            sync.Mutex
}

type SendNotificationParams struct {
	TransactionID int64
	UserID        int64
	Amount        float64
	Balance       float64
	Kind          string
	CreatedAt     int64
}

func (s *Service) SendNotification(ctx context.Context, params SendNotificationParams) error {
	fmt.Println("got the notif", params)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.allNotifications) > 10 {
		s.allNotifications = []SendNotificationParams{}
	}
	s.allNotifications = append(s.allNotifications, params)
	return nil
}

// just for test purposes
func (s *Service) GetAllNotifications() []SendNotificationParams {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.allNotifications
}

func NewService() *Service {
	return &Service{}
}

package notification

import (
	"context"
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
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if len(s.allNotifications) > 10 {
		//empty slices without re allocation
		s.allNotifications = s.allNotifications[:0]
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

package notification

import "context"

type Service struct {
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
	return nil
}

func NewService() *Service {
	return &Service{}
}

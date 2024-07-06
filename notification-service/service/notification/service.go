package notification

import "context"

type Service struct {
}

type SendNotificationParams struct {
	UserID        int64
	TransactionID int64
	Amount        float64
}

func (s *Service) SendNotification(ctx context.Context, params SendNotificationParams) error {
	return nil
}

func NewService() *Service {
	return &Service{}
}

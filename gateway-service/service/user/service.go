package user

import (
	"context"
	"log/slog"

	userpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/user"
)

type Service struct {
	userClient userpb.UserClient
	logger     *slog.Logger
}

type RegisterParams struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confrimPassword" validate:"required"`
}

type RegisterResponse struct{}

type UserData struct {
	ID    int64
	Email string
}

func (s *Service) Register(
	ctx context.Context,
	params RegisterParams,
) error {
	_, err := s.userClient.Register(
		ctx,
		&userpb.RegisterRequest{
			Email:           params.Email,
			Password:        params.Password,
			ConfirmPassword: params.ConfirmPassword,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUserData(
	ctx context.Context,
	token string,
	path string,
) (UserData, error) {
	res, err := s.userClient.GetUserDataByToken(
		ctx,
		&userpb.GetUserDataByTokenRequest{
			Token: token,
			Path:  path,
		},
	)
	if err != nil {
		return UserData{}, err
	}
	return UserData{
		ID:    res.Id,
		Email: res.Email,
	}, nil
}

func NewService(
	userClient userpb.UserClient,
	logger *slog.Logger,
) *Service {
	return &Service{
		userClient: userClient,
		logger:     logger,
	}
}

package usergrpc

import (
	"context"

	userpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/user"
	"github.com/majidmohsenifar/heli-tech/user-service/service/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	userService *user.Service
	userpb.UnimplementedUserServer
}

func (s *server) Register(
	ctx context.Context,
	req *userpb.RegisterRequest,
) (resp *userpb.RegisterResponse, err error) {
	//TODO: validate email, password, and confirm password
	err = s.userService.Register(user.RegisterParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == user.ErrEmailAlreadyExist {
		return nil, status.Error(codes.Code(423), err.Error())
	}
	resp.Ok = true
	return resp, nil
}

func (s *server) Login(
	ctx context.Context,
	req *userpb.LoginRequest,
) (resp *userpb.LoginResponse, err error) {
	//TODO: validate password, and email
	token, err := s.userService.Login(user.LoginParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == user.ErrUserNotFound {
		return nil, status.Error(codes.Code(404), err.Error())
	}
	resp.Token = token
	return resp, nil
}

func (s *server) GetUserDataByToken(
	ctx context.Context,
	req *userpb.GetUserDataByTokenRequest,
) (resp *userpb.GetUserDataByTokenResponse, err error) {
	result, err := s.userService.GetUserDataByToken(user.GetUserDataByTokenParams{
		Token: req.Token,
		Path:  req.Token,
	})

	resp.Email = result.Email
	resp.Id = result.ID
}

func NewServer(
	userService *user.Service,
) userpb.UserServer {
	return &server{
		userService: userService,
	}
}

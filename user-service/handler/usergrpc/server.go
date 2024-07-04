package usergrpc

import (
	"context"
	"net/mail"
	"strings"

	userpb "github.com/majidmohsenifar/heli-tech/data-contracts/proto/user"
	"github.com/majidmohsenifar/heli-tech/user-service/service/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	authService *auth.Service
	userpb.UnimplementedUserServer
}

func (s *server) Register(
	ctx context.Context,
	req *userpb.RegisterRequest,
) (resp *userpb.RegisterResponse, err error) {
	if strings.EqualFold(strings.Trim(req.Email, " "), "") {
		return nil, status.Error(codes.Code(400), "email is empty")
	}
	_, err = mail.ParseAddress(req.Email)
	if err != nil {
		return nil, status.Error(codes.Code(400), "invalid email")
	}
	password := strings.Trim(req.Password, " ")
	confirmPassword := strings.Trim(req.ConfirmPassword, " ")
	if strings.EqualFold(password, "") {
		return nil, status.Error(codes.Code(400), "password is empty")
	}
	if strings.EqualFold(confirmPassword, "") {
		return nil, status.Error(codes.Code(400), "confirmPassword is empty")
	}
	if password != confirmPassword {
		return nil, status.Error(codes.Code(400), "password and confirmPassword are not the same")
	}

	err = s.authService.Register(ctx, auth.RegisterParams{
		Email:    req.Email,
		Password: password,
	})
	if err == auth.ErrEmailAlreadyExist {
		return nil, status.Error(codes.Code(423), err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Code(500), "something went wrong")
	}
	return &userpb.RegisterResponse{Ok: true}, nil
}

func (s *server) Login(
	ctx context.Context,
	req *userpb.LoginRequest,
) (resp *userpb.LoginResponse, err error) {
	if strings.EqualFold(strings.Trim(req.Email, " "), "") {
		return nil, status.Error(codes.Code(400), "email is empty")
	}
	_, err = mail.ParseAddress(req.Email)
	if err != nil {
		return nil, status.Error(codes.Code(400), "invalid email")
	}
	if strings.EqualFold(strings.Trim(req.Password, " "), "") {
		return nil, status.Error(codes.Code(400), "password is empty")
	}
	token, err := s.authService.Login(ctx, auth.LoginParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == auth.ErrInvalidUsernameOrPassword {
		return nil, status.Error(codes.Code(401), err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Code(500), "something went wrong")
	}
	return &userpb.LoginResponse{Token: token}, nil
}

func (s *server) GetUserDataByToken(
	ctx context.Context,
	req *userpb.GetUserDataByTokenRequest,
) (resp *userpb.GetUserDataByTokenResponse, err error) {
	if strings.EqualFold(strings.Trim(req.Token, " "), "") {
		return nil, status.Error(codes.Code(400), "token is empty")
	}
	if strings.EqualFold(strings.Trim(req.Path, " "), "") {
		return nil, status.Error(codes.Code(400), "path is empty")
	}
	result, err := s.authService.GetUserDataByToken(ctx, auth.GetUserDataByTokenParams{
		Token: req.Token,
		Path:  req.Path,
	})
	if err == auth.ErrAccessDenied {
		return nil, status.Error(codes.Code(403), err.Error())
	}
	if err != nil {
		return nil, status.Error(codes.Code(500), "something went wrong")
	}
	resp.Email = result.Email
	resp.Id = result.ID
	return resp, err
}

func NewServer(
	authService *auth.Service,
) userpb.UserServer {
	return &server{
		authService: authService,
	}
}

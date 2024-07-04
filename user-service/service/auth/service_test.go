package auth_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/majidmohsenifar/heli-tech/user-service/core"
	"github.com/majidmohsenifar/heli-tech/user-service/logger"
	"github.com/majidmohsenifar/heli-tech/user-service/mocks"
	"github.com/majidmohsenifar/heli-tech/user-service/repository"
	"github.com/majidmohsenifar/heli-tech/user-service/service/auth"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Register_UsernameAlreadyExist(t *testing.T) {
	assert := assert.New(t)
	repo := new(mocks.MockQuerier)
	repo.EXPECT().GetUserByEmail(mock.Anything, mock.Anything, "test@test.com").Once().Return(repository.User{}, nil)
	passwordEncoder := core.NewPasswordEncoder()
	logger := logger.NewLogger()
	authService := auth.NewService(
		nil,
		repo,
		passwordEncoder,
		nil,
		logger,
		nil,
	)
	err := authService.Register(context.Background(), auth.RegisterParams{
		Email:    "test@test.com",
		Password: "123456789",
	})
	assert.Equal(err, auth.ErrEmailAlreadyExist)
	repo.AssertExpectations(t)
}

func TestService_Register_CannotGeneratePasswordHash(t *testing.T) {
	assert := assert.New(t)
	repo := new(mocks.MockQuerier)
	repo.EXPECT().GetUserByEmail(mock.Anything, mock.Anything, "test@test.com").Once().Return(repository.User{}, pgx.ErrNoRows)
	passwordEncoder := core.NewPasswordEncoder()
	logger := logger.NewLogger()
	authService := auth.NewService(
		nil,
		repo,
		passwordEncoder,
		nil,
		logger,
		nil,
	)
	err := authService.Register(context.Background(), auth.RegisterParams{
		Email:    "test@test.com",
		Password: strings.Repeat("a", 73),
	})
	assert.Equal(err, errors.New("cannot hash the password"))
	repo.AssertExpectations(t)
}

func TestService_Register_DefaultRoleDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	repo := new(mocks.MockQuerier)
	repo.EXPECT().GetUserByEmail(mock.Anything, mock.Anything, "test@test.com").Once().Return(repository.User{}, pgx.ErrNoRows)
	repo.EXPECT().GetRoleByCode(mock.Anything, mock.Anything, auth.RoleEndUser).Once().Return(repository.Role{}, pgx.ErrNoRows)
	passwordEncoder := core.NewPasswordEncoder()
	roleRouteManager := auth.NewRoleRouteManager(nil, repo)
	logger := logger.NewLogger()
	authService := auth.NewService(
		nil,
		repo,
		passwordEncoder,
		nil,
		logger,
		roleRouteManager,
	)
	err := authService.Register(context.Background(), auth.RegisterParams{
		Email:    "test@test.com",
		Password: "123456789",
	})
	assert.Equal(err, errors.New("default role does not exist"))
	repo.AssertExpectations(t)
}

func TestService_Register_CannotCreateUser(t *testing.T) {
	assert := assert.New(t)
	dbMock, err := pgxmock.NewPool()
	assert.Nil(err)
	defer dbMock.Close()
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	passwordEncoder := core.NewPasswordEncoder()
	repo := new(mocks.MockQuerier)
	repo.EXPECT().GetUserByEmail(mock.Anything, mock.Anything, "test@test.com").Once().Return(repository.User{}, pgx.ErrNoRows)
	repo.EXPECT().GetRoleByCode(mock.Anything, mock.Anything, auth.RoleEndUser).Once().Return(repository.Role{ID: 1, Code: auth.RoleEndUser}, nil)
	repo.EXPECT().CreateUser(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.CreateUserParams)
			if p.Email != "test@test.com" {
				return false
			}
			if passwordEncoder.CompareHashAndPassword(p.Password, "123456789") != nil {
				return false
			}
			return true
		}),
	).Once().Return(repository.User{}, errors.New("db error"))
	roleRouteManager := auth.NewRoleRouteManager(nil, repo)
	logger := logger.NewLogger()
	authService := auth.NewService(
		dbMock,
		repo,
		passwordEncoder,
		nil,
		logger,
		roleRouteManager,
	)
	err = authService.Register(context.Background(), auth.RegisterParams{
		Email:    "test@test.com",
		Password: "123456789",
	})
	assert.Equal(err, errors.New("cannot create user"))
	repo.AssertExpectations(t)
}

func TestService_Register_CannotAddDefaultRole(t *testing.T) {
	assert := assert.New(t)
	dbMock, err := pgxmock.NewPool()
	assert.Nil(err)
	defer dbMock.Close()
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	passwordEncoder := core.NewPasswordEncoder()
	repo := new(mocks.MockQuerier)
	repo.EXPECT().GetUserByEmail(mock.Anything, mock.Anything, "test@test.com").Once().Return(repository.User{}, pgx.ErrNoRows)
	repo.EXPECT().GetRoleByCode(mock.Anything, mock.Anything, auth.RoleEndUser).Once().Return(repository.Role{ID: 1, Code: auth.RoleEndUser}, nil)
	repo.EXPECT().CreateUser(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.CreateUserParams)
			if p.Email != "test@test.com" {
				return false
			}
			if passwordEncoder.CompareHashAndPassword(p.Password, "123456789") != nil {
				return false
			}
			return true
		}),
	).Once().Return(repository.User{ID: 1}, nil)
	repo.EXPECT().AddRoleToUser(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.AddRoleToUserParams)
			if p.UserID != 1 {
				return false
			}
			if p.RoleID != 1 {
				return false
			}
			return true
		}),
	).Once().Return(errors.New("db error"))
	roleRouteManager := auth.NewRoleRouteManager(nil, repo)
	logger := logger.NewLogger()
	authService := auth.NewService(
		dbMock,
		repo,
		passwordEncoder,
		nil,
		logger,
		roleRouteManager,
	)
	err = authService.Register(context.Background(), auth.RegisterParams{
		Email:    "test@test.com",
		Password: "123456789",
	})
	assert.Equal(err, errors.New("cannot add role to user"))
	repo.AssertExpectations(t)
}

func TestService_Register_Successful(t *testing.T) {
	assert := assert.New(t)
	dbMock, err := pgxmock.NewPool()
	assert.Nil(err)
	defer dbMock.Close()
	dbMock.ExpectBegin()
	dbMock.ExpectCommit()
	passwordEncoder := core.NewPasswordEncoder()
	repo := new(mocks.MockQuerier)
	repo.EXPECT().GetUserByEmail(mock.Anything, mock.Anything, "test@test.com").Once().Return(repository.User{}, pgx.ErrNoRows)
	repo.EXPECT().GetRoleByCode(mock.Anything, mock.Anything, auth.RoleEndUser).Once().Return(repository.Role{ID: 1, Code: auth.RoleEndUser}, nil)
	repo.EXPECT().CreateUser(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.CreateUserParams)
			if p.Email != "test@test.com" {
				return false
			}
			if passwordEncoder.CompareHashAndPassword(p.Password, "123456789") != nil {
				return false
			}
			return true
		}),
	).Once().Return(repository.User{ID: 1}, nil)
	repo.EXPECT().AddRoleToUser(
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(input interface{}) bool {
			p := input.(repository.AddRoleToUserParams)
			if p.UserID != 1 {
				return false
			}
			if p.RoleID != 1 {
				return false
			}
			return true
		}),
	).Once().Return(nil)
	roleRouteManager := auth.NewRoleRouteManager(nil, repo)
	logger := logger.NewLogger()
	authService := auth.NewService(
		dbMock,
		repo,
		passwordEncoder,
		nil,
		logger,
		roleRouteManager,
	)
	err = authService.Register(context.Background(), auth.RegisterParams{
		Email:    "test@test.com",
		Password: "123456789",
	})
	assert.Nil(err)
	repo.AssertExpectations(t)
}

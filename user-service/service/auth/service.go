package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/majidmohsenifar/heli-tech/user-service/core"
	"github.com/majidmohsenifar/heli-tech/user-service/repository"
	"github.com/majidmohsenifar/heli-tech/user-service/service/jwt"
)

var (
	ErrEmailAlreadyExist         = errors.New("email already exist")
	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
	ErrInvalidToken              = errors.New("invalid token")
	ErrAccessDenied              = errors.New("access denied")
)

type Service struct {
	db               core.PgxInterface
	repo             repository.Querier
	passwordEncoder  *core.PasswordEncoder
	jwtService       *jwt.Service
	logger           *slog.Logger
	roleRouteManager *RoleRouteManager
}

type RegisterParams struct {
	Email    string
	Password string
}

type LoginParams struct {
	Email    string
	Password string
}

type GetUserDataByTokenParams struct {
	Token string
	Path  string
}

type GetUserDataByTokenResponse struct {
	Email string
	ID    int64
}

func (s *Service) Register(ctx context.Context, params RegisterParams) error {
	_, err := s.repo.GetUserByEmail(ctx, s.db, params.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		s.logger.Error("cannot check if email already exist", err)
		return fmt.Errorf("cannot check if email already exist")
	}
	if err == nil {
		return ErrEmailAlreadyExist
	}
	encodedPasswordBytes, err := s.passwordEncoder.GenerateFromPassword(params.Password)
	if err != nil {
		s.logger.Error("cannot generate password", err)
		return fmt.Errorf("cannot hash the password")
	}

	//get the default role
	role, err := s.roleRouteManager.GetRoleByCode(ctx, RoleEndUser)
	if err != nil && !errors.Is(err, ErrDefaultRoleDoesNotExist) {
		s.logger.Error("cannot get endUser role", err)
		return fmt.Errorf("cannot get the default role")
	}
	if errors.Is(err, ErrDefaultRoleDoesNotExist) {
		return ErrDefaultRoleDoesNotExist
	}
	dbTx, err := s.db.Begin(ctx)
	if err != nil {
		dbTx.Rollback(ctx)
		s.logger.Error("cannot start db transaction", err)
		return fmt.Errorf("something went wrong")
	}

	u, err := s.repo.CreateUser(ctx, dbTx, repository.CreateUserParams{
		Email:    params.Email,
		Password: string(encodedPasswordBytes),
	})
	if err != nil {
		dbTx.Rollback(ctx)
		s.logger.Error("cannot create user", err)
		return fmt.Errorf("cannot create user")
	}
	err = s.repo.AddRoleToUser(ctx, dbTx, repository.AddRoleToUserParams{
		UserID: u.ID,
		RoleID: role.ID,
	})
	if err != nil {
		dbTx.Rollback(ctx)
		s.logger.Error("cannot add role to user", err)
		return fmt.Errorf("cannot add role to user")
	}
	err = dbTx.Commit(ctx)
	if err != nil {
		dbTx.Rollback(ctx)
		s.logger.Error("cannot commit db transaction", err)
		return fmt.Errorf("something went wrong")
	}
	return nil
}

func (s *Service) Login(ctx context.Context, params LoginParams) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, s.db, params.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		s.logger.Error("cannot get user from db", err)
		return "", fmt.Errorf("something went wrong")
	}
	if err == pgx.ErrNoRows {
		return "", ErrInvalidUsernameOrPassword
	}
	err = s.passwordEncoder.CompareHashAndPassword(user.Password, params.Password)
	if err != nil {
		return "", ErrInvalidUsernameOrPassword
	}
	token, err := s.jwtService.GenerateToken(user.Email)
	if err != nil {
		s.logger.Error("cannot generate token", err)
		return "", fmt.Errorf("cannot generate token")
	}
	return token, nil
}

func (s *Service) GetUserDataByToken(ctx context.Context, params GetUserDataByTokenParams) (GetUserDataByTokenResponse, error) {
	email, err := s.jwtService.GetUsernameFromToken(params.Token)
	if err != nil {
		return GetUserDataByTokenResponse{}, ErrInvalidToken
	}
	user, err := s.repo.GetUserByEmail(ctx, s.db, email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		s.logger.Error("cannot get user by email", err)
		return GetUserDataByTokenResponse{}, fmt.Errorf("something went wrong")
	}
	if err == pgx.ErrNoRows {
		return GetUserDataByTokenResponse{}, ErrAccessDenied
	}
	hasAccess, err := s.roleRouteManager.HasUserAccessToRoute(ctx, user.ID, params.Path)
	if err != nil {
		s.logger.Error("cannot check user access", err)
		return GetUserDataByTokenResponse{}, fmt.Errorf("cannot check user access")
	}
	if !hasAccess {
		return GetUserDataByTokenResponse{}, ErrAccessDenied
	}
	return GetUserDataByTokenResponse{
		Email: email,
		ID:    user.ID,
	}, nil
}

func NewService(
	db core.PgxInterface,
	repo repository.Querier,
	passwordEncoder *core.PasswordEncoder,
	jwtService *jwt.Service,
	logger *slog.Logger,
	roleRouteManager *RoleRouteManager,
) *Service {
	return &Service{
		db:               db,
		repo:             repo,
		passwordEncoder:  passwordEncoder,
		jwtService:       jwtService,
		logger:           logger,
		roleRouteManager: roleRouteManager,
	}
}
